# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Group SetNotFoundHandle** - Support custom 404 handler for router groups (#229)
  - Added `SetNotFoundHandle(handler StandardHandle)` method to Group interface
  - Group-level 404 handler takes priority over app-level handler
  - Example: `apiGroup.SetNotFoundHandle(func(ctx Context) error { ... })`

- **Trailing Slash Redirect Configuration** - Control trailing slash redirect behavior (#245)
  - Added `EnabledRedirectTrailingSlash` configuration option
  - Default is `false` to match net/http behavior
  - Set to `true` to enable automatic 301 redirect for trailing slashes

- **CORS Preflight Support** - Fixed CORS preflight request handling (#250)
  - Modified `DefaultAutoOPTIONSHandler` to include CORS headers
  - Better integration with CORS middleware

### Fixed
- **Security Vulnerabilities** - Upgraded dependencies to fix security issues
  - Upgraded `gopkg.in/yaml.v2` to `v3.0.1` (fixes DoS vulnerability)
  - Upgraded `golang.org/x/net` to `v0.33.0` (fixes XSS and proxy bypass)

### Changed
- **Go Version** - Maintained Go 1.21 compatibility
- **Code Quality** - Improved error handling and nil checks

## [1.0.0] - Previous Release

### Features
- Support for go mod
- Static routing, parameter routing, group routing
- Route support for file/directory services
- HttpModule support
- Middleware support (App, Group, Router levels)
- STRING/JSON/JSONP/HTML output formats
- Built-in Mock capability
- Custom Context support
- Timeout Hook integration
- Global HTTP error handling
- Global logging
- Hijack and WebSocket support
- Built-in Cache support
- Built-in Session support with Redis failover
- Built-in TLS support
- Third-party template engine support
- Modular configuration
- Built-in statistics

---

For more details, see [GitHub Releases](https://github.com/devfeel/dotweb/releases).
