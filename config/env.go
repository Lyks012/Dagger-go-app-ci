package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	GITHUB_TOKEN        = ""
	GIT_REPOSITORY_URL  = ""
	GIT_REPOSITORY_SSH  = ""
	GIT_REPOSITORY_NAME = ""
	GIT_BRANCH          = ""
	PUBLISH_ADDRESS     = ""
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Could not load env:", err)
	}

	GITHUB_TOKEN = getStringEnv("GITHUB_TOKEN", GITHUB_TOKEN)
	GIT_REPOSITORY_URL = getStringEnv("GIT_REPOSITORY_URL", GIT_REPOSITORY_URL)
	GIT_REPOSITORY_SSH = getStringEnv("GIT_REPOSITORY_SSH", GIT_REPOSITORY_SSH)
	GIT_REPOSITORY_NAME = getStringEnv("GIT_REPOSITORY_NAME", GIT_REPOSITORY_NAME)
	GIT_BRANCH = getStringEnv("GIT_BRANCH", GIT_BRANCH)
	PUBLISH_ADDRESS = getStringEnv("PUBLISH_ADDRESS", PUBLISH_ADDRESS)

}

func getStringEnv(envKey string, defaultValue string) string {
	envValue, isSet := os.LookupEnv(envKey)
	if isSet {
		return envValue
	}
	return defaultValue
}
