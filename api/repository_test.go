package api

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestRepositoryList_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, RepositoryList{Items: []*RepositoryModel{
			{Repo: "https://github.com/org/repo"},
		}})
	}))
	defer server.Close()

	r := NewRepository(client)
	list, err := r.List(nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Items) != 1 || list.Items[0].Repo != "https://github.com/org/repo" {
		t.Errorf("unexpected list: %+v", list)
	}
}

func TestRepositoryGet_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, RepositoryModel{Repo: "https://github.com/org/repo", Name: "my-repo"})
	}))
	defer server.Close()

	r := NewRepository(client)
	repo, err := r.Get("https://github.com/org/repo", nil)
	if err != nil {
		t.Fatal(err)
	}
	if repo.Name != "my-repo" {
		t.Errorf("unexpected repo: %+v", repo)
	}
}

func TestRepositoryCreate_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var repo RepositoryModel
		json.NewDecoder(r.Body).Decode(&repo)
		jsonResponse(w, 201, repo)
	}))
	defer server.Close()

	r := NewRepository(client)
	repo, err := r.Create(&RepositoryModel{Repo: "https://github.com/org/new", Name: "new-repo"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if repo.Name != "new-repo" {
		t.Errorf("unexpected repo: %+v", repo)
	}
}

func TestRepositoryUpdate_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, RepositoryModel{Repo: "https://github.com/org/repo", Name: "updated"})
	}))
	defer server.Close()

	r := NewRepository(client)
	repo, err := r.Update(&RepositoryModel{Repo: "https://github.com/org/repo", Name: "updated"})
	if err != nil {
		t.Fatal(err)
	}
	if repo.Name != "updated" {
		t.Errorf("unexpected repo: %+v", repo)
	}
}

func TestRepositoryDelete_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	r := NewRepository(client)
	if err := r.Delete("https://github.com/org/repo", nil); err != nil {
		t.Fatal(err)
	}
}

func TestRepositoryListApps_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, RepositoryAppsList{Items: []*RepoApp{
			{RepoURL: "https://github.com/org/repo", Path: "."},
		}})
	}))
	defer server.Close()

	r := NewRepository(client)
	apps, err := r.ListApps("https://github.com/org/repo", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(apps.Items) != 1 || apps.Items[0].Path != "." {
		t.Errorf("unexpected apps: %+v", apps)
	}
}

func TestRepositoryGetHelmCharts_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, HelmChartsResponse{Items: []*HelmChart{
			{Name: "nginx", Versions: []string{"1.0.0"}},
		}})
	}))
	defer server.Close()

	r := NewRepository(client)
	charts, err := r.GetHelmCharts("https://charts.helm.sh/stable", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(charts.Items) != 1 || charts.Items[0].Name != "nginx" {
		t.Errorf("unexpected charts: %+v", charts)
	}
}

func TestRepositoryListRefs_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, RefsResponse{Branches: []string{"main"}, Tags: []string{"v1.0"}})
	}))
	defer server.Close()

	r := NewRepository(client)
	refs, err := r.ListRefs("https://github.com/org/repo", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(refs.Branches) != 1 || refs.Branches[0] != "main" {
		t.Errorf("unexpected refs: %+v", refs)
	}
}

func TestRepositoryValidateAccess_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	r := NewRepository(client)
	if err := r.ValidateAccess(&RepoAccessQuery{Repo: "https://github.com/org/repo"}); err != nil {
		t.Fatal(err)
	}
}

func TestRepositoryListOCITags_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, OCITagsResponse{Tags: []string{"latest", "v1.0"}})
	}))
	defer server.Close()

	r := NewRepository(client)
	tags, err := r.ListOCITags("oci://registry.example.com/repo", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(tags.Tags) != 2 {
		t.Errorf("unexpected tags: %+v", tags)
	}
}

