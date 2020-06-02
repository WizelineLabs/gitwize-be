### Init Setup
- [Install GoLang](https://golang.org/doc/install)
- Go version: 1.14

- [Pre-commit hook](https://pre-commit.com/)

  Install using homebrew `brew install pre-commit` or using pip `pip install pre-commit`

  Install git hook `pre-commit install`

### Run All Test Cases
# On local environment

`docker-compose -f ./docker/docker-compose.yaml up`

`export GW_DATABASE_SECRET_LOCAL=P@ssword123`

`go test ./...`

To run tests without cache:

`go test -count=1 ./...`

# On cloud environment
`GW_DEPLOY_ENV=DEV GW_DATABASE_SECRET_DEV=database_secret go test ./...`

### Run application
# On local environment
`docker-compose -f ./docker/docker-compose.yaml up`
`GW_DATABASE_SECRET_LOCAL=P@ssword123 go run application.go`
# On cloud environment
`GW_DEPLOY_ENV=DEV GW_DATABASE_SECRET_DEV=database_secret go run application.go`
