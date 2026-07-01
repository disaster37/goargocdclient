package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestApplicationList_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, ApplicationList{
			Items: []*ApplicationModel{
				{ObjectMeta: ObjectMeta{Name: "myapp"}},
			},
		})
	}))
	defer server.Close()

	a := NewApplication(client)
	list, err := a.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Items) != 1 || list.Items[0].Name != "myapp" {
		t.Errorf("unexpected list: %+v", list)
	}
}

func TestApplicationGet_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, ApplicationModel{ObjectMeta: ObjectMeta{Name: "myapp"}})
	}))
	defer server.Close()

	a := NewApplication(client)
	app, err := a.Get("myapp")
	if err != nil {
		t.Fatal(err)
	}
	if app.Name != "myapp" {
		t.Errorf("unexpected app: %+v", app)
	}
}

func TestApplicationCreate_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var app ApplicationModel
		json.NewDecoder(r.Body).Decode(&app)
		jsonResponse(w, 201, app)
	}))
	defer server.Close()

	a := NewApplication(client)
	app, err := a.Create(&ApplicationModel{ObjectMeta: ObjectMeta{Name: "newapp"}})
	if err != nil {
		t.Fatal(err)
	}
	if app.Name != "newapp" {
		t.Errorf("unexpected app: %+v", app)
	}
}

func TestApplicationUpdate_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, ApplicationModel{ObjectMeta: ObjectMeta{Name: "myapp"}})
	}))
	defer server.Close()

	a := NewApplication(client)
	app, err := a.Update(&ApplicationModel{ObjectMeta: ObjectMeta{Name: "myapp"}})
	if err != nil {
		t.Fatal(err)
	}
	if app.Name != "myapp" {
		t.Errorf("unexpected app: %+v", app)
	}
}

func TestApplicationDelete_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	a := NewApplication(client)
	if err := a.Delete("myapp", nil); err != nil {
		t.Fatal(err)
	}
}

func TestApplicationSync_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	a := NewApplication(client)
	if err := a.Sync("myapp", &SyncOptions{Prune: true}); err != nil {
		t.Fatal(err)
	}
}

func TestApplicationRollback_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	a := NewApplication(client)
	if err := a.Rollback("myapp", &RollbackOptions{ID: 1}); err != nil {
		t.Fatal(err)
	}
}

func TestApplicationTerminateOperation_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	a := NewApplication(client)
	if err := a.TerminateOperation("myapp"); err != nil {
		t.Fatal(err)
	}
}

func TestApplicationPatch_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, ApplicationModel{ObjectMeta: ObjectMeta{Name: "myapp"}})
	}))
	defer server.Close()

	a := NewApplication(client)
	app, err := a.Patch("myapp", map[string]string{"key": "value"}, "application/merge-patch+json")
	if err != nil {
		t.Fatal(err)
	}
	if app.Name != "myapp" {
		t.Errorf("unexpected app: %+v", app)
	}
}

func TestApplicationResourceTree_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, ApplicationTree{Nodes: []ResourceNode{{Name: "pod1"}}})
	}))
	defer server.Close()

	a := NewApplication(client)
	tree, err := a.ResourceTree("myapp")
	if err != nil {
		t.Fatal(err)
	}
	if len(tree.Nodes) != 1 || tree.Nodes[0].Name != "pod1" {
		t.Errorf("unexpected tree: %+v", tree)
	}
}

func TestApplicationManagedResources_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, ManagedResourcesResponse{Items: []*ResourceDiff{{Name: "svc"}}})
	}))
	defer server.Close()

	a := NewApplication(client)
	resp, err := a.ManagedResources("myapp")
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Items) != 1 || resp.Items[0].Name != "svc" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestApplicationGetManifests_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, ManifestResponse{Manifests: []string{"apiVersion: v1"}})
	}))
	defer server.Close()

	a := NewApplication(client)
	resp, err := a.GetManifests("myapp", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Manifests) != 1 {
		t.Errorf("unexpected manifests: %+v", resp)
	}
}

func TestApplicationRevisionMetadata_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, RevisionMetadata{Author: "dev", Message: "commit"})
	}))
	defer server.Close()

	a := NewApplication(client)
	meta, err := a.RevisionMetadata("myapp", "abc123")
	if err != nil {
		t.Fatal(err)
	}
	if meta.Author != "dev" {
		t.Errorf("unexpected metadata: %+v", meta)
	}
}