func TestRepositoryGetWrite_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, RepositoryModel{Repo: "https://github.com/org/repo", Name: "write-repo"})
	}))
	defer server.Close()

	r := NewRepository(client)
	repo, err := r.GetWrite("https://github.com/org/repo", nil)
	if err != nil {
		t.Fatal(err)
	}
	if repo.Name != "write-repo" {
		t.Errorf("unexpected repo: %+v", repo)
	}
}

func TestRepositoryCreateWriteRepo_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var repo RepositoryModel
		json.NewDecoder(r.Body).Decode(&repo)
		jsonResponse(w, 201, repo)
	}))
	defer server.Close()

	r := NewRepository(client)
	repo, err := r.CreateWriteRepository(&RepositoryModel{Repo: "https://github.com/org/write", Name: "write"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if repo.Name != "write" {
		t.Errorf("unexpected repo: %+v", repo)
	}
}

func TestRepositoryUpdateWriteRepo_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, RepositoryModel{Repo: "https://github.com/org/write", Name: "updated-write"})
	}))
	defer server.Close()

	r := NewRepository(client)
	repo, err := r.UpdateWriteRepository(&RepositoryModel{Repo: "https://github.com/org/write", Name: "updated-write"})
	if err != nil {
		t.Fatal(err)
	}
	if repo.Name != "updated-write" {
		t.Errorf("unexpected repo: %+v", repo)
	}
}

func TestRepositoryDeleteWriteRepo_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	r := NewRepository(client)
	if err := r.DeleteWriteRepository("https://github.com/org/write", nil); err != nil {
		t.Fatal(err)
	}
}

func TestRepositoryValidateWriteAccess_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	r := NewRepository(client)
	if err := r.ValidateWriteAccess(&RepoAccessQuery{Repo: "https://github.com/org/repo"}); err != nil {
		t.Fatal(err)
	}
}

func TestRepositoryGet_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 404, APIError{Code: 404, Message: "not found"})
	}))
	defer server.Close()

	r := NewRepository(client)
	_, err := r.Get("https://unknown", nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestRepositoryGetAppDetails_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, RepoAppDetails{Type: "Helm"})
	}))
	defer server.Close()

	r := NewRepository(client)
	details, err := r.GetAppDetails(&RepoAppDetailsQuery{
		Source: ApplicationSource{RepoURL: "https://github.com/org/repo"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if details.Type != "Helm" {
		t.Errorf("unexpected details: %+v", details)
	}
}

func TestRepositoryListWriteRepositories_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, RepositoryList{Items: []*RepositoryModel{
			{Repo: "https://github.com/org/write-repo"},
		}})
	}))
	defer server.Close()

	r := NewRepository(client)
	list, err := r.ListWriteRepositories(nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Items) != 1 || list.Items[0].Repo != "https://github.com/org/write-repo" {
		t.Errorf("unexpected list: %+v", list)
	}
}

