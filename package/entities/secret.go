package entities

type Secret struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Type       string `yaml:"type"`
}

type Data struct {
	Key   string
	Value string
}