func TestApplicationGet_SyncWindows(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, SyncWindows{Windows: []*SyncWindow{{Kind: "allow"}}})
	}))
	defer server.Close()

	a := NewApplication(client)
	windows, err := a.GetSyncWindows("myapp")
	if err != nil {
		t.Fatal(err)
	}
	if len(windows.Windows) != 1 || windows.Windows[0].Kind != "allow" {
		t.Errorf("unexpected windows: %+v", windows)
	}
}

func TestApplicationListResourceEvents_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, ResourceEventList{Items: []ResourceEvent{{Reason: "Created"}}})
	}))
	defer server.Close()

	a := NewApplication(client)
	events, err := a.ListResourceEvents("myapp", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(events.Items) != 1 {
		t.Errorf("unexpected events: %+v", events)
	}
}

func TestApplicationListLinks_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, LinksResponse{Items: []LinkItem{{URL: "http://example.com"}}})
	}))
	defer server.Close()

	a := NewApplication(client)
	links, err := a.ListLinks("myapp")
	if err != nil {
		t.Fatal(err)
	}
	if len(links.Items) != 1 {
		t.Errorf("unexpected links: %+v", links)
	}
}

func TestApplicationGetResource_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, ApplicationResourceResponse{Manifest: "{}"})
	}))
	defer server.Close()

	a := NewApplication(client)
	resp, err := a.GetResource(&ApplicationResourceRequest{
		Name: "myapp", Namespace: "default", ResourceName: "pod1",
		Version: "v1", Kind: "Pod",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Manifest != "{}" {
		t.Errorf("unexpected manifest: %s", resp.Manifest)
	}
}

func TestApplicationListResourceActions_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, ResourceActionsList{Actions: []ResourceAction{{Name: "restart"}}})
	}))
	defer server.Close()

	a := NewApplication(client)
	actions, err := a.ListResourceActions(&ApplicationResourceRequest{
		Name: "myapp", ResourceName: "pod1", Version: "v1", Kind: "Pod",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(actions.Actions) != 1 || actions.Actions[0].Name != "restart" {
		t.Errorf("unexpected actions: %+v", actions)
	}
}

func TestApplicationGet_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 404, APIError{Code: 404, Message: "not found"})
	}))
	defer server.Close()

	a := NewApplication(client)
	_, err := a.Get("nonexistent")
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationList_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "internal error"})
	}))
	defer server.Close()

	a := NewApplication(client)
	_, err := a.List()
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationCreate_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "internal error"})
	}))
	defer server.Close()

	a := NewApplication(client)
	_, err := a.Create(&ApplicationModel{ObjectMeta: ObjectMeta{Name: "newapp"}})
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationUpdate_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "internal error"})
	}))
	defer server.Close()

	a := NewApplication(client)
	_, err := a.Update(&ApplicationModel{ObjectMeta: ObjectMeta{Name: "myapp"}})
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationDelete_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "internal error"})
	}))
	defer server.Close()

	a := NewApplication(client)
	err := a.Delete("myapp", nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationSync_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "internal error"})
	}))
	defer server.Close()

	a := NewApplication(client)
	err := a.Sync("myapp", &SyncOptions{Prune: true})
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationRollback_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "internal error"})
	}))
	defer server.Close()

	a := NewApplication(client)
	err := a.Rollback("myapp", &RollbackOptions{ID: 1})
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationTerminateOperation_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "internal error"})
	}))
	defer server.Close()

	a := NewApplication(client)
	err := a.TerminateOperation("myapp")
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationPatch_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "internal error"})
	}))
	defer server.Close()

	a := NewApplication(client)
	_, err := a.Patch("myapp", map[string]string{"key": "value"}, "application/json")
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationResourceTree_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "internal error"})
	}))
	defer server.Close()

	a := NewApplication(client)
	_, err := a.ResourceTree("myapp")
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationManagedResources_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "internal error"})
	}))
	defer server.Close()

	a := NewApplication(client)
	_, err := a.ManagedResources("myapp")
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationGetManifests_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "internal error"})
	}))
	defer server.Close()

	a := NewApplication(client)
	_, err := a.GetManifests("myapp", nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationRevisionMetadata_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "internal error"})
	}))
	defer server.Close()

	a := NewApplication(client)
	_, err := a.RevisionMetadata("myapp", "abc123")
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationGetSyncWindows_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "internal error"})
	}))
	defer server.Close()

	a := NewApplication(client)
	_, err := a.GetSyncWindows("myapp")
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationListResourceEvents_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "internal error"})
	}))
	defer server.Close()

	a := NewApplication(client)
	_, err := a.ListResourceEvents("myapp", nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationListLinks_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "internal error"})
	}))
	defer server.Close()

	a := NewApplication(client)
	_, err := a.ListLinks("myapp")
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationGetResource_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "internal error"})
	}))
	defer server.Close()

	a := NewApplication(client)
	_, err := a.GetResource(&ApplicationResourceRequest{
		Name: "myapp", Namespace: "default", ResourceName: "pod1",
		Version: "v1", Kind: "Pod",
	})
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationListResourceActions_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "internal error"})
	}))
	defer server.Close()

	a := NewApplication(client)
	_, err := a.ListResourceActions(&ApplicationResourceRequest{
		Name: "myapp", ResourceName: "pod1", Version: "v1", Kind: "Pod",
	})
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplication_NetworkError(t *testing.T) {
	client := newFailingClient()
	a := NewApplication(client)
	_, err := a.List()
	if err == nil {
		t.Error("expected error")
	}
	_, err = a.Get("myapp")
	if err == nil {
		t.Error("expected error")
	}
	_, err = a.Create(&ApplicationModel{ObjectMeta: ObjectMeta{Name: "newapp"}})
	if err == nil {
		t.Error("expected error")
	}
	_, err = a.Update(&ApplicationModel{ObjectMeta: ObjectMeta{Name: "myapp"}})
	if err == nil {
		t.Error("expected error")
	}
	e := a.Delete("myapp", nil)
	if e == nil {
		t.Error("expected error")
	}
	e = a.Sync("myapp", &SyncOptions{Prune: true})
	if e == nil {
		t.Error("expected error")
	}
	e = a.Rollback("myapp", &RollbackOptions{ID: 1})
	if e == nil {
		t.Error("expected error")
	}
	e = a.TerminateOperation("myapp")
	if e == nil {
		t.Error("expected error")
	}
	_, err = a.Patch("myapp", map[string]string{"key": "value"}, "application/json")
	if err == nil {
		t.Error("expected error")
	}
	_, err = a.GetResource(&ApplicationResourceRequest{Name: "myapp", Namespace: "default", ResourceName: "pod1", Version: "v1", Kind: "Pod"})
	if err == nil {
		t.Error("expected error")
	}
	_, err = a.PatchResource(&ApplicationResourcePatchRequest{
		ApplicationResourceRequest: ApplicationResourceRequest{Name: "myapp", Namespace: "default", ResourceName: "pod1", Version: "v1", Kind: "Pod"},
		Patch: "{}", PatchType: "application/json",
	})
	if err == nil {
		t.Error("expected error")
	}
	e = a.DeleteResource(&ApplicationResourceDeleteRequest{
		ApplicationResourceRequest: ApplicationResourceRequest{Name: "myapp", Namespace: "default", ResourceName: "pod1", Version: "v1", Kind: "Pod"},
	})
	if e == nil {
		t.Error("expected error")
	}
	_, err = a.ListResourceActions(&ApplicationResourceRequest{Name: "myapp", ResourceName: "pod1", Version: "v1", Kind: "Pod"})
	if err == nil {
		t.Error("expected error")
	}
	e = a.RunResourceAction(&ApplicationResourceActionRequest{
		ApplicationResourceRequest: ApplicationResourceRequest{Name: "myapp", Namespace: "default", ResourceName: "pod1", Version: "v1", Kind: "Pod"},
		Action: "restart",
	})
	if e == nil {
		t.Error("expected error")
	}
	_, err = a.GetManifests("myapp", nil)
	if err == nil {
		t.Error("expected error")
	}
	_, err = a.ResourceTree("myapp")
	if err == nil {
		t.Error("expected error")
	}
	_, err = a.ManagedResources("myapp")
	if err == nil {
		t.Error("expected error")
	}
	_, err = a.RevisionMetadata("myapp", "abc123")
	if err == nil {
		t.Error("expected error")
	}
	_, err = a.GetSyncWindows("myapp")
	if err == nil {
		t.Error("expected error")
	}
	_, err = a.ListResourceEvents("myapp", nil)
	if err == nil {
		t.Error("expected error")
	}
	_, err = a.ListLinks("myapp")
	if err == nil {
		t.Error("expected error")
	}
	_, err = a.ListResourceLinks(&ApplicationResourceRequest{Name: "myapp", Namespace: "default", ResourceName: "pod1", Version: "v1", Kind: "Pod"})
	if err == nil {
		t.Error("expected error")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = a.Watch(ctx, &WatchOptions{})
	if err == nil {
		t.Error("expected error")
	}
	_, err = a.WatchResourceTree(ctx, "myapp")
	if err == nil {
		t.Error("expected error")
	}
	_, err = a.PodLogs(ctx, &PodLogsOptions{Name: "myapp", Namespace: "default", ResourceName: "pod1", Kind: "Pod"})
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationPatchResource_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, ApplicationResourceResponse{Manifest: "{}"})
	}))
	defer server.Close()

	a := NewApplication(client)
	resp, err := a.PatchResource(&ApplicationResourcePatchRequest{
		ApplicationResourceRequest: ApplicationResourceRequest{
			Name: "myapp", Namespace: "default", ResourceName: "pod1",
			Version: "v1", Kind: "Pod",
		},
		Patch: "{}", PatchType: "application/json",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Manifest != "{}" {
		t.Errorf("unexpected manifest: %s", resp.Manifest)
	}
}

func TestApplicationPatchResource_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "internal error"})
	}))
	defer server.Close()

	a := NewApplication(client)
	_, err := a.PatchResource(&ApplicationResourcePatchRequest{
		ApplicationResourceRequest: ApplicationResourceRequest{
			Name: "myapp", Namespace: "default", ResourceName: "pod1",
			Version: "v1", Kind: "Pod",
		},
		Patch: "{}", PatchType: "application/json",
	})
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationDeleteResource_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	a := NewApplication(client)
	err := a.DeleteResource(&ApplicationResourceDeleteRequest{
		ApplicationResourceRequest: ApplicationResourceRequest{
			Name: "myapp", Namespace: "default", ResourceName: "pod1",
			Version: "v1", Kind: "Pod",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestApplicationDeleteResource_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "internal error"})
	}))
	defer server.Close()

	a := NewApplication(client)
	err := a.DeleteResource(&ApplicationResourceDeleteRequest{
		ApplicationResourceRequest: ApplicationResourceRequest{
			Name: "myapp", Namespace: "default", ResourceName: "pod1",
			Version: "v1", Kind: "Pod",
		},
	})
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationDeleteResource_WithForce(t *testing.T) {
	force := true
	var queryChecked bool
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("force") == "true" {
			queryChecked = true
		}
		w.WriteHeader(200)
	}))
	defer server.Close()

	a := NewApplication(client)
	err := a.DeleteResource(&ApplicationResourceDeleteRequest{
		ApplicationResourceRequest: ApplicationResourceRequest{
			Name: "myapp", Namespace: "default", ResourceName: "pod1",
			Version: "v1", Kind: "Pod",
		},
		Force: &force,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !queryChecked {
		t.Error("force query param was not set")
	}
}

func TestApplicationRunResourceAction_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	a := NewApplication(client)
	err := a.RunResourceAction(&ApplicationResourceActionRequest{
		ApplicationResourceRequest: ApplicationResourceRequest{
			Name: "myapp", Namespace: "default", ResourceName: "pod1",
			Version: "v1", Kind: "Pod",
		},
		Action: "restart",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestApplicationRunResourceAction_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "internal error"})
	}))
	defer server.Close()

	a := NewApplication(client)
	err := a.RunResourceAction(&ApplicationResourceActionRequest{
		ApplicationResourceRequest: ApplicationResourceRequest{
			Name: "myapp", Namespace: "default", ResourceName: "pod1",
			Version: "v1", Kind: "Pod",
		},
		Action: "restart",
	})
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationListResourceLinks_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, LinksResponse{Items: []LinkItem{{URL: "http://example.com"}}})
	}))
	defer server.Close()

	a := NewApplication(client)
	links, err := a.ListResourceLinks(&ApplicationResourceRequest{
		Name: "myapp", Namespace: "default", ResourceName: "pod1",
		Version: "v1", Kind: "Pod",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(links.Items) != 1 {
		t.Errorf("unexpected links: %+v", links)
	}
}

func TestApplicationListResourceLinks_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "internal error"})
	}))
	defer server.Close()

	a := NewApplication(client)
	_, err := a.ListResourceLinks(&ApplicationResourceRequest{
		Name: "myapp", Namespace: "default", ResourceName: "pod1",
		Version: "v1", Kind: "Pod",
	})
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationDelete_WithCascade(t *testing.T) {
	cascade := true
	var queryChecked bool
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("cascade") == "true" {
			queryChecked = true
		}
		w.WriteHeader(200)
	}))
	defer server.Close()

	a := NewApplication(client)
	err := a.Delete("myapp", &ApplicationDeleteOptions{Cascade: &cascade})
	if err != nil {
		t.Fatal(err)
	}
	if !queryChecked {
		t.Error("cascade query param was not set")
	}
}

