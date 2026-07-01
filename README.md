# go-argocd-client

A lightweight Go client for the [ArgoCD](https://argo-cd.readthedocs.io/) API. All sub-clients are interface-based for testability and follow Go best practices.

[![Go Reference](https://pkg.go.dev/badge/disaster37/goargocdclient.svg)](https://pkg.go.dev/disaster37/goargocdclient)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Installation

```bash
go get disaster37/goargocdclient
```

## Quick Start

```go
package main

import (
    "fmt"
    "os"

    goargocdclient "disaster37/goargocdclient"
)

func main() {
    // Authenticate with a token
    client, err := goargocdclient.New("https://argocd.example.com",
        goargocdclient.WithToken("your-auth-token"),
    )
    if err != nil {
        fmt.Printf("failed to create client: %v\n", err)
        os.Exit(1)
    }

    // List all applications
    apps, err := client.Application().List()
    if err != nil {
        fmt.Printf("failed to list applications: %v\n", err)
        os.Exit(1)
    }
    for _, app := range apps.Items {
        fmt.Printf("- %s (%s)\n", app.Name, app.Status.Health.Status)
    }
}
```

## Authentication

Two authentication methods are supported:

**With token (recommended):**
```go
client, err := goargocdclient.New("https://argocd.example.com",
    goargocdclient.WithToken("eyJhbGciOi..."),
)
```

**With username and password:**
```go
client, err := goargocdclient.New("https://argocd.example.com",
    goargocdclient.WithUsernamePassword("admin", "password"),
)
```

### Additional Options
```go
client, err := goargocdclient.New("https://argocd.example.com",
    goargocdclient.WithToken("your-token"),
    goargocdclient.WithTimeout(60 * time.Second), // default: 30s
    goargocdclient.WithInsecure(),                // skip TLS verification
)
```

## API Sub-Clients

The client exposes the following sub-clients, each corresponding to an ArgoCD API endpoint:

| Sub-client | Description |
|---|---|
| `.Application()` | Manage applications (CRUD, sync, rollback, manifests, resources, watch) |
| `.ApplicationSet()` | Manage application sets (CRUD, generate, watch) |
| `.Account()` | List/get accounts, manage tokens, check permissions |
| `.Cluster()` | Manage clusters (CRUD, rotate auth, invalidate cache) |
| `.Project()` | Manage projects (CRUD, tokens, events, links) |
| `.Repository()` | Manage repositories (CRUD, apps, charts, refs, validate) |
| `.RepoCreds()` | Manage repository credentials (CRUD) |
| `.Certificate()` | Manage TLS certificates |
| `.GPGKey()` | Manage GPG keys |
| `.Notification()` | List notification triggers, services, and templates |
| `.Session()` | Create/delete sessions, get user info |
| `.Settings()` | Get ArgoCD settings and plugins |
| `.Version()` | Get ArgoCD version info |

## Examples

### Managing Applications

```go
// Create an application
app, err := client.Application().Create(&api.ApplicationModel{
    ObjectMeta: api.ObjectMeta{
        Name: "my-app",
    },
    Spec: api.ApplicationSpec{
        Project: "default",
        Source: &api.ApplicationSource{
            RepoURL:        "https://github.com/example/repo",
            Path:           "kustomize/overlays/prod",
            TargetRevision: "main",
        },
        Destination: api.ApplicationDestination{
            Server:    "https://kubernetes.default.svc",
            Namespace: "default",
        },
    },
})

// Get an application
app, err := client.Application().Get("my-app")

// Sync an application
err := client.Application().Sync("my-app", &api.SyncOptions{
    Prune: true,
})

// Rollback to a specific deployment
err := client.Application().Rollback("my-app", &api.RollbackOptions{
    ID: 2,
})

// Delete with cascade
cascade := true
err := client.Application().Delete("my-app", &api.ApplicationDeleteOptions{
    Cascade: &cascade,
})

// Patch an application
_, err = client.Application().Patch("my-app", map[string]any{
    "spec": map[string]any{
        "source": map[string]any{
            "targetRevision": "release-v1.2",
        },
    },
}, "application/merge-patch+json")

// Get application manifests
manifests, err := client.Application().GetManifests("my-app", nil)

// Get resource tree
tree, err := client.Application().ResourceTree("my-app")

// Check the sync status
app, err := client.Application().Get("my-app")
if app.Status.Sync.Status == api.SyncStatusCodeSynced {
    fmt.Println("Application is synced")
}
```

### Watching Application Changes

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

events, err := client.Application().Watch(ctx, nil)
if err != nil {
    panic(err)
}
for event := range events {
    fmt.Printf("[%s] %s\n", event.Type, event.Application.Name)
}
```

### Streaming Pod Logs

```go
ctx := context.Background()

logs, err := client.Application().PodLogs(ctx, &api.PodLogsOptions{
    Name:         "my-app",
    Namespace:    "default",
    ResourceName: "my-app-deployment-abc123",
    Kind:         "Pod",
    Follow:       true,
    TailLines:    100,
})
if err != nil {
    panic(err)
}
for entry := range logs {
    fmt.Printf("[%s] %s\n", entry.PodName, entry.Content)
}
```

### Managing Clusters

```go
cluster, err := client.Cluster().Create(&api.ClusterModel{
    Server: "https://kubernetes.example.com",
    Name:   "production-cluster",
    Config: api.ClusterConfig{
        TLSClientConfig: api.TLSClientConfig{
            Insecure: false,
            CAData:   "base64-encoded-ca-cert",
        },
    },
    Project: "default",
})

// Rotate auth credentials
err := client.Cluster().RotateAuth(cluster.Server)

// Invalidate cache to force refresh
err := client.Cluster().InvalidateCache(cluster.Server)
```

### Managing Repositories

```go
repo, err := client.Repository().Create(&api.RepositoryModel{
    Repo:     "https://github.com/example/repo",
    Username: "git-user",
    Password: "git-token",
})

charts, err := client.Repository().GetHelmCharts(repo.Repo)

refs, err := client.Repository().ListRefs(repo.Repo)
fmt.Printf("Branches: %v\n", refs.Branches)
fmt.Printf("Tags: %v\n", refs.Tags)
```

### Managing Projects

```go
project, err := client.Project().Create(&api.ProjectModel{
    ObjectMeta: api.ObjectMeta{
        Name: "my-project",
    },
    Spec: api.ProjectSpec{
        Description: "My team project",
        SourceRepos: []string{"*"},
        Destinations: []api.ApplicationDestination{
            {
                Server:    "https://kubernetes.default.svc",
                Namespace: "*",
            },
        },
    },
})

// Create a JWT token for a project role
tokenResp, err := client.Project().CreateToken("my-project", "read-only", &api.ProjectTokenCreateOptions{
    ID:        "ci-token",
    ExpiresIn: 86400, // 24 hours
})
fmt.Printf("Token: %s\n", tokenResp.Token)
```

### Managing ApplicationSets

```go
appSet, err := client.ApplicationSet().Create(&api.ApplicationSetModel{
    ObjectMeta: api.ObjectMeta{
        Name: "production-apps",
    },
    Spec: api.ApplicationSetSpec{
        Generators: []api.ApplicationSetGenerator{
            {
                List: &api.ListGenerator{
                    Elements: []map[string]string{
                        {"cluster": "us-east-1", "url": "https://1.1.1.1"},
                        {"cluster": "us-west-1", "url": "https://2.2.2.2"},
                    },
                },
            },
        },
        Template: api.ApplicationSetTemplate{
            ApplicationSetTemplateMeta: api.ApplicationSetTemplateMeta{
                Name: "app-{{cluster}}",
            },
            Spec: api.ApplicationSpec{
                Project: "default",
                Source: &api.ApplicationSource{
                    RepoURL:        "https://github.com/example/repo",
                    Path:           "helm/app",
                    TargetRevision: "main",
                },
                Destination: api.ApplicationDestination{
                    Server: "{{url}}",
                },
            },
        },
    },
})
```

### Error Handling

All API errors are returned as `*api.APIError`:

```go
import (
    "errors"

    goargocdclient "disaster37/goargocdclient"
    "disaster37/goargocdclient/api"
)

app, err := client.Application().Get("non-existent")
if err != nil {
    if api.IsNotFound(err) {
        fmt.Println("Application not found")
    } else if api.IsUnauthorized(err) {
        fmt.Println("Authentication required")
    } else if api.IsConflict(err) {
        fmt.Println("Resource conflict")
    } else {
        var apiErr *api.APIError
        if errors.As(err, &apiErr) {
            fmt.Printf("API error %d: %s\n", apiErr.Code, apiErr.Message)
        }
    }
}
```

### Get ArgoCD Server Info

```go
version, err := client.Version().Get()
fmt.Printf("ArgoCD %s (built %s, Go %s)\n", version.Version, version.BuildDate, version.GoVersion)

settings, err := client.Settings().Get()
fmt.Printf("Controller namespace: %s\n", settings.ControllerNamespace)

user, err := client.Session().GetUserInfo()
fmt.Printf("Logged in as: %s\n", user.Username)
```

## License

MIT - see [LICENSE](LICENSE)
