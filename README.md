# Backup Service

Ce projet est un service de sauvegarde automatique pour plusieurs bases de données (MySQL et PostgreSQL). Il peut être exécuté en tant que service planifié avec `cron` et gère l'envoi des sauvegardes vers un serveur distant via SFTP.

## Structure du Projet

```bash
backup-service/
├── cmd/
│   └── main.go                 # Point d'entrée principal
├── config/
│   └── config.yaml             # Fichier de configuration YAML
├── logs/
│   └── backup.log              # Fichier de log
├── services/
│   ├── database.go             # Gestion des sauvegardes des bases de données
│   ├── files.go                # Gestion des fichiers locaux et archivage
│   ├── sftp.go                 # Gestion des transferts SFTP
│   ├── scheduler.go            # Gestionnaire de planification et multitâches
│   ├── cleanup.go              # Suppression des fichiers anciens
├── shared/
│   └── config.go               # Structures partagées
├── go.mod                      # Dépendances Go
├── go.sum                      # Sommes de contrôle des dépendances
├── README.md                   # Documentation du projet

```

## Installation

1. Clonez le projet :
```bash
git clone https://github.com/votre-repository/backup-service.git
cd backup-service
```

2. Installez les dépendances Go :
```bash
go mod tidy
```

3. Compilez le projet :
```bash
go build -o backup-service cmd/main.go
```

## Configuration

Modifiez `config/config.yaml` pour définir :
- Bases de données
- Répertoire des sauvegardes
- Paramètres du serveur SFTP
- Planning au format cron

## Exécution

### En arrière-plan :

#### Création d'un Service Linux :

Pour exécuter le binaire comme un service en arrière-plan, créez un fichier d'unité système pour `systemd`. Voici un exemple :
Fichier d'unité système (`/etc/systemd/system/backup-service.service`)

```bash
[Unit]
Description=Service de sauvegarde planifiée
After=network.target

[Service]
ExecStart=/chemin/vers/votre/binaire/backup-service -C /chemin/vers/config.yaml
Restart=always
User=votre_utilisateur
WorkingDirectory=/chemin/vers/votre/repertoire

[Install]
WantedBy=multi-user.target

```

- `ExecStart` : Spécifie le chemin du binaire et le chemin du fichier de configuration YAML.
- `Restart` : Redémarre le service automatiquement en cas de crash.
- `User` : L'utilisateur sous lequel le service s'exécute (veillez à utiliser un utilisateur avec des permissions appropriées).
- `WorkingDirectory` : Répertoire où se trouve votre binaire.


#### Activation et Lancement du Service :

1. Rechargez les fichiers de configuration `systemd` :
```bash
sudo systemctl daemon-reload
```

2. Activez le service pour qu'il démarre automatiquement au démarrage du système :
```bash
sudo systemctl enable backup-service
```

3. Lancez le service :
```bash
sudo systemctl start backup-service
```

4. Vérifiez l'état du service :
```bash
sudo systemctl status backup-service
```

### En mode manuel :

Lancez le programme avec :
```bash
./backup-service -C config/config.yaml
```

## Journalisation

Les logs sont enregistrés dans `logs/backup.log`.

Les logs peuvent être suivis avec la commande suivante (Si vous avez crée un service):
```bash
journalctl -u backup-service -f
```

## Dépendances

- `github.com/robfig/cron/v3`
- `github.com/pkg/sftp`
- `golang.org/x/crypto/ssh`
- `gopkg.in/yaml.v2`
