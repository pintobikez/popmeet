package structures

type SecurityConfig struct {
	CipherKey string `yaml:"cipherkey"`
	TTL       int    `yaml:"ttl"`
}

type DatabaseConfig struct {
	Host   string `yaml:"host,omitempty"`
	User   string `yaml:"user,omitempty"`
	Pw     string `yaml:"pw,omitempty"`
	Port   int    `yaml:"port,omitempty"`
	Schema string `yaml:"schema,omitempty"`
}
