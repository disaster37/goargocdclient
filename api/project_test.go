package api

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestProjectList_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, ProjectList{Items: []*ProjectModel{
			{ObjectMeta: ObjectMeta{Name: "default"}},
		}})
	}))
	defer server.Close()

	p := NewProject(client)
	list, err := p.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Items) != 1 || list.Items[0].Name != "default" {
		t.Errorf("unexpected list: %+v", list)
	}
}

func TestProjectGet_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, ProjectModel{ObjectMeta: ObjectMeta{Name: "default"}})
	}))
	defer server.Close()

	p := NewProject(client)
	proj, err := p.Get("default")
	if err != nil {
		t.Fatal(err)
	}
	if proj.Name != "default" {
		t.Errorf("unexpected project: %+v", proj)
	}
}

func TestProjectCreate_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var proj ProjectModel
		json.NewDecoder(r.Body).Decode(&proj)
		jsonResponse(w, 201, proj)
	}))
	defer server.Close()

	p := NewProject(client)
	proj, err := p.Create(&ProjectModel{ObjectMeta: ObjectMeta{Name: "newproj"}})
	if err != nil {
		t.Fatal(err)
	}
	if proj.Name != "newproj" {
		t.Errorf("unexpected project: %+v", proj)
	}
}

func TestProjectUpdate_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, ProjectModel{ObjectMeta: ObjectMeta{Name: "updated"}})
	}))
	defer server.Close()

	p := NewProject(client)
	proj, err := p.Update(&ProjectModel{ObjectMeta: ObjectMeta{Name: "updated"}})
	if err != nil {
		t.Fatal(err)
	}
	if proj.Name != "updated" {
		t.Errorf("unexpected project: %+v", proj)
	}
}

func TestProjectDelete_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	p := NewProject(client)
	if err := p.Delete("test"); err != nil {
		t.Fatal(err)
	}
}

func TestProjectGetDetailed_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, ProjectDetailed{
			ProjectModel: ProjectModel{ObjectMeta: ObjectMeta{Name: "default"}},
		})
	}))
	defer server.Close()

	p := NewProject(client)
	detail, err := p.GetDetailed("default")
	if err != nil {
		t.Fatal(err)
	}
	if detail.Name != "default" {
		t.Errorf("unexpected detailed: %+v", detail)
	}
}

func TestProjectCreateToken_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, TokenResponse{Token: "proj-token"})
	}))
	defer server.Close()

	p := NewProject(client)
	resp, err := p.CreateToken("default", "read-only", &ProjectTokenCreateOptions{ID: "id1", ExpiresIn: 3600})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Token != "proj-token" {
		t.Errorf("unexpected token: %s", resp.Token)
	}
}

func TestProjectDeleteToken_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	p := NewProject(client)
	if err := p.DeleteToken("default", "read-only", 12345); err != nil {
		t.Fatal(err)
	}
}

func TestProjectListEvents_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, ResourceEventList{Items: []ResourceEvent{{Reason: "Updated"}}})
	}))
	defer server.Close()

	p := NewProject(client)
	events, err := p.ListEvents("default")
	if err != nil {
		t.Fatal(err)
	}
	if len(events.Items) != 1 {
		t.Errorf("unexpected events: %+v", events)
	}
}

func TestProjectGetSyncWindowsState_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, SyncWindows{Windows: []*SyncWindow{{Kind: "deny"}}})
	}))
	defer server.Close()

	p := NewProject(client)
	windows, err := p.GetSyncWindowsState("default")
	if err != nil {
		t.Fatal(err)
	}
	if len(windows.Windows) != 1 || windows.Windows[0].Kind != "deny" {
		t.Errorf("unexpected windows: %+v", windows)
	}
}

func TestProjectListLinks_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, LinksResponse{Items: []LinkItem{{URL: "http://example.com"}}})
	}))
	defer server.Close()

	p := NewProject(client)
	links, err := p.ListLinks("default")
	if err != nil {
		t.Fatal(err)
	}
	if len(links.Items) != 1 {
		t.Errorf("unexpected links: %+v", links)
	}
}

