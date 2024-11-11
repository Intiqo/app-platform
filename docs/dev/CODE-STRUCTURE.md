# Code Structure

We make use of the [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) with small tweaks in naming conventions.

## Directory / File Structure

Below is the generated code structure with explanation on the thought process behind this code structure.

### bin

The application binary is built and copied to this folder.

### cmd

As per recommendation from the golang community, we store the entry file `main.go` inside the cmd folder.

### config

Any tool configuration files go here. For example, the nginx configuration and open api SDK generator configuration.

### docs

Most of the documentation related to the project is stored in this folder.

### internal

- This is the folder where all the code related to the platform including APIs, Services, Repositories, Dependent packages are stored. For example, code to send SMS, code to generate OTP, generate JWT token, upload file to S3, read secrets from AWS etc. are stored in individual folders under the `pkg` directory.
- Under the `domain` folder, we have all the models, constants and default values stored.
- The `dependency` folder contains the code related to dependency manager, which in our case is `wire`.
- The `database` folder contains all the code related to connecting to database.
- The `http` folder contains all the code related to the App's APIs.
- The `repository` folder contains all the code related to talking to the database. So, we'll typically create new repository for each table we create. At times, we may decide to use the same repository across multiple tables. This is decided on a case by case basis.
- The `service` folder contains all the business logic of the application. So, when developing a new feature, we either add new methods to existing services, or create new services for the new requirements.

### scripts

- This folder contains all the scripts required to run various aspects of the codebase including migrations, tests, docker and so on.

### tests

- This folder contains all the tests of the platform.

### others

- The other important files in the project root include the `Docker` related files, Makefiles each for Unix and Windows, and sample env file.

## Code Tree

```txt

```
