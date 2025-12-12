# Contributing to auth-go

Thank you for your interest in contributing to auth-go! This document provides guidelines and instructions for contributing to this project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Commit Message Format](#commit-message-format)
- [Pull Request Process](#pull-request-process)
- [Coding Guidelines](#coding-guidelines)
- [Testing](#testing)
- [Reporting Bugs](#reporting-bugs)
- [Suggesting Features](#suggesting-features)

## Code of Conduct

Be respectful, inclusive, and constructive in all interactions. We're building this together as a community.

## Getting Started

### Prerequisites

- **Go 1.18 or later** (Go 1.21+ recommended for development)
- **Docker and Docker Compose** (for integration tests)
- **Git**

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR-USERNAME/auth-go.git
   cd auth-go
   ```

3. Add upstream remote:
   ```bash
   git remote add upstream https://github.com/supabase-community/auth-go.git
   ```

## Development Setup

### Install Dependencies

```bash
go mod download
```

### Run Unit Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run Integration Tests

Integration tests require Docker:

```bash
cd integration_test
docker-compose up -d
go test -v
docker-compose down
```

### Code Formatting

Always format your code before committing:

```bash
# Format code
go fmt ./...

# Run goimports (if installed)
goimports -w .

# Run linter
go vet ./...
```

## Commit Message Format

**IMPORTANT:** Starting from v1.5.0, this project follows [Conventional Commits](https://www.conventionalcommits.org/).

### Format

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

### Type

Must be one of the following:

- **feat**: A new feature (triggers MINOR version bump)
- **fix**: A bug fix (triggers PATCH version bump)
- **docs**: Documentation only changes
- **test**: Adding missing tests or correcting existing tests
- **refactor**: Code change that neither fixes a bug nor adds a feature
- **perf**: Code change that improves performance
- **chore**: Changes to build process, dependencies, or auxiliary tools
- **ci**: Changes to CI configuration files and scripts

### Scope

Optional. Can be anything specifying the place of the commit change:
- `admin` - Admin endpoints
- `auth` - Authentication methods
- `types` - Type definitions
- `client` - Client configuration
- `mfa` - Multi-factor authentication
- `saml` - SAML/SSO functionality

### Description

- Use imperative, present tense: "add" not "added" nor "adds"
- Don't capitalize first letter
- No period (.) at the end
- Limit to 72 characters

### Examples

#### Feature
```bash
git commit -m "feat: add retry logic with exponential backoff"
```

```bash
git commit -m "feat(admin): add pagination support for user list

Adds Page and PerPage parameters to AdminListUsersRequest to enable
efficient pagination of large user lists.

Closes #123"
```

#### Bug Fix
```bash
git commit -m "fix(admin): handle strconv.Atoi errors in audit endpoint

Previously errors from strconv.Atoi were silently ignored in
adminaudit.go:77. This now properly handles the error and returns
it to the caller.

Fixes #456"
```

#### Documentation
```bash
git commit -m "docs: add security best practices to README"
```

#### Refactoring
```bash
git commit -m "refactor: rename client.client to client.httpClient

Improves code clarity by using a more descriptive field name."
```

#### Breaking Changes

Add `!` after the type or include `BREAKING CHANGE:` in the footer:

```bash
git commit -m "feat!: add context.Context to all methods

BREAKING CHANGE: All client methods now require context.Context as
the first parameter. This enables proper timeout and cancellation
support.

Migration guide: Pass context.Background() for top-level calls, or
propagate context from your application.

Closes #789"
```

### Why Conventional Commits?

- **Automatic versioning**: Commit type determines version bump (feat=minor, fix=patch, BREAKING CHANGE=major)
- **Automatic CHANGELOG**: CHANGELOG.md is generated from commits
- **Clear history**: Easy to understand what changed and why
- **Consistent**: Enforced format across all contributors

## Pull Request Process

### 1. Create a Feature Branch

```bash
git checkout -b feat/my-new-feature
# or
git checkout -b fix/bug-description
```

### 2. Make Your Changes

- Write clear, concise code
- Add tests for new functionality
- Update documentation as needed
- Follow the coding guidelines below

### 3. Test Your Changes

```bash
# Run all tests
go test ./...

# Run integration tests
cd integration_test
docker-compose up -d
go test -v
docker-compose down
```

### 4. Commit Your Changes

Follow the [Commit Message Format](#commit-message-format):

```bash
git add .
git commit -m "feat: add awesome new feature"
```

### 5. Push to Your Fork

```bash
git push origin feat/my-new-feature
```

### 6. Open a Pull Request

- Go to your fork on GitHub
- Click "New Pull Request"
- Select your feature branch
- Fill out the PR template
- Link related issues

### PR Guidelines

- **Title**: Use conventional commit format (e.g., "feat: add new feature")
- **Description**: Clearly describe what changed and why
- **Tests**: Include tests for new functionality
- **Documentation**: Update README, godoc comments, or examples as needed
- **One concern per PR**: Keep PRs focused on a single feature or fix
- **Breaking changes**: Clearly document breaking changes in PR description

### PR Review Process

1. Maintainers will review within 1 week
2. Address any review comments
3. Once approved, maintainer will merge
4. Your contribution will be included in the next release!

## Coding Guidelines

### Go Style

- Follow [Effective Go](https://golang.org/doc/effective_go)
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for formatting (run `go fmt ./...`)
- Use `go vet` to catch common mistakes

### Godoc Comments

All exported types, functions, and methods must have godoc comments:

```go
// SignInWithEmailPassword authenticates a user with email and password.
// It returns an AuthenticatedResponse containing the user and session tokens.
//
// POST /token
func (c *Client) SignInWithEmailPassword(email, password string) (*types.AuthenticatedResponse, error) {
    // ...
}
```

### Error Handling

- Return errors, don't panic
- Wrap errors with context: `fmt.Errorf("failed to do X: %w", err)`
- Use custom error types when appropriate
- Don't ignore errors (no `_ = err`)

### Naming Conventions

- **Interfaces**: Descriptive names (e.g., `Client`, `Logger`)
- **Structs**: Pascal case (e.g., `AuthError`, `TokenRequest`)
- **Functions**: Camel case (e.g., `signInWithEmail`)
- **Constants**: Pascal case or SCREAMING_SNAKE_CASE
- **Private fields**: camel case (e.g., `httpClient`)
- **Public fields**: Pascal case (e.g., `StatusCode`)

### Code Organization

- Keep functions small and focused
- Group related functionality together
- Use interfaces for abstraction
- Avoid global state

## Testing

### Test Coverage

- Aim for 80%+ code coverage
- Write tests for all new functionality
- Include both success and error cases
- Use table-driven tests for multiple scenarios

### Table-Driven Tests

```go
func TestTokenRequest(t *testing.T) {
    tests := []struct {
        name    string
        input   types.TokenRequest
        want    *types.AuthenticatedResponse
        wantErr bool
    }{
        {
            name: "successful login",
            input: types.TokenRequest{
                GrantType: "password",
                Email:     "user@example.com",
                Password:  "password",
            },
            wantErr: false,
        },
        {
            name: "invalid credentials",
            input: types.TokenRequest{
                GrantType: "password",
                Email:     "wrong@example.com",
                Password:  "wrong",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := client.Token(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Token() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            // Additional assertions...
        })
    }
}
```

### Integration Tests

- Located in `integration_test/` directory
- Use Docker Compose for test environment
- Test against real Supabase Auth instance
- Clean up resources after tests

## Reporting Bugs

### Before Submitting

- Check existing issues to avoid duplicates
- Verify the bug exists in the latest version
- Collect relevant information

### Bug Report Template

Use GitHub Issues with the bug report template. Include:

1. **Description**: Clear description of the bug
2. **Steps to Reproduce**:
   ```
   1. Initialize client with '...'
   2. Call method '...'
   3. See error
   ```
3. **Expected Behavior**: What you expected to happen
4. **Actual Behavior**: What actually happened
5. **Code Snippet**:
   ```go
   client := auth.New("ref", "key")
   // Minimal reproducible example
   ```
6. **Environment**:
   - Go version: `go version`
   - auth-go version: `v1.4.0`
   - Operating System: macOS/Linux/Windows

## Suggesting Features

### Feature Request Template

Use GitHub Issues with the feature request template. Include:

1. **Problem**: What problem does this solve?
2. **Proposed Solution**: How should it work?
3. **Alternatives**: What alternatives have you considered?
4. **Use Case**: Real-world example of usage
5. **Breaking Changes**: Would this break existing code?

### Feature Discussion

- Open a GitHub Discussion for larger features
- Get feedback before implementing
- Consider backward compatibility
- Think about API design

## Questions?

- **GitHub Issues**: For bugs and feature requests
- **GitHub Discussions**: For questions and ideas
- **Discord**: Join the Supabase Discord for community support

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

**Thank you for contributing to auth-go!** 🎉