func TestApplicationDelete_WithPropagationPolicy(t *testing.T) {
	var queryChecked bool
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("propagationPolicy") == "background" {
			queryChecked = true
		}
		w.WriteHeader(200)
	}))
	defer server.Close()

	a := NewApplication(client)
	err := a.Delete("myapp", &ApplicationDeleteOptions{PropagationPolicy: "background"})
	if err != nil {
		t.Fatal(err)
	}
	if !queryChecked {
		t.Error("propagationPolicy query param was not set")
	}
}

func TestApplicationPatch_EmptyPatchType(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, ApplicationModel{ObjectMeta: ObjectMeta{Name: "myapp"}})
	}))
	defer server.Close()

	a := NewApplication(client)
	app, err := a.Patch("myapp", map[string]string{"key": "value"}, "")
	if err != nil {
		t.Fatal(err)
	}
	if app.Name != "myapp" {
		t.Errorf("unexpected app: %+v", app)
	}
}

func TestApplicationGetManifests_WithOpts(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, ManifestResponse{Manifests: []string{"apiVersion: v1"}})
	}))
	defer server.Close()

	a := NewApplication(client)
	resp, err := a.GetManifests("myapp", &ApplicationManifestQuery{Revision: "abc123"})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Manifests) != 1 {
		t.Errorf("unexpected manifests: %+v", resp)
	}
}

