# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.4] - 2024-01-16

### Fixed

- Slack API pagination bug

## [0.1.3] - 2024-01-16

### Fixed

- Slack token regex

## [0.1.2] - 2024-01-16

### Changed

- Removed middleware from get user endpoint

## [0.1.1] - 2024-01-16

### Fixed

- Filter by in `getPermissions` 

## [0.1.0] - 2024-01-12

### Added

- Config file for Slack
- Api handlers
- Docs
- CI

### Changed

- Improved logging, wrapped errors
- Some functions were renamed (particularly under internal/slack)
- Slack package methods now use pqueue

### Fixed

- User verifying module

### Deprecated

- `/get_input` and `/users/unverified` endpoints were deleted
- `helpers` package was removed

[0.1.4]: https://github.com/acs-dl/slack-module-svc/compare/v0.1.3...v0.1.4
[0.1.3]: https://github.com/acs-dl/slack-module-svc/compare/v0.1.2...v0.1.3
[0.1.2]: https://github.com/acs-dl/slack-module-svc/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/acs-dl/slack-module-svc/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/acs-dl/slack-module-svc/releases/tag/v0.1.0
