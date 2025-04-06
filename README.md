# Backup Service

Backup Service est une application Go permettant de sauvegarder des bases de données et des fichiers locaux, avec des options de transfert vers un serveur distant via SFTP. Elle inclut également une planification `cron` et une gestion des fichiers anciens.

## Fonctionnalités

- **Bases de données** :
  - Sauvegarde des bases de données MySQL et PostgreSQL.
  - Connexion aux serveurs via adresse (nom d'hôte ou IP).
- **Fichiers locaux** :
  - Archivage et transfert des fichiers/dossiers locaux.
- **Planification** :
  - Sauvegardes automatiques basées sur un horaire configuré.
- **Nettoyage** :
  - Option de suppression automatique des anciens fichiers après un certain nombre de jours.
- **Transfert SFTP** :
  - Envoi des sauvegardes vers un serveur distant (facultatif).

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

## Prérequis

1. **Go** : Version 1.19 ou supérieure.
2. **Commandes externes** :
   - `mysqldump` pour les bases de données MySQL.
   - `pg_dump` pour les bases de données PostgreSQL.
3. **Dépendances Go** : Installez les dépendances via `go mod tidy`.

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

#### Lancer le Service

Lancez le service avec le fichier de configuration spécifié :
```bash
./backup-service -C config/config.yaml
```
#### Exécution Immédiate

Pour exécuter immédiatement sans planification, utilisez le drapeau --run-now :
```bash
./backup-service --run-now -C config/config.yaml
```

## Exemple de Configuration

Voici un exemple de configuration pour `config.yaml` :
```yaml
databases:
  - name: "example_db"
    type: "mysql"
    user: "root"
    password: "example_password"
    address: "localhost"
    remote_directory: "/remote/mysql_backups/"
files:
  - path: "/path/to/local/file.txt"
    remote_directory: "/remote/files/"
backup:
  directory: "/local/backups/"
  retention_days: 7
  max_concurrency: 3
  enable_cleanup: true
server:
  address: "sftp.example.com:22"
  user: "example_user"
  password: "example_password"
schedule: "0 3 * * *"
```

## Nettoyage des Fichiers Anciens

La suppression des fichiers anciens est contrôlée par le champ `enable_cleanup`. Par défaut, elle est activée (`true`). Pour désactiver cette option, modifiez la configuration YAML :

```yaml
backup:
  enable_cleanup: false
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
