# Add Missing API Parameters to go-argocd-client

## Goal
Add all missing query/body parameters from the ArgoCD API proto definitions to every method in the go-argocd-client, introducing option structs where needed and validator structs for input validation.

## Design Decisions

1. **Option structs pattern**: For methods that currently take no parameters or only positional args but need additional optional params, introduce `*XxxOptions` structs (pointer to allow nil = no options). This follows the existing pattern used by `ApplicationDeleteOptions`, `SyncOptions`, `WatchOptions`, etc.
2. **Backward compatibility**: All new parameters are optional. Methods that currently accept positional args will gain an additional `opts *XxxOptions` parameter. Methods that already accept an options struct will have fields added to it.
3. **Validators**: Add `Validate() error` methods on option structs where required fields exist (e.g., `PodLogsOptions` requires `Name` and `ResourceName`).
4. **Query params vs body**: Follow the proto HTTP annotations - GET requests use query params, POST/PUT/PATCH use body.

## Task List

### 1. Application Service (`api/application.go`)

#### 1.1 `List()` -> `List(opts *ApplicationListOptions)`
Add `ApplicationListOptions` struct:
- `Refresh string`
- `Projects []string`
- `ResourceVersion string`
- `Selector string`
- `Repo string`
- `AppNamespace string`
- `Project []string` (legacy)

#### 1.2 `Get(name)` -> `Get(name string, opts *ApplicationGetOptions)`
Add `ApplicationGetOptions` struct:
- `Refresh string`
- `Projects []string`
- `ResourceVersion string`
- `Selector string`
- `Repo string`
- `AppNamespace string`
- `Project []string` (legacy)

#### 1.3 `Create(app)` -> `Create(app *ApplicationModel, opts *ApplicationCreateOptions)`
Add `ApplicationCreateOptions` struct:
- `Upsert *bool`
- `Validate *bool`

#### 1.4 `Update(app)` -> `Update(app *ApplicationModel, opts *ApplicationUpdateOptions)`
Add `ApplicationUpdateOptions` struct:
- `Validate *bool`
- `Project string`

#### 1.5 `Delete(name, opts)` - Add fields to `ApplicationDeleteOptions`:
- `AppNamespace string`
- `Project string`

#### 1.6 `Sync(name, opts)` - Add fields to `SyncOptions`:
- `Manifests []string`
- `Infos []Info`
- `RetryStrategy *RetryStrategy`
- `SyncOptions []string` (the items list)
- `AppNamespace string`
- `Project string`
- `SourcePositions []int64`
- `Revisions []string`

#### 1.7 `Rollback(name, opts)` - Add fields to `RollbackOptions`:
- `DryRun bool`
- `Prune bool`
- `AppNamespace string`
- `Project string`

#### 1.8 `Patch(name, patch, patchType)` -> `Patch(name string, patch any, patchType string, opts *ApplicationPatchOptions)`
Add `ApplicationPatchOptions`:
- `AppNamespace string`
- `Project string`

#### 1.9 `TerminateOperation(name)` -> `TerminateOperation(name string, opts *TerminateOperationOptions)`
Add `TerminateOperationOptions`:
- `AppNamespace string`
- `Project string`

#### 1.10 Add `AppNamespace` and `Project` to these existing structs:
- `ApplicationResourceRequest`
- `ApplicationResourcePatchRequest`
- `ApplicationResourceDeleteRequest` (also add `Orphan *bool`)

#### 1.11 `RunResourceAction` - Add `AppNamespace`, `Project` to `ApplicationResourceActionRequest`

#### 1.12 `GetManifests` - Add fields to `ApplicationManifestQuery`:
- `Project string`
- `SourcePositions []int64`
- `Revisions []string`
- `NoCache *bool`

#### 1.13 `ResourceTree(name)` -> `ResourceTree(name string, opts *ResourcesQuery)`
Add `ResourcesQuery` struct:
- `Namespace string`
- `Name string`
- `Version string`
- `Group string`
- `Kind string`
- `AppNamespace string`
- `Project string`

#### 1.14 `ManagedResources(name)` -> `ManagedResources(name string, opts *ResourcesQuery)`
Reuse the same `ResourcesQuery` struct.

#### 1.15 `RevisionMetadata(name, revision)` -> `RevisionMetadata(name, revision string, opts *RevisionMetadataOptions)`
Add `RevisionMetadataOptions`:
- `AppNamespace string`
- `Project string`
- `SourceIndex *int32`
- `VersionId *int32`

