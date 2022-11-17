package main

import (
	"context"
	"errors"
	"fmt"
	appconfig "golang-app-ci/config"
	"os"
	"time"

	"dagger.io/dagger"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/storage/memory"
)

func main() {
	imageTag, err := getLastCommitHash()
	if err != nil {
		fmt.Println("Could not get last commit tag from github:", err)
		fmt.Println("Using default commit tag :", imageTag)
	}

	ctx := context.Background()

	outputFile, err := os.Create("CI.out")
	if err != nil {
		fmt.Println("Could not create ouput file:", err)
		os.Exit(1)
	}

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(outputFile))
	if err != nil {
		fmt.Println("Could not connect to dagger client:", err)
		os.Exit(1)
	}
	fmt.Println("Client connected successfully !")

	codeSrc := client.Git(appconfig.GIT_REPOSITORY_SSH).Branch(appconfig.GIT_BRANCH).Tree()
	src := client.Host().Workdir()
	codeSrc = codeSrc.WithFile("Dockerfile", src.File("Dockerfile"))

	// test step
	fmt.Println("Testing ...")

	goBaseImage := client.Container().From("golang:latest").WithMountedDirectory("/src", codeSrc).WithWorkdir("/src")

	testResult := goBaseImage.Exec(dagger.ContainerExecOpts{
		Args: []string{"go", "test", "./..."},
	})

	testExitCode, err := testResult.ExitCode(ctx)
	if testExitCode != 0 {
		fmt.Println("Could not perform tests!")
		os.Exit(1)
	}
	if err != nil {
		fmt.Println("Test step failed:", err)
		// os.Exit(1)
	}

	fmt.Println("Testing finished.")

	// vuln scan with trivy
	fmt.Println("Scanning with trivy")

	trivyContainer := client.Container().From("bitnami/trivy:0.34.0-debian-11-r4").WithEnvVariable("GITHUB_TOKEN", appconfig.GITHUB_TOKEN).WithEnvVariable("CACHEBUSTER", time.Now().String())

	resultVunlScan := trivyContainer.Exec(dagger.ContainerExecOpts{
		Args: []string{"repo", "--no-progress", "--branch", appconfig.GIT_BRANCH, appconfig.GIT_REPOSITORY_URL},
	})

	vulnScanExitcode, err := resultVunlScan.ExitCode(ctx)
	if vulnScanExitcode != 0 {
		fmt.Println("Could not perform vulnerability scan.")
		os.Exit(1)
	}
	if err != nil {
		fmt.Println("Vuln scan step failed:", err)
		// os.Exit(1)
	}

	fmt.Println("Scan completed.")

	// publish
	fmt.Println("Publishing image in docker registry ...")

	goApp := client.Container().Build(codeSrc)

	address, err := goApp.Publish(ctx, appconfig.PUBLISH_ADDRESS+":"+imageTag, dagger.ContainerPublishOpts{})
	if err != nil {
		fmt.Println("Could not publish container to:", appconfig.PUBLISH_ADDRESS)
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Container published at:", address)
}

func getLastCommitHash() (string, error) {
	commitHash := "1.0.0"
	defaultHashUsed := true

	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{appconfig.GIT_REPOSITORY_SSH},
	})

	refs, err := rem.List(&git.ListOptions{})
	if err != nil {
		return commitHash, err
	}

	for _, ref := range refs {
		if ref.Name().String() == ("refs/heads/" + appconfig.GIT_BRANCH) {
			commitHash = ref.Hash().String()
			defaultHashUsed = false
		}
	}

	if defaultHashUsed {
		return commitHash, errors.New("Did not find refs/heads/" + appconfig.GIT_BRANCH + ". Maybe branch does not exist!")
	}
	return commitHash, nil
}
