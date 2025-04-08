## database

This package provides a database service for storing and retrieving data related to chat sessions. It uses a PostgreSQL database to store the data and offers functions for creating and managing tables, as well as handling database connections.

### internal/database/db.go
This file contains the implementation for the ChatSessionGraph struct, which is used to store and manage chat session data.

### internal/database/init.go
This file initializes the database connection and creates a new Handler instance. It requires a connection string to the PostgreSQL database as input.

### internal/database/service.go
This file defines the Service struct, which is responsible for managing database interactions. It includes functions for creating tables, handling database connections, and handling errors related to database operations.

### Configuration
The database package requires a connection string to the PostgreSQL database. This connection string can be provided through environment variables, command-line arguments, or configuration files.

### Edge Cases
The database package can be launched by providing the connection string to the PostgreSQL database. If the connection string is not provided, the package will attempt to use default values.