func TestApplicationListResourceEvents_WithOpts(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, ResourceEventList{Items: []ResourceEvent{{Reason: "Created"}}})
	}))
	defer server.Close()

	a := NewApplication(client)
	events, err := a.ListResourceEvents("myapp", &ApplicationResourceEventsQuery{
		ResourceName: "pod1", ResourceUID: "uid123",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(events.Items) != 1 {
		t.Errorf("unexpected events: %+v", events)
	}
}

func TestApplicationWatch_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(200)
		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Error("expected http.Flusher")
			return
		}
		// Send non-matching line first
		fmt.Fprintf(w, "event: ping\n\n")
		flusher.Flush()
		// Then bad JSON line
		fmt.Fprintf(w, "data: not-json\n\n")
		flusher.Flush()
		// Then valid data
		data, _ := json.Marshal(struct {
			Result ApplicationWatchEvent `json:"result"`
		}{
			Result: ApplicationWatchEvent{
				Type:        SyncStatusCodeSynced,
				Application: &ApplicationModel{ObjectMeta: ObjectMeta{Name: "myapp"}},
			},
		})
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	a := NewApplication(client)
	ch, err := a.Watch(ctx, &WatchOptions{Revision: "abc123"})
	if err != nil {
		t.Fatal(err)
	}
	select {
	case event, ok := <-ch:
		if !ok {
			t.Error("channel closed unexpectedly")
			return
		}
		if event.Application.Name != "myapp" {
			t.Errorf("unexpected event: %+v", event)
		}
		cancel()
	case <-time.After(2 * time.Second):
		t.Error("timed out waiting for watch event")
	}
}

