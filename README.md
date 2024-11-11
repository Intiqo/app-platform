# Introduction

This is a template to create golang APIs.

Welcome to App Platform codebase. Here, we explain how you can setup the App API on your local system using docker.

## Developing

If you are a App API developer, head to [Developing Guide](docs/dev/DEVELOPING.md)

## Setup

This is the backend platform for App built with the following set of tools and technologies:

- [Golang](https://golang.org) - Programming Language
- [Postgresql](https://www.postgresql.org) - Database
- [OpenAPI](https://www.openapis.org) - API Specification
- [Pgx](https://github.com/jackc/pgx) - Database Driver
- [Echo](https://echo.labstack.com/) - Web framework
- [Docker](https://www.docker.com) - Containerization

## Audience

This guide is for people who want to run App API through a Docker container.

## Dependencies & Pre-requisites

Ensure you have the following installed and running:

### Docker

We use `Docker` to run the App API in a container.

- Install [Docker](https://www.docker.com) for your OS
- Once installed, start Docker and ensure it is running

### Jq

We use `jq` in some of our scripts, so, [install it from here for your OS](https://jqlang.github.io/jq/download/).

### MKCert

We use `mkcert` to generate SSL certificates that are trusted on the local development system. This allows us to serve requests over https for the APIs.

- Install [mkcert](https://github.com/FiloSottile/mkcert) by following the instructions for your OS
- Run the below command to make mkcert a root signing authority on your system

```bash
mkcert -install
```

- Now, run the below command to generate the SSL certificate for App. You'll need to do this from the project root

```bash
mkcert local.api.app.co
```

- Create a new directory called `.certs` and move the certificates generated from the above command to this new directory.

### Domain Mapping

- To serve requests over https on a proper domain, we need to map the domain `local.api.app.co` onto the local system network
- Copy the below content to `/etc/hosts` on Linux or macOS. And, on Windows, you need to copy it to `C:\Windows\system32\drivers\etc\hosts` file.

```txt
127.0.0.1 local.api.app.co
```

## Environment Configuration

We need to provide some configuration details for App API to run. So, run the below command to initiate the setup process:

### Linux / macOS

```bash
make setup
```

### Windows

```bat
make -f Makefile.win setup
```

Note: On Windows, you'll need to install the `make` package supported by a package manager `chocolatey`

### Alternative

As an alternative, you can simply copy the `sample.env` as `.env` and update all the properties as per your system configuration.

## Update Configuration

- Update the configuration values in the .env file that is generated.
- Change the value for `DB_HOST` to `appdb`  since that is the container name configured in docker.
- For some of the properties (eg: MSG91 related configuration), you'll need to get the appropriate values from those systems, or, talk to the Team lead.

## AWS Setup

We use AWS services for various features such as uploading files to S3, sending SMS using AWS SNS and so on. So, to ensure all the features work fine, we need to setup AWS SSO and have docker read the environment variables related to AWS keys before starting the App platform service.

Follow the below steps to have AWS profile for App setup on your system:

- Request your Team Lead to give you AWS IAM Federated access if you don't have it already
- [Install the AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html) for your OS.
- After the AWS CLI is installed, make sure you open a new terminal window to get the `aws` command working.
- Now type the below command in your terminal

```bash
aws configure sso
```

- Enter the below details for each of the steps:

```bash
SSO session name (Recommended): App
SSO start URL [None]: https://myapp.awsapps.com/start/#/
SSO region [None]: ap-south-1
SSO registration scopes [sso:account:access]: <just press enter>
```

- After the above, a browser tab will be opened so you need to login to your AWS account created in Step 1 and click on "Allow Access"
- Then, go back to the terminal and enter the following details:

```bash
CLI default client Region [None]: ap-south-1
CLI default output format [None]: json
CLI profile name [AdministratorAccess-217199241424]: App
```

That's it! That should have the AWS App account mapped on your system.

## Running

We use [Docker Compose](docker-compose.yaml) to spin up the following services:

### Infrastructure

- Database (PostgreSQL)
- Nginx Server (Used to serve https requests through reverse proxy)
- API Service

### Command Line

Run the below command (from the project root):

---
macOS / Linux

```bash
sh scripts/start-docker.sh
```

Alternatively, you can also run the below command if your system supports the `make` command.

```bash
make start-docker
```

---
Windows

```cmd
scripts\start-docker.bat
```

Alternatively, you can also run the below command if your system supports the `make` command.

```cmd
make -f Makefile.win start-docker
```

- When you run the above command, we try to get a new AWS sso session, so, you'll be taken to the browser to login to your AWS account, do so and click on "Allow Access" to get new tokens.
- Check Docker Desktop and open the `api` service logs to see if the server has started. If everything is successful, you should see a message stating `API Server Started` in the logs.

## Accessing the API Documentation

- If everything works, the platform should be up & running at `https://local.api.app.co`
- Open the [Swagger Documentation](https://local.api.app.co/swagger/index.html) in your web browser to view the API documentation
