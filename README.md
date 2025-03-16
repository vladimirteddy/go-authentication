# Description

The project for authentication with GO.
This project provides a robust authentication system built with Go, featuring:

## Features

- User registration and login functionality
- Secure password hashing and storage
- Email verification
- Session management
- PostgreSQL database integration using Goose for migrations

## Migration Management

Database migrations are handled using Goose with support for:

- Running migrations up/down
- Migrating to specific versions
- Environment-specific configurations (DEV/UAT/PROD)
- GitHub Actions workflow for automated deployments

## Getting Started

1. Set your database connection string in environment variable `GOOSE_DBSTRING`
2. Run migrations: `make up`
3. To rollback: `make down` \
   Also used up-to and down-to \
   Please refer to:https://github.com/pressly/goose?tab=readme-ov-file#up-to