func TestApplicationWatch_SSEWatchErr(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(APIError{Code: 400, Message: "bad request"})
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	a := NewApplication(client)
	_, err := a.Watch(ctx, &WatchOptions{})
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationWatch_SSEError(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("internal error"))
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	a := NewApplication(client)
	_, err := a.Watch(ctx, &WatchOptions{})
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationWatchResourceTree_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(200)
		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Error("expected http.Flusher")
			return
		}
		fmt.Fprintf(w, "event: ping\n\n")
		flusher.Flush()
		data, _ := json.Marshal(struct {
			Result ApplicationTree `json:"result"`
		}{
			Result: ApplicationTree{Nodes: []ResourceNode{{Name: "pod1"}}},
		})
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	a := NewApplication(client)
	ch, err := a.WatchResourceTree(ctx, "myapp")
	if err != nil {
		t.Fatal(err)
	}
	select {
	case tree, ok := <-ch:
		if !ok {
			t.Error("channel closed unexpectedly")
			return
		}
		if len(tree.Nodes) != 1 || tree.Nodes[0].Name != "pod1" {
			t.Errorf("unexpected tree: %+v", tree)
		}
		cancel()
	case <-time.After(2 * time.Second):
		t.Error("timed out waiting for resource tree event")
	}
}

