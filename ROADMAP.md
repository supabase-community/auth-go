# auth-go Library Roadmap to Production

**Document Version:** 2.0
**Created:** December 12, 2025
**Updated:** December 12, 2025
**Target:** v2.0.0 Production Release
**Current Version:** v1.4.0 (Pre-release)

## Executive Summary

This roadmap outlines the path from pre-release (v1.4.0) to production-ready (v2.0.0+) for the auth-go library, following Go open source best practices. The plan prioritizes **non-breaking changes first** to deliver value to users sooner, with breaking changes deferred to v2.0.0.

**Strategy:** Ship incremental improvements in v1.x releases before major v2.0.0 refactor
**Estimated Timeline:** 4-5 releases over 6-9 months
**Breaking Changes:** Deferred to Phase 3 (v2.0.0)

---

## Table of Contents

1. [Phase 0: Release Automation](#phase-0-release-automation-v150)
2. [Phase 1: Non-Breaking Improvements](#phase-1-non-breaking-improvements-v150---v170)
3. [Phase 2: Documentation & Community](#phase-2-documentation--community-v180)
4. [Phase 3: Breaking Changes](#phase-3-breaking-changes-v200)
5. [Phase 4: Post-v2 Polish](#phase-4-post-v2-polish-v210)
6. [Open Source Best Practices](#open-source-best-practices)
7. [Release Strategy](#release-strategy)
8. [Community Engagement](#community-engagement)

---

## Phase 0: Release Automation (v1.5.0)

**Goal:** Automate release process and CHANGELOG generation
**Status:** ✅ Complete
**Breaking Changes:** NO
**Duration:** 1-2 weeks

### 0.1 Setup release-please

**Priority:** Critical (Foundation for all future releases)
**Files:** `.github/workflows/release-please.yml`, `CHANGELOG.md`, `CONTRIBUTING.md`

#### Why release-please?

[release-please](https://github.com/googleapis/release-please) is Google's tool for automating releases based on conventional commits. It:
- Automatically generates CHANGELOG.md from commit messages
- Creates GitHub releases and git tags
- Manages semantic version bumps
- Creates release PRs for review before publishing
- Used by major Go projects (googleapis, google-cloud-go)

#### Implementation Steps

1. **Adopt Conventional Commits**

   Commit format: `<type>(<scope>): <description>`

   ```bash
   # Examples:
   feat: add retry logic with exponential backoff
   fix(admin): handle strconv.Atoi errors properly
   docs: add security best practices section
   test: add unit tests for token endpoint

   # Breaking changes (triggers MAJOR version):
   feat!: add context.Context to all methods
   # OR
   feat: add context support

   BREAKING CHANGE: All methods now require context as first parameter
   ```

   **Commit Types:**
   - `feat`: New feature → MINOR bump (1.5.0 → 1.6.0)
   - `fix`: Bug fix → PATCH bump (1.5.0 → 1.5.1)
   - `docs`: Documentation only → No version bump
   - `test`: Tests only → No version bump
   - `chore`: Maintenance → No version bump
   - `refactor`: Code refactoring → PATCH/MINOR depending on scope
   - `perf`: Performance improvement → PATCH bump
   - `!` or `BREAKING CHANGE:` → MAJOR bump (1.x.x → 2.0.0)

2. **Create GitHub Action Workflow**

   `.github/workflows/release-please.yml`:
   ```yaml
   name: release-please

   on:
     push:
       branches:
         - main

   permissions:
     contents: write
     pull-requests: write

   jobs:
     release-please:
       runs-on: ubuntu-latest
       steps:
         - uses: googleapis/release-please-action@v4
           with:
             release-type: go
             package-name: auth-go
             changelog-types: |
               [
                 {"type":"feat","section":"Features","hidden":false},
                 {"type":"fix","section":"Bug Fixes","hidden":false},
                 {"type":"perf","section":"Performance","hidden":false},
                 {"type":"docs","section":"Documentation","hidden":false},
                 {"type":"refactor","section":"Code Refactoring","hidden":false},
                 {"type":"test","section":"Tests","hidden":true},
                 {"type":"chore","section":"Miscellaneous","hidden":true}
               ]
   ```

3. **Initialize CHANGELOG.md** with historical releases

4. **Update CONTRIBUTING.md** with conventional commit guidelines

5. **Test the automation**:
   - Make a PR with conventional commits
   - Merge PR → release-please creates "Release PR"
   - Review Release PR → Merge → automatic GitHub release created

#### How It Works

```
Developer                  release-please Bot              GitHub
─────────                  ──────────────                  ──────

PR with commits
  "feat: add feature"  →   Scans commits since
  "fix: fix bug"           last release
        ↓                          ↓
    Merged to main       Calculates version bump
                         (feat → minor, fix → patch)
                                   ↓
                         Creates "Release PR"
                         - Updated CHANGELOG.md
                         - Bumped version
        ←───────────────────────┘
   Review Release PR
        ↓
   Merge Release PR     →   Creates GitHub Release
                            Creates git tag
                            Publishes to registry
```

### 0.2 Phase 0 Deliverables

- [x] Conventional Commits documented in CONTRIBUTING.md
- [x] release-please workflow configured
- [x] CHANGELOG.md initialized
- [x] CI validates PR titles (optional but recommended)
- [x] First automated release (v1.5.0) published

**Success Criteria:** ✅ One complete release cycle through release-please

---

## Phase 1: Non-Breaking Improvements (v1.5.0 - v1.7.0)

**Goal:** Add production features without breaking existing code
**Status:** Not Started
**Breaking Changes:** NO
**Duration:** 2-3 months

### 1.1 Structured Error Handling (v1.6.0)

**Priority:** High
**Breaking:** NO - Additive only

#### Implementation

Create `types/error.go`:
```go
// AuthError represents a structured error from Supabase Auth API
type AuthError struct {
    StatusCode int
    Message    string
    ErrorCode  string                 // Supabase error code
    Details    map[string]interface{}
}

func (e *AuthError) Error() string {
    if e.ErrorCode != "" {
        return fmt.Sprintf("auth error (status %d, code %s): %s",
            e.StatusCode, e.ErrorCode, e.Message)
    }
    return fmt.Sprintf("auth error (status %d): %s", e.StatusCode, e.Message)
}
```

Create `endpoints/errors.go`:
```go
func handleErrorResponse(resp *http.Response) error {
    body, _ := io.ReadAll(resp.Body)

    // Try parsing as Supabase error
    var apiErr struct {
        Error     string `json:"error"`
        ErrorCode string `json:"error_code"`
    }

    if json.Unmarshal(body, &apiErr) == nil && apiErr.Error != "" {
        return &types.AuthError{
            StatusCode: resp.StatusCode,
            Message:    apiErr.Error,
            ErrorCode:  apiErr.ErrorCode,
        }
    }

    // Fallback to text
    return &types.AuthError{
        StatusCode: resp.StatusCode,
        Message:    string(body),
    }
}
```

**Non-Breaking Strategy:**
- Add new `AuthError` type
- Keep existing `fmt.Errorf` returns for now
- Let users opt-in to structured errors with `errors.As()`
- In v2.0.0, all errors will be `AuthError`

**Usage (v1.6.0):**
```go
resp, err := client.SignUp(req)
if err != nil {
    // Still works (error is still an error)
    log.Printf("Error: %v", err)

    // NEW: Can now inspect structured errors
    var authErr *types.AuthError
    if errors.As(err, &authErr) {
        if authErr.StatusCode == 429 {
            // Handle rate limiting
        }
    }
}
```

### 1.2 Comprehensive Unit Testing (v1.6.0)

**Priority:** Critical
**Target:** 80%+ coverage

#### Implementation

1. **Create test infrastructure** (`endpoints/testing.go`):
```go
type MockHTTPClient struct {
    DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
    if m.DoFunc != nil {
        return m.DoFunc(req)
    }
    return &http.Response{
        StatusCode: 200,
        Body:       io.NopCloser(bytes.NewBufferString("{}")),
    }, nil
}
```

2. **Table-driven tests** for each endpoint
3. **Test coverage** for error handling paths
4. **Benchmark tests** for critical operations

### 1.3 Configurable HTTP Client (v1.6.0)

**Priority:** High
**Breaking:** NO - Use functional options pattern

#### Implementation

```go
// New options pattern (backward compatible)
type ClientOption func(*clientOptions)

type clientOptions struct {
    httpClient  *http.Client
    timeout     time.Duration
    headers     map[string]string
    retryPolicy *RetryPolicy
}

// WithTimeout sets request timeout
func WithTimeout(d time.Duration) ClientOption {
    return func(o *clientOptions) {
        o.timeout = d
    }
}

// WithHTTPClient uses custom HTTP client
func WithHTTPClient(client *http.Client) ClientOption {
    return func(o *clientOptions) {
        o.httpClient = client
    }
}

// Update New() to accept options
func New(projectRef, apiKey string, opts ...ClientOption) Client {
    options := &clientOptions{
        timeout: 10 * time.Second, // default
    }

    for _, opt := range opts {
        opt(options)
    }

    // ... create client with options
}
```

**Usage:**
```go
// Old way still works (backward compatible)
client := auth.New("ref", "key")

// NEW: Configure timeout
client := auth.New("ref", "key",
    auth.WithTimeout(30 * time.Second))

// NEW: Custom HTTP client
httpClient := &http.Client{
    Transport: &http.Transport{
        MaxIdleConns: 100,
    },
}
client := auth.New("ref", "key",
    auth.WithHTTPClient(httpClient))
```

### 1.4 Retry Logic with Exponential Backoff (v1.7.0)

**Priority:** High
**Breaking:** NO - Opt-in via WithRetryPolicy

#### Implementation

```go
type RetryPolicy struct {
    MaxRetries           int
    MinWait              time.Duration
    MaxWait              time.Duration
    RetryableStatusCodes []int
}

func DefaultRetryPolicy() *RetryPolicy {
    return &RetryPolicy{
        MaxRetries: 3,
        MinWait:    100 * time.Millisecond,
        MaxWait:    2 * time.Second,
        RetryableStatusCodes: []int{429, 500, 502, 503, 504},
    }
}

// Opt-in via option
func WithRetryPolicy(policy *RetryPolicy) ClientOption { ... }
```

**Usage:**
```go
// No retry (current behavior, default for v1.x)
client := auth.New("ref", "key")

// NEW: Opt-in to retry
client := auth.New("ref", "key",
    auth.WithRetryPolicy(auth.DefaultRetryPolicy()))
```

### 1.5 Logger Interface (v1.7.0)

**Priority:** Medium
**Breaking:** NO - Opt-in

```go
type Logger interface {
    LogRequest(ctx context.Context, req *http.Request)
    LogResponse(ctx context.Context, resp *http.Response, duration time.Duration)
    LogError(ctx context.Context, err error)
}

func WithLogger(logger Logger) ClientOption { ... }
```

### 1.6 Code Quality Fixes (v1.6.0 - v1.7.0)

**Priority:** Medium
**Breaking:** NO

- Rename `client.client` → `client.httpClient`
- Fix `adminaudit.go:77` - handle strconv.Atoi errors
- Extract provider constants (token.go:88-93)
- Add input validation helpers

### 1.7 SAML ACS Completion (v1.7.0)

**Priority:** High
**Breaking:** NO - Just adding missing implementation

Implement typed request/response for POST /sso/saml/acs (currently unimplemented).

### 1.8 Phase 1 Deliverables

- [ ] `AuthError` type with structured error handling
- [ ] 80%+ unit test coverage
- [ ] Configurable HTTP client (timeout, custom client, headers)
- [ ] Retry logic (opt-in)
- [ ] Logger interface (opt-in)
- [ ] Code quality improvements
- [ ] SAML ACS fully implemented
- [ ] All changes backward compatible
- [ ] All integration tests pass

**Releases:** v1.6.0, v1.7.0 (Multiple minor releases, all non-breaking)

---

## Phase 2: Documentation & Community (v1.8.0)

**Goal:** Improve developer experience and community engagement
**Status:** Not Started
**Breaking Changes:** NO
**Duration:** 1-2 months

### 2.1 Godoc Examples

Create example functions:
```go
func ExampleClient_SignInWithEmailPassword() { ... }
func ExampleClient_SignUp() { ... }
func ExampleClient_EnrollFactor() { ... }
```

### 2.2 Examples Directory

```
examples/
├── basic-auth/main.go
├── mfa-totp/main.go
├── admin-user-mgmt/main.go
├── oauth-flow/main.go
├── custom-config/main.go
└── README.md
```

### 2.3 CONTRIBUTING.md

Comprehensive contribution guide:
- Development setup
- Running tests
- Commit message format
- PR process

### 2.4 Issue/PR Templates

- `.github/ISSUE_TEMPLATE/bug_report.md`
- `.github/ISSUE_TEMPLATE/feature_request.md`
- `.github/pull_request_template.md`

### 2.5 Security Documentation

Add to README:
- Token storage best practices
- Service role key warnings
- Error logging cautions
- HTTPS requirements

### 2.6 Phase 2 Deliverables

- [ ] 10+ runnable godoc examples
- [ ] 5+ examples in examples/ directory
- [ ] Comprehensive CONTRIBUTING.md
- [ ] Issue and PR templates
- [ ] Security best practices documented
- [ ] README updated with new features

**Release:** v1.8.0 (Documentation improvements)

---

## Phase 3: Breaking Changes (v2.0.0)

**Goal:** Add context.Context support and finalize production-ready API
**Status:** Not Started
**Breaking Changes:** YES
**Duration:** 1-2 months

### Why Wait Until v2.0.0?

- Users have had several v1.x releases to adopt new features
- Community feedback incorporated
- Breaking changes batched together in single major version
- Clear migration path

### 3.1 Context Support (BREAKING)

**Impact:** All 49+ methods change signatures

#### Changes

```go
// v1.x
func (c *Client) SignInWithEmailPassword(email, password string) (*AuthenticatedResponse, error)

// v2.0.0
func (c *Client) SignInWithEmailPassword(ctx context.Context, email, password string) (*AuthenticatedResponse, error)
```

Update:
- `endpoints/request.go`: `newRequest()` → `newRequestWithContext()`
- All 49+ endpoint methods
- All integration tests
- README examples

### 3.2 Make Structured Errors Default (BREAKING)

**Impact:** Error types change

```go
// v1.x: errors are generic
_, err := client.SignUp(req)
// err is fmt.Errorf string

// v2.0.0: errors are always AuthError
_, err := client.SignUp(ctx, req)
// err is *AuthError (can still use as regular error)
```

### 3.3 Make Retry Default (BREAKING)

**Impact:** Default behavior changes

```go
// v1.x: No retry by default
client := auth.New("ref", "key")

// v2.0.0: Retry enabled by default
client := auth.New("ref", "key")
// uses DefaultRetryPolicy()

// Disable if needed:
client := auth.New("ref", "key",
    auth.WithRetryPolicy(nil))
```

### 3.4 Minimum Go Version → 1.21

**Impact:** go.mod changes

Update to Go 1.21 for:
- Better error handling
- Standard library improvements
- Security updates

### 3.5 Module Path (if needed)

**Impact:** Import path changes

If breaking changes require v2 module path:
```
github.com/supabase-community/auth-go/v2
```

But for this library, prefer same path with v2.x.x tags.

### 3.6 Phase 3 Deliverables

- [ ] All methods accept `context.Context`
- [ ] All errors are `*AuthError`
- [ ] Retry enabled by default
- [ ] Minimum Go 1.21
- [ ] Migration guide (v1 → v2)
- [ ] CHANGELOG with breaking changes section
- [ ] Deprecation notices in v1.8.x
- [ ] All tests updated and passing

**Release:** v2.0.0 (MAJOR VERSION - Breaking Changes)

### 3.7 Migration Guide

Provide in `MIGRATION.md`:

```markdown
# Migrating from v1.x to v2.0

## Context Support

**Before (v1.x):**
```go
resp, err := client.SignInWithEmailPassword("user@example.com", "password")
```

**After (v2.0):**
```go
ctx := context.Background()
resp, err := client.SignInWithEmailPassword(ctx, "user@example.com", "password")
```

## Automated Migration

```bash
# Find and replace (example - customize per method)
find . -name "*.go" -exec sed -i '' 's/client.SignInWithEmailPassword(\([^,]*\), \([^)]*\))/client.SignInWithEmailPassword(context.Background(), \1, \2)/g' {} \;
```

## Error Handling

**Before (v1.x):**
```go
_, err := client.SignUp(req)
if err != nil {
    log.Printf("Error: %v", err)
}
```

**After (v2.0):**
```go
_, err := client.SignUp(ctx, req)
if err != nil {
    var authErr *types.AuthError
    if errors.As(err, &authErr) {
        // Access structured error
        log.Printf("Status: %d, Message: %s", authErr.StatusCode, authErr.Message)
    }
}
```
```

---

## Phase 4: Post-v2 Polish (v2.1.0+)

**Goal:** Iterate on v2.0 based on feedback
**Status:** Not Started
**Breaking Changes:** NO
**Duration:** Ongoing

### 4.1 Potential Additions

- Pagination helpers
- Batch operations
- Convenience methods (e.g., `SignUpAndConfirm`)
- Performance optimizations
- Additional observability hooks

### 4.2 Community Driven

Features based on:
- User feedback on v2.0
- GitHub issues
- Real-world usage patterns

---

## Open Source Best Practices

### Go-Specific

1. **API Design**
   - Accept interfaces, return structs
   - Context as first parameter
   - Options pattern for configuration
   - Functional options for flexibility

2. **Error Handling**
   - Custom error types
   - Error wrapping with `%w`
   - Sentinel errors as variables
   - Document error conditions

3. **Testing**
   - Table-driven tests
   - Subtests with `t.Run()`
   - Mocking interfaces
   - Integration tests
   - Benchmarks

4. **Documentation**
   - Godoc for all exports
   - Runnable examples
   - Clear README
   - Migration guides

5. **Versioning**
   - Semantic versioning
   - CHANGELOG.md
   - Clear deprecation policy
   - Migration guides for breaking changes

### Community

1. **Process**
   - CONTRIBUTING.md
   - Issue templates
   - PR templates
   - Code of Conduct

2. **Quality**
   - CI/CD
   - Code coverage (80%+)
   - Linting (golangci-lint)
   - Security scanning
   - Dependency updates

3. **Communication**
   - GitHub Issues for bugs
   - GitHub Discussions for Q&A
   - Timely response (48 hours)
   - Release notes

---

## Release Strategy

### Version Numbering

Following [Semantic Versioning](https://semver.org/):

- **MAJOR** (v1 → v2): Breaking changes
- **MINOR** (v1.5 → v1.6): New features, backward compatible
- **PATCH** (v1.5.0 → v1.5.1): Bug fixes

### Release Process (with release-please)

1. **Development**
   ```bash
   # Make changes
   git commit -m "feat: add new feature"
   git push
   # Create PR, merge
   ```

2. **Automatic Release PR**
   - release-please creates PR with:
     - Updated CHANGELOG.md
     - Version bump
   - PR stays open, accumulates changes

3. **Publishing**
   ```bash
   # When ready to release:
   # 1. Review Release PR
   # 2. Merge Release PR
   # → release-please automatically:
   #    - Creates GitHub release
   #    - Creates git tag
   #    - Closes Release PR
   ```

### Support Policy

- **Latest major (v2.x)**: Full support
- **Previous major (v1.x)**: Security fixes for 6 months after v2.0.0
- **Older versions**: No support

---

## Community Engagement

### Channels

1. **GitHub Issues**: Bug reports, feature requests
2. **GitHub Discussions**: Q&A, ideas
3. **Supabase Discord**: Community support
4. **Social Media**: Release announcements

### Maintainer Responsibilities

- Respond to issues within 48 hours
- Review PRs within 1 week
- Security fixes within 24 hours
- Regular releases (monthly for minor, as-needed for patches)
- Keep documentation updated

### Contributor Recognition

- CONTRIBUTORS.md file
- Credit in release notes
- GitHub badges

---

## Timeline Summary

| Phase | Release | Duration | Breaking | Key Features |
|-------|---------|----------|----------|--------------|
| Phase 0 | v1.5.0 | 1-2 weeks | NO | Release automation, CHANGELOG |
| Phase 1 | v1.6.0-v1.7.0 | 2-3 months | NO | Errors, tests, config, retry, logging |
| Phase 2 | v1.8.0 | 1-2 months | NO | Docs, examples, community |
| Phase 3 | v2.0.0 | 1-2 months | **YES** | Context support, default retry |
| Phase 4 | v2.1.0+ | Ongoing | NO | Polish, community features |

**Total Time to v2.0.0:** ~6-9 months
**Total Time to Production-Ready:** v2.0.0

---

## Next Steps

1. **Start with Phase 0** (release-please setup)
   - Create `.github/workflows/release-please.yml`
   - Initialize CHANGELOG.md
   - Update CONTRIBUTING.md
   - Test with first automated release

2. **Create GitHub issues** for Phase 1 tasks
   - One issue per major feature
   - Label by version (v1.6.0, v1.7.0, v2.0.0)
   - Create milestones

3. **Set up project board**
   - Track progress across phases
   - Organize by priority

4. **Seek contributors**
   - Label "good first issues"
   - Document development setup
   - Welcome community PRs

---

## References

- [Conventional Commits](https://www.conventionalcommits.org/)
- [release-please](https://github.com/googleapis/release-please)
- [Semantic Versioning](https://semver.org/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Keep a Changelog](https://keepachangelog.com/)

---

**Maintained By:** auth-go maintainers
**Last Updated:** December 12, 2025
**Next Review:** After each phase completion