#### 1.16 `GetSyncWindows(name)` -> `GetSyncWindows(name string, opts *SyncWindowsOptions)`
Add `SyncWindowsOptions`:
- `AppNamespace string`
- `Project string`

#### 1.17 `ListResourceEvents` - Add fields to `ApplicationResourceEventsQuery`:
- `ResourceNamespace string`
- `AppNamespace string`
- `Project string`

#### 1.18 `ListLinks(name)` -> `ListLinks(name string, opts *ListLinksOptions)`
Add `ListLinksOptions`:
- `Namespace string`
- `Project string`

#### 1.19 `ListResourceLinks` - add `AppNamespace`, `Project` (uses `ApplicationResourceRequest`)

#### 1.20 `Watch` - Add fields to `WatchOptions`:
- `Projects []string`
- `ResourceVersion string`
- `Selector string`
- `Repo string`
- `Project []string`
- `Refresh string`

#### 1.21 `WatchResourceTree(ctx, name)` -> `WatchResourceTree(ctx, name string, opts *ResourcesQuery)`
Reuse `ResourcesQuery`.

#### 1.22 `PodLogs` - Add fields to `PodLogsOptions`:
- `PodName string`
- `UntilTime string`
- `Filter string`
- `AppNamespace string`
- `Project string`
- `MatchCase *bool`

### 2. Cluster Service (`api/cluster.go`)

#### 2.1 `List()` -> `List(opts *ClusterQueryOptions)`
Add `ClusterQueryOptions`:
- `Server string`
- `Name string`

#### 2.2 `Get(server)` -> `Get(server string, opts *ClusterQueryOptions)`
Reuse `ClusterQueryOptions` (for `name` field when server is empty).

#### 2.3 `Create(cluster)` -> `Create(cluster *ClusterModel, opts *ClusterCreateOptions)`
Add `ClusterCreateOptions`:
- `Upsert bool`

#### 2.4 `Update(cluster)` -> `Update(cluster *ClusterModel, opts *ClusterUpdateOptions)`
Add `ClusterUpdateOptions`:
- `UpdatedFields []string`

#### 2.5 `Delete(server)` -> `Delete(server string, opts *ClusterQueryOptions)`
Reuse `ClusterQueryOptions`.

#### 2.6 `RotateAuth(server)` -> `RotateAuth(server string, opts *ClusterQueryOptions)`
#### 2.7 `InvalidateCache(server)` -> `InvalidateCache(server string, opts *ClusterQueryOptions)`

### 3. Project Service (`api/project.go`)

#### 3.1 `Create(project)` -> `Create(project *ProjectModel, opts *ProjectCreateOptions)`
Add `ProjectCreateOptions`:
- `Upsert bool`

#### 3.2 `DeleteToken` - Add `ID string` parameter:
`DeleteToken(project, role string, iat int64, id string)`

### 4. Repository Service (`api/repository.go`)

#### 4.1 `List()` -> `List(opts *RepositoryQueryOptions)`
Add `RepositoryQueryOptions`:
- `Repo string`
- `ForceRefresh bool`
- `AppProject string`

#### 4.2 `Get(repo)` -> `Get(repo string, opts *RepositoryQueryOptions)`
#### 4.3 `Delete(repo)` -> `Delete(repo string, opts *RepositoryQueryOptions)`
#### 4.4 `GetHelmCharts(repo)` -> `GetHelmCharts(repo string, opts *RepositoryQueryOptions)`
#### 4.5 `ListRefs(repo)` -> `ListRefs(repo string, opts *RepositoryQueryOptions)`
#### 4.6 `ListOCITags(repo)` -> `ListOCITags(repo string, opts *RepositoryQueryOptions)`

#### 4.7 `Create(repo)` -> `Create(repo *RepositoryModel, opts *RepositoryCreateOptions)`
Add `RepositoryCreateOptions`:
- `Upsert bool`
- `CredsOnly bool`

#### 4.8 `ListApps(repo)` -> `ListApps(repo string, opts *RepoAppsQueryOptions)`
Add `RepoAppsQueryOptions`:
- `Revision string`
- `AppName string`
- `AppProject string`

#### 4.9 Add fields to `RepoAppDetailsQuery`:
- `SourceIndex *int32`
- `VersionId *int32`

