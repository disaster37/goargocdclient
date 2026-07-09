package api

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestClusterList_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, ClusterList{Items: []*ClusterModel{
			{Name: "in-cluster", Server: "https://kubernetes.default.svc"},
		}})
	}))
	defer server.Close()

	c := NewCluster(client)
	list, err := c.List(nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Items) != 1 || list.Items[0].Name != "in-cluster" {
		t.Errorf("unexpected list: %+v", list)
	}
}

func TestClusterList_WithOpts(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("server") != "https://example.com" {
			t.Errorf("expected server query param, got %q", q.Get("server"))
		}
		if q.Get("name") != "my-cluster" {
			t.Errorf("expected name query param, got %q", q.Get("name"))
		}
		jsonResponse(w, 200, ClusterList{Items: []*ClusterModel{
			{Name: "my-cluster", Server: "https://example.com"},
		}})
	}))
	defer server.Close()

	c := NewCluster(client)
	list, err := c.List(&ClusterQueryOptions{Server: "https://example.com", Name: "my-cluster"})
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Items) != 1 || list.Items[0].Name != "my-cluster" {
		t.Errorf("unexpected list: %+v", list)
	}
}

func TestClusterGet_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawPath == "" {
			t.Error("expected URL Path to be encoded, but RawPath is empty")
		}
		jsonResponse(w, 200, ClusterModel{Name: "in-cluster", Server: "https://kubernetes.default.svc"})
	}))
	defer server.Close()

	c := NewCluster(client)
	cluster, err := c.Get("https://kubernetes.default.svc", nil)
	if err != nil {
		t.Fatal(err)
	}
	if cluster.Name != "in-cluster" {
		t.Errorf("unexpected cluster: %+v", cluster)
	}
}

func TestClusterGet_WithOpts(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("name") != "my-cluster" {
			t.Errorf("expected name query param, got %q", q.Get("name"))
		}
		if q.Get("id.type") != "name" {
			t.Errorf("expected id.type query param, got %q", q.Get("id.type"))
		}
		jsonResponse(w, 200, ClusterModel{Name: "my-cluster", Server: "https://example.com"})
	}))
	defer server.Close()

	c := NewCluster(client)
	cluster, err := c.Get("https://example.com", &ClusterQueryOptions{Name: "my-cluster", IdType: "name"})
	if err != nil {
		t.Fatal(err)
	}
	if cluster.Name != "my-cluster" {
		t.Errorf("unexpected cluster: %+v", cluster)
	}
}

func TestClusterCreate_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var cluster ClusterModel
		json.NewDecoder(r.Body).Decode(&cluster)
		jsonResponse(w, 201, cluster)
	}))
	defer server.Close()

	c := NewCluster(client)
	cluster, err := c.Create(&ClusterModel{Name: "new-cluster", Server: "https://1.2.3.4"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if cluster.Name != "new-cluster" {
		t.Errorf("unexpected cluster: %+v", cluster)
	}
}
func TestClusterCreate_WithUpsert(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("upsert") != "true" {
			t.Errorf("expected upsert query param, got %q", r.URL.Query().Get("upsert"))
		}
		var cluster ClusterModel
		json.NewDecoder(r.Body).Decode(&cluster)
		jsonResponse(w, 201, cluster)
	}))
	defer server.Close()

	c := NewCluster(client)
	cluster, err := c.Create(&ClusterModel{Name: "new-cluster", Server: "https://1.2.3.4"}, &ClusterCreateOptions{Upsert: true})
	if err != nil {
		t.Fatal(err)
	}
	if cluster.Name != "new-cluster" {
		t.Errorf("unexpected cluster: %+v", cluster)
	}
}

func TestClusterUpdate_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawPath == "" {
			t.Error("expected URL Path to be encoded, but RawPath is empty")
		}
		jsonResponse(w, 200, ClusterModel{Name: "updated", Server: "https://1.2.3.4"})
	}))
	defer server.Close()

	c := NewCluster(client)
	cluster, err := c.Update(&ClusterModel{Name: "updated", Server: "https://1.2.3.4"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if cluster.Name != "updated" {
		t.Errorf("unexpected cluster: %+v", cluster)
	}
}
func TestClusterUpdate_WithUpdatedFields(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("updatedFields") != "name" {
			t.Errorf("expected updatedFields query param, got %q", r.URL.Query().Get("updatedFields"))
		}
		jsonResponse(w, 200, ClusterModel{Name: "updated", Server: "https://1.2.3.4"})
	}))
	defer server.Close()

	c := NewCluster(client)
	cluster, err := c.Update(&ClusterModel{Name: "updated", Server: "https://1.2.3.4"}, &ClusterUpdateOptions{UpdatedFields: []string{"name"}})
	if err != nil {
		t.Fatal(err)
	}
	if cluster.Name != "updated" {
		t.Errorf("unexpected cluster: %+v", cluster)
	}
}

