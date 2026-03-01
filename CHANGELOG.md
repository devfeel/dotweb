# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.8.0] - 2026-03-01

### Added
- GitHub Actions CI workflow
- More code examples (JWT, file-upload, CORS)
- ExistsRouter method to Router interface

### Changed
- Upgrade Go version to 1.21
- Upgrade dependencies
- Improve README with badges and quick start guide

### Fixed
- Issue #245: Return 404 instead of 301 for trailing slash in parameterized routes
- Issue #250: Support AutoOPTIONS for route groups

### Removed
- Deprecated Travis CI configuration

## [1.7.21] - 2023-04-15

### Added
- SessionManager.RemoveSessionState
- HttpContext.DestroySession()

## [1.7.20] - 2022-08-11

### Fixed
- Delete minor unreachable code

## [1.7.19] - 2021-04-20

### Added
- SetReadTimeout, SetReadHeaderTimeout, SetIdleTimeout, SetWriteTimeout in HttpServer

### Fixed
- deepcopy middleware issue

## [1.7.18] - 2021-04-20

### Fixed
- GetRandString returning same result
