# FB-Server Changelog

All notable changes to this project will be documented in this file.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased] v0.1.1

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

[Unreleased](https://github.com/DoRightt/fb-server/compare/v0.1.0...main)

[v0.1.0]: https://github.com/DoRightt/fb-server/tree/v0.1.0
