# Document of covid backend service

Run Via docker-compose up

this project use ELK for logging Error
the environment is in docker-compose.yml (incase you want edit can edit on this file)

the structure is combine between Clean Architecture, Hexagonal Architecture and Layer Architecture.

# Project Structure

## Description

This README outlines the structure and organization of the project, detailing its adherence to Clean Code principles and integration of various architectural concepts, including hexagonal architecture, layer architecture, and Domain-Driven Design (DDD). The aim is to provide a clear overview of the project's layout and architectural decisions to facilitate understanding and collaboration among developers.

## Architectural Concepts

The project structure has been designed to embody several key architectural concepts:

- **Clean Code**: Prioritizing readability, maintainability, and simplicity in code implementation.
- **Hexagonal Architecture**: Separating core business logic from external dependencies to enhance modularity and testability.
- **Layer Architecture**: Structuring code into distinct layers (e.g., presentation, domain, infrastructure) to promote separation of concerns and maintainability.
- **Domain-Driven Design (DDD)**: Focusing on modeling the domain and its behavior, with a clear separation between domain logic and infrastructure concerns.

## Layer Project

│ docker-compose.yml // Case local service use scripts/dev for Dockerfile build image
| .gitlab-ci.yml // Script for build image and deploy pod containner (CICD) use script build/deploy for Dockerfile build image and map script of deploy
│ go.mod // Init package dependency
│ go.sum // Init package dependency
│ README.md
│
├───apidocs
│ doc.json // Script Create Doc in this case create with Swagger
│
├───assets // store assets file like image,temp
├───migrations
│ V_1_1.sql // Script database for migrations
│
├───scripts // Script for deploy and build also crt and key
│ │ nginx.conf
│ │
│ ├───build
│ │ base.Dockerfile
│ │ Dockerfile
│ │
│ ├───deployment
│ │ configmap.yaml
│ │ deploy-backend.yaml
│ │ sonar.sh
│ │
│ ├───dev
│ │ Dockerfile
│ │
│ └───ssl
│ nginx-selfsigned.crt
│ nginx-selfsigned.key
│
└───src
├───api
│ ├───example // Presentation layer of example domain
│ │ └───handler
│ │ handler.go
│ │ handler_test.go
│ │ model.go
│ │
│ ├───middleware // Helper Middleware on Presentation layer
│ │ middleware.go
│ │
│ └───requestheader // Helper on Presentation layer
│ request_header.go
│
├───cmd
│ │ main.go // Init asset, session and call the application
│ │
│ └───config // Config for init
│ config.go
│ config_test.go
│
├───pkg
│ ├───entity // Entity struct date
│ │ └───example.go
│ │  
 │ ├───example // Package Business Logic Layer and Data Access Layer
│ │ │ model.go
│ │ │ repository.go // Data Access Layer
│ │ │ service.go // Business Logic Layer
│ │ │ service_test.go
│ │ │
│ │ └───mocks // mock Business Logic Layer and Data Access Layer for testing
│ │ Repository.go
│ │ UseCase.go
│ │
│ └───utils // Lib and dependency create package to easy to mock usecase for test
│ ├───convert
│ │ convert.go
│ │ convert_test.go
│ │
│ ├───error
│ │ error.go
│ │ error_test.go
│ │ mapping.go
│ │ mapping_test.go
│ │
│ ├───logger
│ │ │ elk.go
│ │ │ logger.go
│ │ │
│ │ └───mocks
│ │ Logger.go
│ │
│ ├───mocker
│ │ service.go
│ │
│ ├───recover
│ │ recover.go
│ │
│ └───sqlquery
│ sql_query.go
│ sql_query_test.go
│
└───routers // Port and adapter call on main to create server for RESTAPI
│ router.go
│
└───docs
docs.go // Auto gen Swagger
