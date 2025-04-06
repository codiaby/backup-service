package shared

type DatabaseConfig struct {
	Name            string `yaml:"name"`             // Nom de la base de données
	Type            string `yaml:"type"`             // Type de base de données (MySQL, PostgreSQL, etc.)
	User            string `yaml:"user"`             // Nom d'utilisateur pour la connexion à la base de données
	Password        string `yaml:"password"`         // Mot de passe pour la connexion à la base de données
	Address         string `yaml:"address"`          // Adresse de la base de données (ex: localhost:3306)
	RemoteDirectory string `yaml:"remote_directory"` // Répertoire distant pour stocker la sauvegarde
}

type FileConfig struct {
	Path            string `yaml:"path"`             // Chemin du fichier ou du répertoire à sauvegarder
	RemoteDirectory string `yaml:"remote_directory"` // Répertoire distant pour stocker la sauvegarde
}

type ServerConfig struct {
	Address         string `yaml:"address"`          // Adresse IP ou nom d'hôte du serveur distant
	User            string `yaml:"user"`             // Nom d'utilisateur pour la connexion SSH
	Password        string `yaml:"password"`         // Mot de passe pour la connexion SSH
	RemoteDirectory string `yaml:"remote_directory"` // Répertoire distant pour stocker la sauvegarde
}