func TestApplicationWatchResourceTree_SSEError(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("internal error"))
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	a := NewApplication(client)
	_, err := a.WatchResourceTree(ctx, "myapp")
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationPodLogs_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(200)
		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Error("expected http.Flusher")
			return
		}
		fmt.Fprintf(w, "event: ping\n\n")
		flusher.Flush()
		data, _ := json.Marshal(struct {
			Result LogEntry `json:"result"`
		}{
			Result: LogEntry{Content: "log line 1"},
		})
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	a := NewApplication(client)
	ch, err := a.PodLogs(ctx, &PodLogsOptions{
		Name: "myapp", Namespace: "default", ResourceName: "pod1",
		Kind: "Pod",
	})
	if err != nil {
		t.Fatal(err)
	}
	select {
	case logEntry, ok := <-ch:
		if !ok {
			t.Error("channel closed unexpectedly")
			return
		}
		if logEntry.Content != "log line 1" {
			t.Errorf("unexpected log entry: %+v", logEntry)
		}
		cancel()
	case <-time.After(2 * time.Second):
		t.Error("timed out waiting for log entry")
	}
}

func TestApplicationPodLogs_SSEError(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(APIError{Code: 400, Message: "bad request"})
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	a := NewApplication(client)
	_, err := a.PodLogs(ctx, &PodLogsOptions{
		Name: "myapp", Namespace: "default", ResourceName: "pod1",
		Kind: "Pod",
	})
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationPodLogs_NilOpts(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()

	a := NewApplication(client)
	_, err := a.PodLogs(context.Background(), nil)
	if err == nil {
		t.Error("expected error for nil opts")
	}
}

func TestApplicationPodLogs_WithOptionalParams(t *testing.T) {
	paramsChecked := false
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("container") == "mycontainer" &&
			q.Get("tailLines") == "100" &&
			q.Get("sinceSeconds") == "3600" &&
			q.Get("sinceTime") == "2024-01-01T00:00:00Z" &&
			q.Get("previous") == "true" {
			paramsChecked = true
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(200)
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	a := NewApplication(client)
	_, err := a.PodLogs(ctx, &PodLogsOptions{
		Name: "myapp", Namespace: "default", ResourceName: "pod1",
		Kind: "Pod", Container: "mycontainer", TailLines: 100,
		SinceSeconds: 3600, SinceTime: "2024-01-01T00:00:00Z", Previous: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !paramsChecked {
		t.Error("optional params were not set")
	}
}

func TestApplicationWatch_CtxDone(t *testing.T) {
	done := make(chan struct{})
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(200)
		flusher, ok := w.(http.Flusher)
		if !ok {
			return
		}
		data, _ := json.Marshal(struct {
			Result ApplicationWatchEvent `json:"result"`
		}{
			Result: ApplicationWatchEvent{
				Type:        SyncStatusCodeSynced,
				Application: &ApplicationModel{ObjectMeta: ObjectMeta{Name: "myapp"}},
			},
		})
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
		<-done
	}))
	defer func() {
		close(done)
		server.Close()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	a := NewApplication(client)
	ch, err := a.Watch(ctx, &WatchOptions{})
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)
	cancel()
	_ = <-ch
}

func TestApplicationWatch_NilOpts(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(200)
	}))
	defer server.Close()

	a := NewApplication(client)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch, err := a.Watch(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if ch == nil {
		t.Error("expected non-nil channel")
	}
}
