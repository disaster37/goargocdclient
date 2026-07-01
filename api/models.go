package api

type ObjectMeta struct {
	Name              string            `json:"name,omitempty"`
	Namespace         string            `json:"namespace,omitempty"`
	Labels            map[string]string `json:"labels,omitempty"`
	Annotations       map[string]string `json:"annotations,omitempty"`
	UID               string            `json:"uid,omitempty"`
	ResourceVersion   string            `json:"resourceVersion,omitempty"`
	CreationTimestamp string            `json:"creationTimestamp,omitempty"`
	DeletionTimestamp string            `json:"deletionTimestamp,omitempty"`
	Generation        int64             `json:"generation,omitempty"`
}

type TypeMeta struct {
	Kind       string `json:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty"`
}

type ListMeta struct {
	ResourceVersion string `json:"resourceVersion,omitempty"`
	SelfLink        string `json:"selfLink,omitempty"`
}