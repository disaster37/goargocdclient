# Plan: Achieve 100% Test Coverage

## Current State
- Root package: **94.6%** (only `performLogin` at 77.8%)
- API package: **64.2%** (many methods at 0-71%)
- Overall: **65.5%**

## Coverage Gap Analysis

### Category 1: Completely untested methods (0%)
| File | Method |
|------|--------|
| `application.go` | `PatchResource`, `DeleteResource`, `RunResourceAction`, `ListResourceLinks` |
| `application.go` | `Watch`, `WatchResourceTree`, `PodLogs`, `readSSE` |
| `applicationset.go` | `Watch` |
| `repository.go` | `ListWriteRepositories` |

### Category 2: Missing API error response path (~71% methods)
Most methods test success but not the `resp.IsError()` → `parseError(resp)` branch. Affects nearly every method across all service files.

### Category 3: Missing network/transport error path (~66-85%)
The `if err != nil` branch from resty calls is untested for most methods.

### Category 4: Missing option/branch coverage
| Method | Missing branch |
|--------|---------------|
| `Application.Delete` | `Cascade` and `PropagationPolicy` query params (currently 40%) |
| `Application.Patch` | Empty `patchType` path (currently 81.8%) |
| `Application.GetManifests` | With non-nil `opts` query params (currently 70%) |
| `Application.ListResourceEvents` | With non-nil `opts` query params (currently 70%) |
| `Application.DeleteResource` | With `Force` option |
| `Application.PodLogs` | All optional params (container, tailLines, sinceSeconds, sinceTime, previous) |
| `parseError` | Empty `apiErr.Message` fallback to `resp.Status()` (currently 80%) |
| `performLogin` | Network error + invalid JSON response (currently 77.8%) |

## Implementation Tasks

### Task 1: Add a reusable helper for network error tests
Add to `testutil_test.go`:
- `newFailingClient()` — returns a resty client pointed at a closed server to trigger transport errors.

### Task 2: Add error-path tests for all API service methods
For each service method that currently only has a `_Success` test, add a corresponding `_Error` test that returns a 500 JSON error response. This covers the `resp.IsError()` branch.

**Files affected**: `account_test.go`, `application_test.go`, `applicationset_test.go`, `certificate_test.go`, `cluster_test.go`, `gpgkey_test.go`, `notification_test.go`, `project_test.go`, `repocreds_test.go`, `repository_test.go`, `session_test.go`, `settings_test.go`, `version_test.go`

### Task 3: Add network error tests for all service methods
For each service, add one test that uses the failing client to trigger a transport-level error. This covers the `if err != nil` branch. One test per service (not per method) is sufficient since the pattern is identical.

### Task 4: Add tests for completely untested methods
- `application_test.go`: `PatchResource`, `DeleteResource`, `RunResourceAction`, `ListResourceLinks` (success + error)
- `repository_test.go`: `ListWriteRepositories` (success + error)

### Task 5: Add SSE/streaming tests
- `application_test.go`: `Watch`, `WatchResourceTree`, `PodLogs` — use httptest server that writes SSE `data: {...}\n\n` lines, read from channel with context timeout
- `applicationset_test.go`: `Watch` — same pattern
- Also cover `readSSE` generic function via these tests
- Test error paths: HTTP 400+ response, nil opts for PodLogs

### Task 6: Add branch coverage tests
- `Application.Delete` with `Cascade` and `PropagationPolicy` options
- `Application.Patch` with empty `patchType`
- `Application.GetManifests` with non-nil `opts`
- `Application.ListResourceEvents` with non-nil `opts`
- `Application.DeleteResource` with `Force` option
- `Application.PodLogs` with all optional params (container, tailLines, sinceSeconds, sinceTime, previous)
- `parseError` with empty message body (triggers `resp.Status()` fallback)
- `performLogin` with invalid JSON response body

### Task 7: Fix client_test.go `performLogin` coverage
- Add test for `performLogin` where the server returns non-JSON body (covers JSON unmarshal failure)

## Test Conventions (Go best practices)
- Use `t.Fatal(err)` for unexpected errors that should stop the test
- Use `t.Error(err)` / `t.Errorf()` for assertion failures
- Use `defer server.Close()` for httptest cleanup
- Keep test names descriptive: `TestXxx_Success`, `TestXxx_Error`, `TestXxx_WithOpts`
- Use the existing `newTestClient` / `jsonResponse` helpers
- No external test dependencies (no testify, gomock, etc.)
- Tests in same package (white-box) — consistent with existing pattern
- Each test is self-contained with its own httptest server

## Validation
```bash
go test ./... -coverprofile=coverage.out -count=1
go tool cover -func=coverage.out | tail -1
# Expected: total: (statements) 100.0%
```
