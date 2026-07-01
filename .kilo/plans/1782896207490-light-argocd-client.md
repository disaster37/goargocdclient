# Light ArgoCD Client - Implementation Plan

## Overview
Implement a lightweight ArgoCD REST API client in Go with no Kubernetes dependencies. Supports all 13 ArgoCD services, streaming endpoints, token + basic auth, custom error types, and auto-pagination. Uses `resty.Client` for HTTP.

## Architecture

```
client.go (package goargocdclient)
  └─ Creates *resty.Client (base URL, auth, TLS, timeouts)
  └─ Calls api.New(client) → returns api.API interface

api/ (package api)
  └─ api.go: API interface with 13 service accessors + New(*resty.Client) factory
  └─ models.go: Shared typed Go structs (ObjectMeta, TypeMeta, etc.)
  └─ errors.go: APIError type + helpers (IsNotFound, IsUnauthorized, IsConflict)
  └─ <service>.go: Interface + Standard implementation per service
  └─ <service>_test.go: Unit tests per service
```

**Dependency**: Only `github.com/go-resty/resty/v2` (already in go.mod)

## Implementation Tasks (ordered)

### Phase 1: Core Infrastructure

1. **`api/models.go`** - Shared models
   - `ObjectMeta` (name, namespace, labels, annotations, uid, resourceVersion, creationTimestamp, deletionTimestamp, generation)
   - `TypeMeta` (kind, apiVersion)
   - `ListMeta` (resourceVersion, selfLink)
   - Common response wrappers

2. **`api/errors.go`** - Error handling
   - `APIError{Code int, Message string, Details string}` implementing `error`
   - `parseError(resp *resty.Response) error` - parses ArgoCD JSON error from response
   - `IsNotFound(err)`, `IsUnauthorized(err)`, `IsConflict(err)` helpers

3. **`client.go`** - Main client entry point
   - `Option` type (functional options)
   - `WithToken(token string)`, `WithUsernamePassword(user, pass string)`, `WithInsecure()`, `WithTimeout(d time.Duration)`
   - `New(serverURL string, opts ...Option) (api.API, error)` - creates resty client, applies options, auto-login if user/pass, returns api.API

4. **`api/api.go`** - Update API interface & factory
   - Add all 13 service accessors: `Account()`, `Application()`, `ApplicationSet()`, `Certificate()`, `Cluster()`, `GPGKey()`, `Notification()`, `Project()`, `RepoCreds()`, `Repository()`, `Session()`, `Settings()`, `Version()`
   - `APIStandard` stores all service instances
   - `New(client *resty.Client)` initializes all services

### Phase 2: Core Services

5. **`api/session.go`** - Session service
   - Interface: `Create(username, password string) (*SessionResponse, error)`, `Delete() error`, `GetUserInfo() (*UserInfo, error)`
   - Endpoints: POST `/api/v1/session`, DELETE `/api/v1/session`, GET `/api/v1/session/userinfo`

6. **`api/account.go`** - Account service
   - Interface: `List()`, `Get(name)`, `CanI(resource, action, subresource)`, `UpdatePassword(current, new, name)`, `CreateToken(name, expiresIn, id)`, `DeleteToken(name, id)`
   - Endpoints: `/api/v1/account`, `/api/v1/account/{name}`, `/api/v1/account/can-i/{resource}/{action}/{subresource}`, `/api/v1/account/password`, `/api/v1/account/{name}/token`, `/api/v1/account/{name}/token/{id}`

7. **`api/application.go`** - Application service
   - CRUD: `List()`, `Get(name)`, `Create(app)`, `Update(app)`, `Delete(name, opts)`
   - Operations: `Sync(name, opts)`, `Rollback(name, opts)`, `TerminateOperation(name)`, `Patch(name, patch, patchType)`
   - Resources: `GetResource(opts)`, `PatchResource(opts)`, `DeleteResource(opts)`, `ListResourceActions(opts)`, `RunResourceAction(opts)`
   - Metadata: `GetManifests(name, opts)`, `ResourceTree(name)`, `ManagedResources(name)`, `RevisionMetadata(name, revision)`, `GetSyncWindows(name)`
   - Events/Links: `ListResourceEvents(name, opts)`, `ListLinks(name)`, `ListResourceLinks(opts)`
   - Streaming: `Watch(ctx, opts)`, `WatchResourceTree(ctx, name)`, `PodLogs(ctx, opts)`
   - Endpoints: `/api/v1/applications/*`

8. **`api/cluster.go`** - Cluster service
   - Interface: `List()`, `Get(server)`, `Create(cluster)`, `Update(cluster)`, `Delete(server)`, `RotateAuth(server)`, `InvalidateCache(server)`
   - Endpoints: `/api/v1/clusters/*`

