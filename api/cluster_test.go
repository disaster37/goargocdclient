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

func TestClusterGet_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

func TestClusterUpdate_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

func TestClusterDelete_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	c := NewCluster(client)
	if err := c.Delete("https://kubernetes.default.svc", nil); err != nil {
		t.Fatal(err)
	}
}

func TestClusterRotateAuth_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	c := NewCluster(client)
	if err := c.RotateAuth("https://kubernetes.default.svc", nil); err != nil {
		t.Fatal(err)
	}
}

func TestClusterInvalidateCache_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	c := NewCluster(client)
	if err := c.InvalidateCache("https://kubernetes.default.svc", nil); err != nil {
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
