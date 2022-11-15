# Dagger-go-app-ci

## _Continous Integration for any golang project_
Run your Continous integration pipeline easily with dagger go.
The steps of this pipeline are:
- Test stage
- Vulnerability scanning with trivy
- Build docker image and push to docker hub
## Requirements

- Docker installed
- Golang

## Steps
First clone the app and move to the root of the app folder

rename the .env.example to .env with the command:
```sh
mv .env.example .env
```

In the .env file, put the correct informations by setting every env variable.

> You can also go to the config/env.go file
> and set up the values directly there

## Execution
After setting everything up, run :
```sh
go run main.go
```

There will be a file called CI.out (you can rename it if you wish :)), where you will be able to follow the CI pipeline steps as they are getting executed.

ENJOY !!!!
