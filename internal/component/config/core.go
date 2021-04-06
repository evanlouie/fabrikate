package config

type Config struct {
	Path            string                 `yaml:"-" json:"-"`
	Serialization   string                 `yaml:"-" json:"-"`
	Namespace       string                 `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	InjectNamespace bool                   `yaml:"injectNamespace,omitempty" json:"injectNamespace,omitempty"`
	Disabled        bool                   `yaml:"disabled,omitempty" json:"disabled,omitempty"`
	Config          map[string]interface{} `yaml:"config,omitempty" json:"config,omitempty"`
	Subcomponents   map[string]Config      `yaml:"subcomponents,omitempty" json:"subcomponents,omitempty"`
}