func TestClusterDelete_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawPath == "" {
			t.Error("expected URL Path to be encoded, but RawPath is empty")
		}
		w.WriteHeader(200)
	}))
	defer server.Close()

	c := NewCluster(client)
	if err := c.Delete("https://kubernetes.default.svc", nil); err != nil {
		t.Fatal(err)
	}
}

func TestClusterDelete_WithOpts(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("name") != "my-cluster" {
			t.Errorf("expected name query param, got %q", q.Get("name"))
		}
		if q.Get("id.type") != "name" {
			t.Errorf("expected id.type query param, got %q", q.Get("id.type"))
		}
		w.WriteHeader(200)
	}))
	defer server.Close()

	c := NewCluster(client)
	if err := c.Delete("https://example.com", &ClusterQueryOptions{Name: "my-cluster", IdType: "name"}); err != nil {
		t.Fatal(err)
	}
}

func TestClusterRotateAuth_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawPath == "" {
			t.Error("expected URL Path to be encoded, but RawPath is empty")
		}
		w.WriteHeader(200)
	}))
	defer server.Close()

	c := NewCluster(client)
	if err := c.RotateAuth("https://kubernetes.default.svc", nil); err != nil {
		t.Fatal(err)
	}
}

func TestClusterRotateAuth_WithOpts(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("name") != "my-cluster" {
			t.Errorf("expected name query param, got %q", q.Get("name"))
		}
		if q.Get("id.type") != "name" {
			t.Errorf("expected id.type query param, got %q", q.Get("id.type"))
		}
		w.WriteHeader(200)
	}))
	defer server.Close()

	c := NewCluster(client)
	if err := c.RotateAuth("https://example.com", &ClusterQueryOptions{Name: "my-cluster", IdType: "name"}); err != nil {
		t.Fatal(err)
	}
}

func TestClusterInvalidateCache_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawPath == "" {
			t.Error("expected URL Path to be encoded, but RawPath is empty")
		}
		w.WriteHeader(200)
	}))
	defer server.Close()

	c := NewCluster(client)
	if err := c.InvalidateCache("https://kubernetes.default.svc", nil); err != nil {
		t.Fatal(err)
	}
}

func TestClusterInvalidateCache_WithOpts(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("name") != "my-cluster" {
			t.Errorf("expected name query param, got %q", q.Get("name"))
		}
		if q.Get("id.type") != "name" {
			t.Errorf("expected id.type query param, got %q", q.Get("id.type"))
		}
		w.WriteHeader(200)
	}))
	defer server.Close()

	c := NewCluster(client)
	if err := c.InvalidateCache("https://example.com", &ClusterQueryOptions{Name: "my-cluster", IdType: "name"}); err != nil {
		t.Fatal(err)
	}
}

func TestClusterGet_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 404, APIError{Code: 404, Message: "not found"})
	}))
	defer server.Close()

	c := NewCluster(client)
	_, err := c.Get("https://unknown", nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestClusterList_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	c := NewCluster(client)
	_, err := c.List(nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestClusterCreate_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	c := NewCluster(client)
	_, err := c.Create(&ClusterModel{Name: "new-cluster", Server: "https://1.2.3.4"}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestClusterUpdate_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	c := NewCluster(client)
	_, err := c.Update(&ClusterModel{Name: "updated", Server: "https://1.2.3.4"}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestClusterDelete_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	c := NewCluster(client)
	err := c.Delete("https://kubernetes.default.svc", nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestClusterRotateAuth_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	c := NewCluster(client)
	err := c.RotateAuth("https://kubernetes.default.svc", nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestClusterInvalidateCache_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	c := NewCluster(client)
	err := c.InvalidateCache("https://kubernetes.default.svc", nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestCluster_NetworkError(t *testing.T) {
	client := newFailingClient()
	c := NewCluster(client)
	_, err := c.List(nil)
	if err == nil {
		t.Error("expected error")
	}
	_, err = c.Get("https://kubernetes.default.svc", nil)
	if err == nil {
		t.Error("expected error")
	}
	_, err = c.Create(&ClusterModel{Name: "new-cluster", Server: "https://1.2.3.4"}, nil)
	if err == nil {
		t.Error("expected error")
	}
	_, err = c.Update(&ClusterModel{Name: "updated", Server: "https://1.2.3.4"}, nil)
	if err == nil {
		t.Error("expected error")
	}
	if err = c.Delete("https://kubernetes.default.svc", nil); err == nil {
		t.Error("expected error")
	}
	if err = c.RotateAuth("https://kubernetes.default.svc", nil); err == nil {
		t.Error("expected error")
	}
	if err = c.InvalidateCache("https://kubernetes.default.svc", nil); err == nil {
		t.Error("expected error")
	}
}

func TestEncodeClusterServer(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"https://kubernetes.default.svc", "https:%2F%2Fkubernetes.default.svc"},
		{"https://example.com:6443", "https:%2F%2Fexample.com:6443"},
		{"plain-string", "plain-string"},
		{"user@host/path", "user@host%2Fpath"},
	}
	for _, tc := range cases {
		got := encodeClusterServer(tc.in)
		if got != tc.want {
			t.Errorf("encodeClusterServer(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
