# gitwize-be

This repository contains the backend code of the GitWize project

---

## Table of Contents

- [Introduction](#introduction)
- [Local Installation](#local-installation)
  - [Prerequisites](#prerequisites)
  - [Run the Project Locally](#run-the-project-locally)
- [Deployment](#deployment)
- [Contribution](#contribution)
  - [Branching Model](#branching-model)
  - [Review Process](#review-process)
- [Versions](#versions)
- [Licenses](#licenses)

---

## Introduction

GitWize is an AWS cloud-based application that allows you to extract valuable insights and metrics from a single GitHub repository using git commit logs. Currently, GitWize integrates Okta as the primary identity management service for end-users. Either if you’re a delivery manager, engineering lead, or engineering manager, GitWize will enable you to:

- Improve your software quality
- Track your project’s progress on a daily, monthly, or quarterly basis 
- Track your team’s productivity, performance, and health
- Identify early potential risks and take early measures that can save money, time, and stress

Once you add a GitHub repository to GitWize, you will unlock the following features:
- Analyze PR statistics such as average PR size and PR rejection rate
- Calculate code change velocity over time in terms of the number of commits
- Display reports in either tabular or chart format.

[⇧ back to top](#table-of-contents)

---

## Local Installation

To set up this project on you local machine, follow the next linux-based instructions.

### Prerequisites

Ensure you comply with the following prerequisites before you follow the rest of the instructions.

* Install [Go](https://go.dev/doc/install) v1.14 or above \
As of August 2022, v1.14 is functional. To check the Go version that you have installed, type in your terminal: `go version`
* Install [pre-commit hook](https://pre-commit.com/#install). A framework for managing and maintaining multi-language pre-commit hooks. To install it using homebrew, type in your terminal: `brew install pre-commit `
* Run in your terminal `pre-commit install` to set up the git hook scripts.  \
Now the pre-commit will run automatically on `git commit`
* Install [Docker ](https://docs.docker.com/get-docker/)

### Run the Project Locally
To run the backend project on your machine:

1. Open Docker Desktop application
2. Clone the gitwize-be repository
3. Go to the gitwize-be root directory
4. Type in your terminal:
``` 
docker-compose -f ./docker/docker-compose.yaml up
```
5. In a different tab, type in your terminal the following command to run all test cases:
```
GW_DATABASE_SECRET_LOCAL=P@ssword123 
go test -count=1 ./...
```
6. In a different tab, type in your terminal the following command to run the project:
```
GW_DATABASE_SECRET_LOCAL=P@ssword123 
go run application.go
```

If you get the following error: 

```    
    # golang.org/x/sys/unix
    ../go/pkg/mod/golang.org/x/sys@v0.0.0-20200302150141-5c8b2ff67527/unix/syscall_darwin.1_13.go:25:3: //go:linkname must refer to declared function or variable
    ../go/pkg/mod/golang.org/x/sys@v0.0.0-20200302150141-5c8b2ff67527/unix/zsyscall_darwin_amd64.1_13.go:27:3: //go:linkname must refer to declared function or variable
    ../go/pkg/mod/golang.org/x/sys@v0.0.0-20200302150141-5c8b2ff67527/unix/zsyscall_darwin_amd64.1_13.go:40:3: //go:linkname must refer to declared function or variable
    ../go/pkg/mod/golang.org/x/sys@v0.0.0-20200302150141-5c8b2ff67527/unix/zsyscall_darwin_amd64.go:28:3: //go:linkname must refer to declared function or variable
    ../go/pkg/mod/golang.org/x/sys@v0.0.0-20200302150141-5c8b2ff67527/unix/zsyscall_darwin_amd64.go:43:3: //go:linkname must refer to declared function or variable
    ../go/pkg/mod/golang.org/x/sys@v0.0.0-20200302150141-5c8b2ff67527/unix/zsyscall_darwin_amd64.go:59:3: //go:linkname must refer to declared function or variable
    ../go/pkg/mod/golang.org/x/sys@v0.0.0-20200302150141-5c8b2ff67527/unix/zsyscall_darwin_amd64.go:75:3: //go:linkname must refer to declared function or variable
    ../go/pkg/mod/golang.org/x/sys@v0.0.0-20200302150141-5c8b2ff67527/unix/zsyscall_darwin_amd64.go:90:3: //go:linkname must refer to declared function or variable
    ../go/pkg/mod/golang.org/x/sys@v0.0.0-20200302150141-5c8b2ff67527/unix/zsyscall_darwin_amd64.go:105:3: //go:linkname must refer to declared function or variable
    ../go/pkg/mod/golang.org/x/sys@v0.0.0-20200302150141-5c8b2ff67527/unix/zsyscall_darwin_amd64.go:121:3: //go:linkname must refer to declared function or variable
    ../go/pkg/mod/golang.org/x/sys@v0.0.0-20200302150141-5c8b2ff67527/unix/zsyscall_darwin_amd64.go:121:3: too many errors
```

Run the following script to fix the `Go` source code version mismatching 

```
go get -u golang.org/x/sys
go: downloading golang.org/x/sys v0.0.0-20220804214406-8e32c043e418
go: upgraded golang.org/x/sys v0.0.0-20200302150141-5c8b2ff67527 => v0.0.0-20220804214406-8e32c043e418
```

As a result you obtain the following welcoming message:

```
    Hello from gitwize BE
    configuration is {Server:{Port:8080} Database:{GwDbName:gitwize GwDbUser:gitwize_user GwDbPassword:P@ssword123 GwDbHost:localhost GwDbPort:3306} Auth:{AuthDisable:true} Cypher:{PassPhase:} Endpoint:{Frontend:http://localhost:8080 SonarQubeServer:http://localhost} SonarQube:{AdminSecret: ScannerPath:/usr/local/sonarqube/sonar-scanner/bin/sonar-scanner PropertiesPath:/usr/local/sonarqube/sonar-scanner/conf/sonar-scanner.properties BaseDirectory:./}}
    [GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

    [GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
     - using env:	export GIN_MODE=release
     - using code:	gin.SetMode(gin.ReleaseMode)

    [GIN-debug] POST   /api/v1/admin/:op_id      --> gitwize-be/src/controller.posAdminOperation (4 handlers)
    [GIN-debug] GET    /api/v1/repositories/     --> gitwize-be/src/controller.getListRepos (5 handlers)
    [GIN-debug] GET    /api/v1/repositories/:id  --> gitwize-be/src/controller.getRepos (5 handlers)
    [GIN-debug] POST   /api/v1/repositories/     --> gitwize-be/src/controller.postRepos (5 handlers)
    [GIN-debug] DELETE /api/v1/repositories/:id  --> gitwize-be/src/controller.delRepos (5 handlers)
    [GIN-debug] GET    /api/v1/repositories/:id/stats --> gitwize-be/src/controller.getStats (5 handlers)
    [GIN-debug] GET    /api/v1/repositories/:id/contributor --> gitwize-be/src/controller.getContributorStats (5 handlers)
    [GIN-debug] GET    /api/v1/repositories/:id/impact/weekly --> gitwize-be/src/controller.getWeeklyImpact (5 handlers)
    [GIN-debug] GET    /api/v1/repositories/:id/code-velocity --> gitwize-be/src/controller.getCodeChangeVelocity (5 handlers)
    [GIN-debug] GET    /api/v1/repositories/:id/trends --> gitwize-be/src/controller.getStatsQuarterlyTrends (5 handlers)
    [GIN-debug] GET    /api/v1/repositories/:id/pullrequest-size --> gitwize-be/src/controller.getPullRequestSize (5 handlers)
    [GIN-debug] GET    /api/v1/repositories/:id/code-quality --> gitwize-be/src/controller.getCodeQuality (5 handlers)
    [GIN-debug] Environment variable PORT is undefined. Using port :8080 by default
    [GIN-debug] Listening and serving HTTP on :8080
```

[⇧ back to top](#table-of-contents)

### Run the Project on the Cloud

To run all test cases, type in your terminal:  
```
GW_DEPLOY_ENV=DEV  
GW_DATABASE_SECRET_DEV=database_secret  
go test ./… 
```

To run the project, type in your terminal: 
```
GW_DEPLOY_ENV=DEV  
GW_DATABASE_SECRET_DEV=database_secret  
go run application.go`
```

[⇧ back to top](#table-of-contents)

---

## Contribution

We encourage you to contribute to GitWize! Please check out our [Contributing guide](CONTRIBUTING)

[⇧ back to top](#table-of-contents)


### Branching Model

The branching model follows this convention: `<feature|bugfix>-<ticket-id>-<short-description>`  \
Example: `feature/GW-1-skeleton-code` and `bugfix/GW-10-some-blocker`

[⇧ back to top](#table-of-contents)


### Review Process

When sending a PR, ensure you include the following information:
```
What does this PR do?
Where should the reviewer start?
Screenshots & link (if appropriate)
Questions
```

[⇧ back to top](#table-of-contents)

---

## Versions

Versioning follows this convention: `<major>.<minor>.<buildnumber>` \
Example: 1.0.1 and 1.0.11

[⇧ back to top](#table-of-contents)

---

## Licenses

@[MIT](LICENSE) and @[CLA](CLA.md)

[⇧ back to top](#table-of-contents)