func TestRepositoryList_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	r := NewRepository(client)
	_, err := r.List(nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestRepositoryCreate_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	r := NewRepository(client)
	_, err := r.Create(&RepositoryModel{Repo: "https://github.com/org/new", Name: "new-repo"}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestRepositoryUpdate_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	r := NewRepository(client)
	_, err := r.Update(&RepositoryModel{Repo: "https://github.com/org/repo", Name: "updated"})
	if err == nil {
		t.Error("expected error")
	}
}

func TestRepositoryDelete_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	r := NewRepository(client)
	err := r.Delete("https://github.com/org/repo", nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestRepositoryListApps_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	r := NewRepository(client)
	_, err := r.ListApps("https://github.com/org/repo", nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestRepositoryGetAppDetails_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	r := NewRepository(client)
	_, err := r.GetAppDetails(&RepoAppDetailsQuery{
		Source: ApplicationSource{RepoURL: "https://github.com/org/repo"},
	})
	if err == nil {
		t.Error("expected error")
	}
}

func TestRepositoryGetHelmCharts_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	r := NewRepository(client)
	_, err := r.GetHelmCharts("https://charts.helm.sh/stable", nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestRepositoryListRefs_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	r := NewRepository(client)
	_, err := r.ListRefs("https://github.com/org/repo", nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestRepositoryListOCITags_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	r := NewRepository(client)
	_, err := r.ListOCITags("oci://registry.example.com/repo", nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestRepositoryValidateAccess_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	r := NewRepository(client)
	err := r.ValidateAccess(&RepoAccessQuery{Repo: "https://github.com/org/repo"})
	if err == nil {
		t.Error("expected error")
	}
}

func TestRepositoryListWriteRepositories_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	r := NewRepository(client)
	_, err := r.ListWriteRepositories(nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestRepositoryGetWrite_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	r := NewRepository(client)
	_, err := r.GetWrite("https://github.com/org/repo", nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestRepositoryCreateWriteRepository_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	r := NewRepository(client)
	_, err := r.CreateWriteRepository(&RepositoryModel{Repo: "https://github.com/org/write", Name: "write"}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestRepositoryUpdateWriteRepository_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	r := NewRepository(client)
	_, err := r.UpdateWriteRepository(&RepositoryModel{Repo: "https://github.com/org/write", Name: "updated-write"})
	if err == nil {
		t.Error("expected error")
	}
}

func TestRepositoryDeleteWriteRepository_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	r := NewRepository(client)
	err := r.DeleteWriteRepository("https://github.com/org/write", nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestRepositoryValidateWriteAccess_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	r := NewRepository(client)
	err := r.ValidateWriteAccess(&RepoAccessQuery{Repo: "https://github.com/org/repo"})
	if err == nil {
		t.Error("expected error")
	}
}

func TestRepository_NetworkError(t *testing.T) {
	client := newFailingClient()
	r := NewRepository(client)
	_, err := r.List(nil)
	if err == nil {
		t.Error("expected error")
	}
	_, err = r.Get("https://github.com/org/repo", nil)
	if err == nil {
		t.Error("expected error")
	}
	_, err = r.Create(&RepositoryModel{Repo: "https://github.com/org/new", Name: "new-repo"}, nil)
	if err == nil {
		t.Error("expected error")
	}
	_, err = r.Update(&RepositoryModel{Repo: "https://github.com/org/repo", Name: "updated"})
	if err == nil {
		t.Error("expected error")
	}
	if err = r.Delete("https://github.com/org/repo", nil); err == nil {
		t.Error("expected error")
	}
	_, err = r.ListApps("https://github.com/org/repo", nil)
	if err == nil {
		t.Error("expected error")
	}
	_, err = r.GetAppDetails(&RepoAppDetailsQuery{Source: ApplicationSource{RepoURL: "https://github.com/org/repo"}})
	if err == nil {
		t.Error("expected error")
	}
	_, err = r.GetHelmCharts("https://charts.helm.sh/stable", nil)
	if err == nil {
		t.Error("expected error")
	}
	_, err = r.ListRefs("https://github.com/org/repo", nil)
	if err == nil {
		t.Error("expected error")
	}
	_, err = r.ListOCITags("oci://registry.example.com/repo", nil)
	if err == nil {
		t.Error("expected error")
	}
	if err = r.ValidateAccess(&RepoAccessQuery{Repo: "https://github.com/org/repo"}); err == nil {
		t.Error("expected error")
	}
	_, err = r.ListWriteRepositories(nil)
	if err == nil {
		t.Error("expected error")
	}
	_, err = r.GetWrite("https://github.com/org/repo", nil)
	if err == nil {
		t.Error("expected error")
	}
	_, err = r.CreateWriteRepository(&RepositoryModel{Repo: "https://github.com/org/write", Name: "write"}, nil)
	if err == nil {
		t.Error("expected error")
	}
	_, err = r.UpdateWriteRepository(&RepositoryModel{Repo: "https://github.com/org/write", Name: "updated-write"})
	if err == nil {
		t.Error("expected error")
	}
	if err = r.DeleteWriteRepository("https://github.com/org/write", nil); err == nil {
		t.Error("expected error")
	}
	if err = r.ValidateWriteAccess(&RepoAccessQuery{Repo: "https://github.com/org/repo"}); err == nil {
		t.Error("expected error")
	}
}
