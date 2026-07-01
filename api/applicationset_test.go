package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestApplicationSetList_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, struct {
			Items []*ApplicationSetModel `json:"items"`
		}{Items: []*ApplicationSetModel{
			{ObjectMeta: ObjectMeta{Name: "myappset"}},
		}})
	}))
	defer server.Close()

	a := NewApplicationSet(client)
	list, err := a.List(nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].Name != "myappset" {
		t.Errorf("unexpected list: %+v", list)
	}
}

func TestApplicationSetGet_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, ApplicationSetModel{ObjectMeta: ObjectMeta{Name: "myappset"}})
	}))
	defer server.Close()

	a := NewApplicationSet(client)
	appset, err := a.Get("myappset", nil)
	if err != nil {
		t.Fatal(err)
	}
	if appset.Name != "myappset" {
		t.Errorf("unexpected appset: %+v", appset)
	}
}

func TestApplicationSetCreate_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var appset ApplicationSetModel
		json.NewDecoder(r.Body).Decode(&appset)
		jsonResponse(w, 201, appset)
	}))
	defer server.Close()

	a := NewApplicationSet(client)
	appset, err := a.Create(&ApplicationSetModel{ObjectMeta: ObjectMeta{Name: "new"}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if appset.Name != "new" {
		t.Errorf("unexpected appset: %+v", appset)
	}
}

func TestApplicationSetDelete_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	a := NewApplicationSet(client)
	if err := a.Delete("myappset", nil); err != nil {
		t.Fatal(err)
	}
}

func TestApplicationSetGenerate_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, struct {
			Items []*ApplicationSetModel `json:"items"`
		}{Items: []*ApplicationSetModel{
			{ObjectMeta: ObjectMeta{Name: "generated-app"}},
		}})
	}))
	defer server.Close()

	a := NewApplicationSet(client)
	generated, err := a.Generate(&ApplicationSetModel{ObjectMeta: ObjectMeta{Name: "generator"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(generated) != 1 || generated[0].Name != "generated-app" {
		t.Errorf("unexpected generated: %+v", generated)
	}
}

func TestApplicationSetResourceTree_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, ApplicationTree{Nodes: []ResourceNode{{Name: "pod1"}}})
	}))
	defer server.Close()

	a := NewApplicationSet(client)
	tree, err := a.ResourceTree("myappset", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(tree.Nodes) != 1 || tree.Nodes[0].Name != "pod1" {
		t.Errorf("unexpected tree: %+v", tree)
	}
}

func TestApplicationSetListResourceEvents_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 200, ResourceEventList{Items: []ResourceEvent{{Reason: "Created"}}})
	}))
	defer server.Close()

	a := NewApplicationSet(client)
	events, err := a.ListResourceEvents("myappset", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(events.Items) != 1 {
		t.Errorf("unexpected events: %+v", events)
	}
}

func TestApplicationSetList_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	a := NewApplicationSet(client)
	_, err := a.List(nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationSetDelete_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 404, APIError{Code: 404, Message: "not found"})
	}))
	defer server.Close()

	a := NewApplicationSet(client)
	if err := a.Delete("nonexistent", nil); err == nil {
		t.Error("expected error")
	}
}

func TestApplicationSetGet_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	a := NewApplicationSet(client)
	_, err := a.Get("myappset", nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationSetCreate_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	a := NewApplicationSet(client)
	_, err := a.Create(&ApplicationSetModel{ObjectMeta: ObjectMeta{Name: "new"}}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationSetGenerate_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	a := NewApplicationSet(client)
	_, err := a.Generate(&ApplicationSetModel{ObjectMeta: ObjectMeta{Name: "generator"}})
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationSetResourceTree_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	a := NewApplicationSet(client)
	_, err := a.ResourceTree("myappset", nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationSetListResourceEvents_Error(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(w, 500, APIError{Code: 500, Message: "error"})
	}))
	defer server.Close()

	a := NewApplicationSet(client)
	_, err := a.ListResourceEvents("myappset", nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationSet_NetworkError(t *testing.T) {
	client := newFailingClient()
	a := NewApplicationSet(client)
	_, err := a.List(nil)
	if err == nil {
		t.Error("expected error")
	}
	_, err = a.Get("myappset", nil)
	if err == nil {
		t.Error("expected error")
	}
	_, err = a.Create(&ApplicationSetModel{ObjectMeta: ObjectMeta{Name: "new"}}, nil)
	if err == nil {
		t.Error("expected error")
	}
	if err = a.Delete("myappset", nil); err == nil {
		t.Error("expected error")
	}
	_, err = a.Generate(&ApplicationSetModel{ObjectMeta: ObjectMeta{Name: "generator"}})
	if err == nil {
		t.Error("expected error")
	}
	_, err = a.ResourceTree("myappset", nil)
	if err == nil {
		t.Error("expected error")
	}
	_, err = a.ListResourceEvents("myappset", nil)
	if err == nil {
		t.Error("expected error")
	}
	_, err = a.Watch(context.Background(), nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestApplicationSetWatch_Success(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(200)
		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Error("expected http.Flusher")
			return
		}
		data, _ := json.Marshal(struct {
			Result ApplicationSetWatchEvent `json:"result"`
		}{
			Result: ApplicationSetWatchEvent{
				Type:           SyncStatusCodeSynced,
				ApplicationSet: &ApplicationSetModel{ObjectMeta: ObjectMeta{Name: "myappset"}},
			},
		})
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	a := NewApplicationSet(client)
	ch, err := a.Watch(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	select {
	case event, ok := <-ch:
		if !ok {
			t.Error("channel closed unexpectedly")
			return
		}
		if event.ApplicationSet.Name != "myappset" {
			t.Errorf("unexpected event: %+v", event)
		}
		cancel()
	case <-time.After(2 * time.Second):
		t.Error("timed out waiting for watch event")
	}
}

func TestApplicationSetWatch_SSEError(t *testing.T) {
	client, server := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("internal error"))
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	a := NewApplicationSet(client)
	_, err := a.Watch(ctx, nil)
	if err == nil {
		t.Error("expected error")
	}
}
