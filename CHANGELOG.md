# FB-Server Changelog

All notable changes to this project will be documented in this file.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased] v0.3.3

## 20 Sep 2024

### Added

-   /health endpoint to check services (auth / event / fighters / gateway) status
-   HealthStatus model for services
-   HealthStatus to proto converter method
-   HealthStatus from proto converter method
-   HealthResponse message for proto file
-   HealthCheck method for Auth/Event/Fighters services in proto file

## 28 Aug 2024

### Changed

-   Fightbettr name was replaced to Pickfighter
-   Fb prefix was replaced to Pf

## Released [v0.3.2]

## 31 Jul 2024

### Added

-   Fighters service: cmd tests
-   Fighters service: pkg/errors tests
-   Fighters service: pkg/cfg test
-   Fighters service: pkg/model tests
-   Fighters service: internal/service/fighters tests
-   Fighters service: internal/service/handler/grpc tests
-   Fighters service: internal/service/controller/fighters tests
-   Fighters service: internal/service/repository/psql tests
-   Fighters service: added new error codes
-   Fighters service: added script for mockgen
-   Fighters service: added directory gen/mocks
-   Fighters service: added viper test config generator
-   Added tests directory in root project with Dockerfile and init.sql for test database creation
-   Added script to run docker container with test db

### Changed

-   Fighters service: Error field ErrCode changed to Code
-   Fighters service: Timestamp field in error struct is string now
-   Fighters service: pkg/utils moved to cmd package
-   Fighters service: psql.New constructor needs config in arguments now
-   Fighters service: changed logs
-   Fighters service: config argument is required for WriteFighterData and DeleteFighterData methods

### Removed

-   Fighters service: removed few error codes

## Released [v0.3.0]

## 12 Jul 2024

### Added

-   Fightbettr Service as Gateway
-   Auth Service
-   Events Service
-   Scrapper Service
-   Fighters Service
-   Registry script
-   Proto script
-   api directory with proto file
-   gen directory with generated grpc files

### Changes

-   Logger
-   Now app works as microservice system

### Removed

-   fb-service app directory
-   Tests
-   Mocs
-   Tests scripts

## Released [v0.2.0]

## 24 March 2024

### Added

-   Tests for /internal/repo/auth
-   Added FbFightersRepo interface

### Changes

-   AuthRepo mocks

## 12 March 2024

### Added

-   Tests for /internal/auth package
-   FbAuth interface
-   AuthRepo mocks

### Changes

-   Script for packages testing. Scrapper package no longer counts

## 6 March 2024

### Added

-   Tests for /internal/common package
-   CommonRepo interface
-   CommonRepo mocks
-   Router interface

### Changes

-   Api mocks have been changed
-   CommonService interface has new methods
-   ApiService mock was regenerated
-   Pgxs mock was regenerated

## 29 Feb 2024

### Added

-   Tests for /internal/services package
-   Mocks for services package
-   Mocks for logger package
-   Mocks for pgxs package
-   Repo interface
-   Logger interface
-   Api interface
-   Added shell scripts for testing and check coverage

### Changes

-   CheckIsAdmin middleware was cleared of unnecessary code

## 23 Feb 2024

### Added

-   Tests for /cmd package

## 22 Feb 2024

### Added

-   Tests for /pkg/logger package
-   Tests for /pkg/errors package

### Changes

-   Change unknown error code
-   Change [User Credentials] eror message

## 21 Feb 2024

### Added

-   Tests for /pkg/pgxs package

### Changes

-   Added main/test values for postgres field in config

## 20 Feb 2024

### Added

-   Tests for /pkg/httplib package
-   Tests for /pkg/sigx package
-   Tests for /pkg/utils package
-   Tests for /pkg/cfg package

## Released [v0.1.0]

## 19 Feb 2024

### Added

-   New error ( EventIsDone = 902)
-   Documentation for repo package methods
-   Documentation for common package methods
-   Documentation for fights / events / fighter structs in model package

### Changed

-   .gitignore changed
-   GetUndoneFights method name has been changed to GetUndoneFightsCount

[Unreleased](https://github.com/DoRightt/fb-app/v0.3.0...main)

[v0.3.2]: https://github.com/DoRightt/fb-app/compare/v0.3.0...v0.3.2
[v0.3.0]: https://github.com/DoRightt/fb-app/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/DoRightt/fb-app/compare/v0.1.0...v0.2.0
[v0.1.0]: https://github.com/DoRightt/fb-app/tree/v0.1.0
