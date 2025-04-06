package shared

type DatabaseConfig struct {
	Name            string `yaml:"name"`
	Type            string `yaml:"type"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	RemoteDirectory string `yaml:"remote_directory"`
}

type FileConfig struct {
	Path            string `yaml:"path"`
	RemoteDirectory string `yaml:"remote_directory"`
}

type ServerConfig struct {
	Address         string `yaml:"address"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	RemoteDirectory string `yaml:"remote_directory"`
}
