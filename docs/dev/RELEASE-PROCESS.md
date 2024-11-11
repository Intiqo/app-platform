# Introduction

This document explains the process we follow to make releases with the App platform.

## Update swagger specs

First, we update the swagger specification to ensure everything is in sync. So, run the below command:

### Linux / macOS

```sh
make doc
```

### Windows

```bat
make -f Makefile.win doc
```

## Upgrading dependencies

Next, we check if there are any [new pull requests from Dependabot](https://github.com/intiqo/app-platform/pulls). If there are, then we need to upgrade the dependencies with the below command for all the pull requests opened by Dependabot.

```sh
go get -u <package>@<new_version>
```

For example:

```sh
go get -u github.com/aws/aws-sdk-go-v2/service/s3@v1.58.0
```

## Updating the build number and version

Then, we update the version & build numbers in the following files:

- [scripts/build.bat](../../scripts/build.bat)
- [scripts/build.sh](../../scripts/build.sh)
- [scripts/cd-build.bat](../../scripts/cd-build.bat)
- [scripts/cd-build.sh](../../scripts/cd-build.sh)

We set the version number as yy.mm.wom. This represents year.month.week_of_month.

For example, if today is Monday, Jul 8th 2024, and we finished some development last week, then, the new version number would be 24.7.1 indicating 2024 as the year, 7 as the month and 1 as the week of the month, in this case being the 1st week of the month. The build number follows the same convention in the format `2407011` with 24 being the year, 07 being the month and 01 being the week of the month. The last digit (1 in this case) indicates the build number which will be incremented as we make any patches against the branch.

## Commit & Branch

### Commit Message

Now that the version & build numbers are set, we commit the code with the following message including the version number:

```sh
release: 24.7.1
```

Push the commit to the remote repositories.

### New Branch

Create a new branch in the below format:

```sh
releases/24.7.1
```

Version number as the branch name and releases is the sub-folder under which the new branch is created. Push the new branch to the upstream repository.

## SDK Release

### Staging / Testing

Once you are in the new branch, we need to publish a new version of the SDK as well as [explained here](./developing.md#publishing). Ensure that you publish the SDK as a beta release since we'll be in a testing phase with this branch by using the below command.

#### Linux / macOS

```sh
APP_SDK_VERSION=24.7.1-beta.1 make sdk-gen
```

#### Windows

```bat
APP_SDK_VERSION=24.7.1-beta.1 make -f Makefile.win sdk-gen
```

### Production

After publishing the staging SDK, we need to publish a new version for the production SDK. So, switch to last week's branch. In our case, that'll be `releases/24.6.4` and run the below command.

#### Linux / macOS

```sh
APP_SDK_VERSION=24.6.4 make sdk-gen
```

#### Windows

```bat
APP_SDK_VERSION=24.6.4 make -f Makefile.win sdk-gen
```

Notice how we’ve left out the `beta` part for the production release.

## GitHub Release

Our final step is to create GitHub releases for the two branches against staging & production to trigger automatic deployments.

So, head to [GitHub Releases](https://github.com/intiqo/app-platform/releases) and create a new release for both staging & production environments.

### Production

Draft a new release on GitHub. For our example, it’d be the following:

- Tag would be `24.6.4`
- Target branch would be `releases/24.6.4`
- Previous tag would be `24.6.3`
- Click on `Generate Release Notes`

SELECT THE OPTION `Set as the latest release`.

DO NOT SELECT THE OPTION `Set as a pre-release` as that will trigger a staging deployment.

Finally, click on the `Publish Release` button.

### Staging

Draft a new release on GitHub. For our example, it’d be the following:

- Tag would be `24.7.1-beta.1`
- Target branch would be `releases/24.7.1`
- Previous tag would be `24.6.4`
- Click on `Generate Release Notes`

SELECT THE OPTION `Set as a pre-release`.

DO NOT SELECT THE OPTION `Set as the latest release` as that will trigger a production deployment.

Finally, click on the `Publish Release` button.

## Patch Releases

At times, we create patch releases to fix critical issues. So, we’ll need to create releases to have the latest code deployed to production & staging.

For patch releases, we follow the same process as regular releases except for the following changes:

We increment the version number for every new patch release created

#### Staging

In our example, if we create two new patch releases for Staging, the new release numbers would be `24.7.1-beta.2` and `24.7.1-beta.3`.

#### Production

In our example, if we create two new patch releases for Production, the new release numbers would be `24.6.4-2` and `24.6.4-3`.

#### New SDK Versions

At times, if the swagger specification has changed, we’ll need to publish new versions of the SDK as well. So, checkout the appropriate branch(es) on your local system, pull the latest code and publish new SDK versions with the new version numbers as explained above.