func TestProjectGet_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 404, APIError{Code: 404, Message: "not found"})
	}))
	defer server.Close()

	p := NewProject(client)
	_, err := p.Get("nonexistent")
	if err == nil {
		t.Error("expected error")
	}
}

func TestProjectGetGlobalProjects_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, ProjectGlobalResponse{
			ProjectModel: ProjectModel{ObjectMeta: ObjectMeta{Name: "default"}},
		})
	}))
	defer server.Close()

	p := NewProject(client)
	resp, err := p.GetGlobalProjects("default")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Name != "default" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestProjectList_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	p := NewProject(client)
	_, err := p.List()
	if err == nil {
		t.Error("expected error")
	}
}

func TestProjectCreate_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	p := NewProject(client)
	_, err := p.Create(&ProjectModel{ObjectMeta: ObjectMeta{Name: "newproj"}})
	if err == nil {
		t.Error("expected error")
	}
}

func TestProjectUpdate_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	p := NewProject(client)
	_, err := p.Update(&ProjectModel{ObjectMeta: ObjectMeta{Name: "updated"}})
	if err == nil {
		t.Error("expected error")
	}
}

func TestProjectDelete_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	p := NewProject(client)
	err := p.Delete("test")
	if err == nil {
		t.Error("expected error")
	}
}

func TestProjectGetDetailed_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	p := NewProject(client)
	_, err := p.GetDetailed("default")
	if err == nil {
		t.Error("expected error")
	}
}

func TestProjectGetGlobalProjects_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	p := NewProject(client)
	_, err := p.GetGlobalProjects("default")
	if err == nil {
		t.Error("expected error")
	}
}

func TestProjectCreateToken_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	p := NewProject(client)
	_, err := p.CreateToken("default", "read-only", &ProjectTokenCreateOptions{ID: "id1", ExpiresIn: 3600})
	if err == nil {
		t.Error("expected error")
	}
}

func TestProjectDeleteToken_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	p := NewProject(client)
	err := p.DeleteToken("default", "read-only", 12345)
	if err == nil {
		t.Error("expected error")
	}
}

func TestProjectListEvents_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	p := NewProject(client)
	_, err := p.ListEvents("default")
	if err == nil {
		t.Error("expected error")
	}
}

func TestProjectGetSyncWindowsState_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	p := NewProject(client)
	_, err := p.GetSyncWindowsState("default")
	if err == nil {
		t.Error("expected error")
	}
}

func TestProjectListLinks_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	p := NewProject(client)
	_, err := p.ListLinks("default")
	if err == nil {
		t.Error("expected error")
	}
}

func TestProject_NetworkError(t *testing.T) {
	client := newFailingClient()
	p := NewProject(client)
	_, err := p.List()
	if err == nil {
		t.Error("expected error")
	}
	_, err = p.Get("default")
	if err == nil {
		t.Error("expected error")
	}
	_, err = p.Create(&ProjectModel{ObjectMeta: ObjectMeta{Name: "newproj"}})
	if err == nil {
		t.Error("expected error")
	}
	_, err = p.Update(&ProjectModel{ObjectMeta: ObjectMeta{Name: "updated"}})
	if err == nil {
		t.Error("expected error")
	}
	if err = p.Delete("test"); err == nil {
		t.Error("expected error")
	}
	_, err = p.GetDetailed("default")
	if err == nil {
		t.Error("expected error")
	}
	_, err = p.GetGlobalProjects("default")
	if err == nil {
		t.Error("expected error")
	}
	_, err = p.CreateToken("default", "read-only", &ProjectTokenCreateOptions{ID: "id1", ExpiresIn: 3600})
	if err == nil {
		t.Error("expected error")
	}
	if err = p.DeleteToken("default", "read-only", 12345); err == nil {
		t.Error("expected error")
	}
	_, err = p.ListEvents("default")
	if err == nil {
		t.Error("expected error")
	}
	_, err = p.GetSyncWindowsState("default")
	if err == nil {
		t.Error("expected error")
	}
	_, err = p.ListLinks("default")
	if err == nil {
		t.Error("expected error")
	}
}