9. **`api/project.go`** - Project service
   - CRUD: `List()`, `Get(name)`, `Create(project)`, `Update(project)`, `Delete(name)`
   - Extended: `GetDetailed(name)`, `GetGlobalProjects(name)`, `CreateToken(project, role, opts)`, `DeleteToken(project, role, iat)`, `ListEvents(name)`, `GetSyncWindowsState(name)`, `ListLinks(name)`
   - Endpoints: `/api/v1/projects/*`

10. **`api/repository.go`** - Repository service
    - CRUD: `List()`, `Get(repo)`, `Create(repo)`, `Update(repo)`, `Delete(repo)`
    - Apps: `ListApps(repo)`, `GetAppDetails(opts)`, `GetHelmCharts(repo)`
    - Refs: `ListRefs(repo)`, `ListOCITags(repo)`
    - Access: `ValidateAccess(opts)`
    - Write: `ListWriteRepositories()`, `GetWrite(repo)`, `CreateWriteRepository(repo)`, `UpdateWriteRepository(repo)`, `DeleteWriteRepository(repo)`, `ValidateWriteAccess(opts)`
    - Endpoints: `/api/v1/repositories/*`, `/api/v1/write-repositories/*`

11. **`api/version.go`** - Version service
    - Interface: `Get() (*VersionInfo, error)`
    - Endpoint: GET `/api/version`

### Phase 3: Additional Services

12. **`api/settings.go`** - Settings service
    - Interface: `Get()`, `GetPlugins()`
    - Endpoints: `/api/v1/settings`, `/api/v1/settings/plugins`

13. **`api/applicationset.go`** - ApplicationSet service
    - CRUD: `List()`, `Get(name)`, `Create(appset)`, `Delete(name)`
    - Extended: `Generate(appset)`, `ResourceTree(name)`, `ListResourceEvents(name)`
    - Streaming: `Watch(ctx)`
    - Endpoints: `/api/v1/applicationsets/*`

14. **`api/certificate.go`** - Certificate service
    - Interface: `List(opts)`, `Create(certs)`, `Delete(opts)`
    - Endpoints: `/api/v1/certificates`

15. **`api/repocreds.go`** - RepoCreds service
    - CRUD: `List()`, `Create(creds)`, `Update(creds)`, `Delete(url)`
    - Write: `ListWrite()`, `CreateWrite(creds)`, `UpdateWrite(creds)`, `DeleteWrite(url)`
    - Endpoints: `/api/v1/repocreds/*`, `/api/v1/write-repocreds/*`

16. **`api/gpgkey.go`** - GPGKey service
    - Interface: `List()`, `Get(keyID)`, `Create(key)`, `Delete(keyID)`
    - Endpoints: `/api/v1/gpgkeys/*`

17. **`api/notification.go`** - Notification service
    - Interface: `ListTriggers()`, `ListServices()`, `ListTemplates()`
    - Endpoints: `/api/v1/notifications/*`

### Phase 4: Testing (100% coverage)

18. **`api/testutil_test.go`** - Test infrastructure
    - `httptest.Server` mock with configurable route handlers
    - `newTestClient(handler http.Handler) *resty.Client` helper
    - JSON response builders for each model type

19. **Unit tests** - One `_test.go` per service + errors + client
    - Test every method (success + error paths)
    - Test error parsing
    - Test streaming (mock SSE responses)
    - Test auth flow in client.go

## Key Implementation Details

### Streaming (Watch/PodLogs)
```go
func (a *ApplicationStandard) Watch(ctx context.Context, opts WatchOptions) (<-chan *ApplicationWatchEvent, error) {
    // resty SetDoNotParseResponse → raw http.Response
    // goroutine reads body line-by-line, parses JSON events
    // sends to channel, closes on ctx.Done() or EOF
}
```

### Error Handling
```go
type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}
// parseError checks resp.StatusCode(), unmarshals error body, returns *APIError
// IsNotFound(err) → errors.As + Code == 404
```

### Auth Flow
```go
// client.go
func New(serverURL string, opts ...Option) (api.API, error) {
    client := resty.New().SetBaseURL(serverURL).SetHeader("Content-Type", "application/json")
    // apply options
    // if username/password: POST /api/v1/session → get token → SetAuthToken(token)
    return api.New(client), nil
}
```

## File Structure
```
client.go, client_test.go
api/api.go, api/models.go, api/errors.go, api/errors_test.go, api/testutil_test.go
api/session.go, api/session_test.go
api/account.go, api/account_test.go
api/application.go, api/application_test.go
api/cluster.go, api/cluster_test.go
api/project.go, api/project_test.go
api/repository.go, api/repository_test.go
api/version.go, api/version_test.go
api/settings.go, api/settings_test.go
api/applicationset.go, api/applicationset_test.go
api/certificate.go, api/certificate_test.go
api/repocreds.go, api/repocreds_test.go
api/gpgkey.go, api/gpgkey_test.go
api/notification.go, api/notification_test.go
```

## Validation
1. `go build ./...` - compiles
2. `go test -v -cover -count=1 ./...` - 100% coverage
3. `go vet ./...` - no issues
