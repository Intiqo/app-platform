# Architecture

## Application Architecture

- We use a customized version of [Clean Architecture](https://www.freecodecamp.org/news/a-quick-introduction-to-clean-architecture-990c014448d2/)
- Interfaces are used to define functionality on a service layer, and on a repository layer
- Handlers are created on the API layer, we don't create interfaces for these
- Database layer is built using [pgx](https://github.com/jackc/pgx) with mostly raw SQL queries. In some cases, we also use [squirrel](https://github.com/Masterminds/squirrel) for constructing SQL queries.
- We use [Swaggo](https://github.com/swaggo/swag) to generate and maintain API documentation using code comments on the controller methods.
- Application configuration is setup using [viper](https://github.com/spf13/viper) and a `.env` file on local systems. We use AWS Secrets typically on production systems.
- Integration tests are written using [testify](https://github.com/stretchr/testify)