#### 4.10 Add missing fields to `RepoAccessQuery`:
- `GitHubAppPrivateKey string`
- `GitHubAppID int64`
- `GitHubAppInstallationID int64`
- `GitHubAppEnterpriseBaseUrl string`
- `GCPServiceAccountKey string`
- `BearerToken string`
- `InsecureOCIForceHttp bool`
- `AzureServicePrincipalClientId string`
- `AzureServicePrincipalClientSecret string`
- `AzureServicePrincipalTenantId string`
- `AzureActiveDirectoryEndpoint string`

#### 4.11 Write variants - same options as read counterparts:
- `GetWrite(repo, opts)`, `DeleteWriteRepository(repo, opts)`, `ListWriteRepositories(opts)`
- `CreateWriteRepository(repo, opts)` with `RepositoryCreateOptions`
- `UpdateWriteRepository(repo, opts)` (no extra params needed)
- `ValidateWriteAccess` - already uses `RepoAccessQuery`

### 5. ApplicationSet Service (`api/applicationset.go`)

#### 5.1 `List()` -> `List(opts *ApplicationSetListOptions)`
Add `ApplicationSetListOptions`:
- `Projects []string`
- `Selector string`
- `AppsetNamespace string`

#### 5.2 `Get(name)` -> `Get(name string, opts *ApplicationSetGetOptions)`
Add `ApplicationSetGetOptions`:
- `AppsetNamespace string`

#### 5.3 `Create(appSet)` -> `Create(appSet *ApplicationSetModel, opts *ApplicationSetCreateOptions)`
Add `ApplicationSetCreateOptions`:
- `Upsert bool`
- `DryRun bool`

#### 5.4 `Delete(name)` -> `Delete(name string, opts *ApplicationSetDeleteOptions)`
Add `ApplicationSetDeleteOptions`:
- `AppsetNamespace string`

#### 5.5 `ResourceTree(name)` -> `ResourceTree(name string, opts *ApplicationSetTreeOptions)`
Add `ApplicationSetTreeOptions`:
- `AppsetNamespace string`

#### 5.6 `ListResourceEvents(name)` -> `ListResourceEvents(name string, opts *ApplicationSetGetOptions)`

#### 5.7 `Watch(ctx)` -> `Watch(ctx context.Context, opts *ApplicationSetWatchOptions)`
Add `ApplicationSetWatchOptions`:
- `Name string`
- `Projects []string`
- `Selector string`
- `AppSetNamespace string`
- `ResourceVersion string`

### 6. Certificate Service (`api/certificate.go`)

#### 6.1 `Create(certs)` -> `Create(certs *CertificateCreateRequest, opts *CertificateCreateOptions)`
Add `CertificateCreateOptions`:
- `Upsert bool`

### 7. GPG Key Service (`api/gpgkey.go`)

#### 7.1 `Create(key)` -> `Create(key *GPGKeyModel, opts *GPGKeyCreateOptions)`
Add `GPGKeyCreateOptions`:
- `Upsert bool`

### 8. RepoCreds Service (`api/repocreds.go`)

#### 8.1 `Create(creds)` -> `Create(creds *RepoCredsModel, opts *RepoCredsCreateOptions)`
Add `RepoCredsCreateOptions`:
- `Upsert bool`

#### 8.2 `CreateWrite(creds)` -> `CreateWrite(creds *RepoCredsModel, opts *RepoCredsCreateOptions)`

### 9. Session Service (`api/session.go`)

#### 9.1 `Create(username, password)` -> `Create(opts *SessionCreateOptions)`
Add `SessionCreateOptions`:
- `Username string`
- `Password string`
- `Token string`
Add `Validate() error` - must have either (username+password) or token.

### 10. Validators

Add `Validate() error` methods on structs where required fields must be checked:
- `PodLogsOptions.Validate()` - Name and ResourceName required
- `ApplicationResourceRequest.Validate()` - Name, ResourceName, Version, Kind required
- `ApplicationResourceDeleteRequest.Validate()` - same + check Force/Orphan not both set
- `SessionCreateOptions.Validate()` - either (username+password) or token required
- `SyncOptions.Validate()` - at least check struct is coherent
- `RollbackOptions.Validate()` - ID must be > 0
- `ResourcesQuery.Validate()` - ApplicationName required (for ResourceTree/ManagedResources/WatchResourceTree)

### 11. Update all test files

Update all existing tests to match new method signatures. Tests should:
- Pass `nil` for new options where existing tests don't need the new params
- Add new tests for each new options struct verifying query params/body are correctly set

## Implementation Order

1. Start with models/option structs (add all new types)
2. Update interfaces
3. Update implementations
4. Update tests
5. Run `go build ./...` and `go test ./...`
