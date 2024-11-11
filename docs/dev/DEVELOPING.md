# Development Setup

- Welcome! We'll be developing the App Platform and it's APIs in this repository.
- Follow the instructions below to set your development environment up and get started.

## Pre-requisites

### Golang

We use golang as the primary programming language for the platform.

- Download and install Golang from [here](https://go.dev/doc/install) for your platform

### Node.js

We use Node.js to create an SDK for the APIs. This is used by the clients to consume the APIs.

- You can download & install Node.js from [the official website](https://nodejs.org/en/download)
- Alternatively, you can use Node Version Manager (nvm) which allows you to manage multiple versions of Node.js on your system
  - For macOS / Linux, check [nvm](https://github.com/nvm-sh/nvm)
  - For Windows, check [nvm-windows](https://github.com/coreybutler/nvm-windows)
- For our requirement, anything above Node.js 18 should work

### PostgreSQL

We use PostgreSQL as the primary database for the platform.

- Install and setup PostgreSQL on your system by following the instructions for your OS [here](https://www.postgresql.org/download/)
- Check for alternative installation options in the above link if you wish to use tools like Homebrew, Docker, etc.

#### Creating a database

- Once you've installed PostgreSQL, you'll need to create a database for App to use
- See [this documentation](https://www.tutorialspoint.com/postgresql/postgresql_create_database.htm) to create a database

There is plenty of documentation available online to help you with any issues that you may run into.

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

### Domain Mapping

- To serve requests over https on a proper domain, we need to map the domain `local.api.app.co` onto the local system network
- Copy the below content to `/etc/hosts` for macOS / Linux. For windows, copy it to `C:\Windows\System32\drivers\etc\hosts`

```txt
127.0.0.1 local.api.app.co
```

### Caddy

We use caddy server to serve https requests through a reverse proxy

- Install Caddy server by following the instructions for your OS [here](https://caddyserver.com/docs/install)
- Depending on which method you took to install caddy, locate or create the `Caddyfile`
- Copy the below content to the Caddyfile

```caddy
local.api.app.co {
    tls /usr/local/etc/certs/local.api.app.co.pem /usr/local/etc/certs/local.api.app.co-key.pem
    reverse_proxy localhost:4455
}
```

- Now, copy the below content to `/etc/hosts`

```hosts
127.0.0.1 local.api.app.co
```

### Windows setup instructions

- On Windows, we need to use [WSL (Windows Subsystem for linux)](https://learn.microsoft.com/en-us/windows/wsl/install) that allows us to run unix commands on Windows and install linux packages that are unavailable on Windows. For example, Redis. Although it is supported by 3rd party packages like choclatey, the latest versions are not available on the package managers and lag behind most of the time.
- So, [install WSL Ubuntu package](https://learn.microsoft.com/en-us/windows/wsl/install#install-wsl-command)
- Also, install [Chocolatey package manager](https://chocolatey.org/install) on Windows to have support to install other tools.
- Install [make tool](https://community.chocolatey.org/packages/make) with chocolatey which will be used to run make commands on Windows.

Subsequently, to run any make command from the project root, use the below format on Windows:

```bat
make -f Makefile.win <command_name>
```

For example:

```bat
make -f Makefile.win doc
```

## AWS Setup

- We use AWS services for various features such as uploading files to S3, sending SMS using AWS SNS and so on.
- So, to ensure all the features work fine, we need to have AWS SSO configured on the host system.
- Once you get AWS IAM access to App's account, setup the SSO using the [instructions here](https://docs.aws.amazon.com/cli/latest/userguide/sso-configure-profile-token.html).
- Make sure that the below details are configured:
  - SSO session name: App
  - SSO start URL: <https://intiqo.awsapps.com/start/#/>
  - SSO region: ap-south-1
  - Just enter a blank line for SSO registration scopes
  - CLI default client Region: ap-south-1
  - CLI default output format: json
  - CLI profile name: App
- The SSO session is designed to work for a few hours before it expires. So, when the session expires (and you get an error running the platform), you'll need to refresh the tokens using the command `aws sso login --profile App` and logging in again on AWS.
- Finally, remove the keys `AWS_REGION`, `AWS_ACCESS_KEY_ID` & `AWS_SECRET_ACCESS_KEY` from the `.env` file.

## Migrations

### Goose

- We use [goose](https://github.com/pressly/goose) for migrations. So install it.
- If you need to make database schema changes, you'll need to add migrations for the same.
- To create a new migration, run the below command:

```bash
goose -dir ./internal/database/migrations create <name_of_migration> sql
```

- Once the migration file is created, fill it in with the schema changes as required
- Do not forget to fill in the down migrations as well
- Once the above is done, run the below command to apply the migrations

```bash
make migrate
```

- In case you need to revert a migration, change `up` to `down` in the `migrate.sh` file and run the below command.

```bash
make migrate
```

Note: Do not forget to change `down` back to `up` once you've run the down migration.

## Dependency Injection

- We use Google's [wire dependency manager](https://github.com/google/wire) to manage dependencies, so, install it.
- When you create a new service or repository or package, or, when new dependencies are added, you'll need to run the below command

```bash
make wire
```

- This will generate a new version of the `wire_gen.go` file and wire the new dependencies to all the modules in the application.

## Running the app

`Note: Ensure that you've started the Caddy server`

### Command Line

```bash
make start
```

- This starts the api service in the background.
- If everything works, the platform should be up & running at `https://local.api.app.co`

### Visual Studio Code

- There are a bunch of [recommended extensions](../../.vscode/extensions.json), ensure you install them for a better experience
- We have configurations for the following in `launch.json`:
  - [Start App Api](../../.vscode/launch.json)
  - If everything works, the platform should be up & running at `https://local.api.app.co`

## Accessing the API

- Open the [API Documentation](https://local.api.app.co/swagger/index.html) in your web browser

## Tests

- To run all the tests, run the below command

```bash
make test-cover
```

## Documentation

### Swagger API Specs

- Swagger API specs essentially represent the complete documentation for `app`'s APIs exposed to clients
- You can access this at `https://local.api.app.co` once you start the app with `make start`
- When you have added new models or APIs or routes, you'll need to run `make doc` to update the API specs and the Swagger documentation

### Client SDK

#### Preface

- We make use of [OpenAPI Generator](https://openapi-generator.tech/) that allows us to generate client SDKs for Web & Mobile.
- Currently, we only support an SDK for Javascript / Typescript with more clients to be supported in the future
- The Javascript / Typescript SDK is published at [this repository](https://github.com/intiqo/app-sdk-ts)

#### Setup

- Install the OpenAPI Generator by referring to the [installation instructions](https://openapi-generator.tech/docs/installation/)
- [Create a new Developer access token in GitHub](https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token)
- Update your `~/.npmrc` file (create one if it doesn't exist) to add the below lines. Replace `<access_token>` with the token you generated above.

```shell
@intiqo:registry=https://npm.pkg.github.com
//npm.pkg.github.com/:_authToken=<access_token>
```

- Check out [this repository](https://github.com/intiqo/app-sdk-ts) next to `app-platform` (both of them should be in the same folder.)

#### Publishing

- Whenever you update the swagger specs (i.e. add new APIs / modify existing APIs), you will need to update the SDK by regenerating it.
- To do so, run the command `APP_SDK_VERSION=<version_number> make sdk-gen` which generates the sdk, builds javascript from typescript and publishes the package to the GitHub private npm registry under `@intiqo/app-sdk-ts` package
- In the above, replace `<version_number>` with the appropriate version number you need. For example: `23.7.2-beta.1` if you are still developing it, or, `23.7.2` if it is ready to be consumed by the clients

### Package Documentation

- Package documentation is essentially for fellow developers who are new to the project and want to understand what packages are available and their corresponding functionality
- To enable this, follow the below set of commands

```shell
# Downloads the godoc package
$ go get golang.org/x/tools/cmd/godoc

# Runs the documentation server against app platform
$ godoc -http=:6060
```

- Once the server is up and running, visit `http://localhost:6060/pkg/github.com/intiqo/app/` to check all the package level documentation for app

## Running Tasks

- You can also run any make command such as `make wire`, or `make doc` using [VSCode Tasks](../../.vscode/tasks.json)

## Architecture & Code Structure

Before getting your environment setup, let us understand how the code is structured and how it works.

- [Read this document](./architecture.md) to understand the high level architecture
- [Read this document](./code-structure.md) to understand how we organize the code

## Release Process

If you are responsible for creating releases, check the [Release Process](./release-process.md) documented.
