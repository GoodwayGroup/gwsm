package lib

type ConfigMap struct {
	Data map[string]string
}

type Result struct {
	Name  string
	JSON  map[string]interface{}
	Error error
}

type KubeSecret struct {
	Kind string
	Type string
	Data map[string]interface{}
}
