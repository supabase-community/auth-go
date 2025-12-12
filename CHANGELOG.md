# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.4.0](https://github.com/supabase-community/auth-go/compare/v1.3.2...v1.4.0) - 2025-06-19

### Features

- Add SignInWithIdToken method for OAuth social login (#12)
  - Supports github, apple, kakao, and keycloak providers
  - Enables ID token grant type for OAuth flows

### Bug Fixes

- Update github.com/golang-jwt/jwt/v4 from 4.4.2 to 4.5.1 (#9)
  - Addresses security vulnerability CVE-2022-29526

### Documentation

- Add pagination support for GET /admin/users (#11)
  - Added Page and PerPage fields to AdminListUsersRequest

## [1.3.2](https://github.com/supabase-community/auth-go/compare/v1.3.1...v1.3.2) - 2024-10-09

### Security

- Bump github.com/golang-jwt/jwt/v4 from 4.5.0 to 4.5.1
  - Security update for JWT library

## [1.3.1](https://github.com/supabase-community/auth-go/compare/v1.3.0...v1.3.1) - 2024-10-09

### Bug Fixes

- Minor bug fixes and improvements

## [1.3.0](https://github.com/supabase-community/auth-go/compare/v1.2.1...v1.3.0) - 2024-09-08

### Features

- Add pagination support for AdminListUsers endpoint (#11)
  - Added Page and PerPage parameters to AdminListUsersRequest
  - Enables efficient pagination of user lists

### Changed

- Repository renamed from gotrue-go to auth-go
  - Updated module path to github.com/supabase-community/auth-go

## [1.2.1](https://github.com/supabase-community/auth-go/compare/v1.2.0...v1.2.1) - 2024-07-28

### Bug Fixes

- Minor bug fixes and stability improvements

## [1.2.0](https://github.com/supabase-community/auth-go/compare/v1.1.0...v1.2.0) - 2024-01-26

### Features

- Additional authentication methods and improvements

## [1.1.0](https://github.com/supabase-community/auth-go/compare/v1.0.1...v1.1.0) - 2023-08-24

### Features

- Enhanced authentication capabilities
- Additional endpoint support

## [1.0.1](https://github.com/supabase-community/auth-go/compare/v1.0.0...v1.0.1) - 2022-11-17

### Bug Fixes

- Bug fixes and stability improvements

## [1.0.0](https://github.com/supabase-community/auth-go/compare/v0.2.0...v1.0.0) - 2022-11-11

### Features

- First stable release
- Complete implementation of Supabase Auth API
- Support for authentication, user management, and admin operations
- MFA (Multi-Factor Authentication) support
- SSO/SAML provider management

## [0.2.0](https://github.com/supabase-community/auth-go/compare/v0.1.0...v0.2.0) - 2022-11-08

### Features

- Beta release with expanded functionality
- Additional endpoint coverage

## [0.1.0](https://github.com/supabase-community/auth-go/releases/tag/v0.1.0) - 2022-11-07

### Features

- Initial beta release
- Basic authentication endpoints
- User management
- Admin operations

---

## Notes

### Pre-release Status

**Current Status:** This library is in pre-release and should not be used in production environments.

From v1.5.0 onwards, this project uses [Conventional Commits](https://www.conventionalcommits.org/) and [release-please](https://github.com/googleapis/release-please) for automated releases and changelog generation.

### Version History Note

You may notice a v2.0.0 tag in the repository from 2022. This was an early tag that has been superseded by the v1.x release series. The project is currently working towards a production-ready v2.0.0 release following the roadmap in ROADMAP.md.
