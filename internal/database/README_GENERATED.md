## database_test

This package contains tests for the database service. It uses a PostgreSQL test instance to verify the functionality of the service, including updating and retrieving the history for various issue IDs, repository names, and models.

### File structure
- db.go
- db_test.go
- init.go
- service.go
- internal/database/db_test.go

### Code summary
The `database_test` package imports necessary libraries for database interaction and testing. It defines two main functions:

1. `expectMessage`: This function retrieves the history for a given issue ID, repository name, and model. It then checks if the message role is AI, if the message part is of type TextContent, and if the message text matches the expected message. If any of these checks fail, it prints a fatal error.

2. `Test_DB`: This function sets up a PostgreSQL test instance and creates a new AI service instance. It then updates the history for several issue ID, repository name, and model combinations, drops the history for specific combinations, and calls `expectMessage` to verify the expected messages. Finally, it retrieves the history for a specific combination and checks if the content is empty or not. If any of these checks fail, it prints a fatal error.

The package tests the database service by interacting with a PostgreSQL test instance and verifying the expected behavior of the service.