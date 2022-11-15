package main

import (
	"context"
	"fmt"
	"golang-app-ci/config"
	"os"

	"dagger.io/dagger"
)

const IMAGE_TAG = "1.1.1"

func main() {
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

	codeSrc := client.Git(config.GIT_REPOSITORY_SSH).Branch(config.GIT_BRANCH).Tree()

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

	trivyContainer := client.Container().From("bitnami/trivy:0.34.0-debian-11-r4").WithEnvVariable("GITHUB_TOKEN", config.GITHUB_TOKEN)

	resultVunlScan := trivyContainer.Exec(dagger.ContainerExecOpts{
		Args: []string{"repo", "--no-progress", "--branch", "main", config.GIT_REPOSITORY_URL},
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

	address, err := goApp.Publish(ctx, config.PUBLISH_ADDRESS+":"+IMAGE_TAG, dagger.ContainerPublishOpts{})
	if err != nil {
		fmt.Println("Could not publish container to:", config.PUBLISH_ADDRESS)
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Container published at:", address)
}
