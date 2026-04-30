# job-match-api

API Go qui analyse un CV en texte et une offre d'emploi via un LLM, et retourne un score de matching, un résumé du profil, les compétences communes, les compétences manquantes, et des recommandations actionnables.

## Stack

- **Go** 1.25
- **Gin** — framework HTTP
- **PostgreSQL** — base de données (via `pgx`)
- **OpenRouter** — accès aux modèles LLM (modèle utilisé : `openai/gpt-oss-120b:free`)
- **Docker Compose** — orchestration locale de la DB
- **Air** — live reload en développement

## Architecture

```
.
├── cmd/api/main.go              # Entrée du serveur HTTP + graceful shutdown
├── internal/
│   ├── server/                  # Configuration serveur + routes
│   ├── llm/                     # Appel OpenRouter + parsing de la réponse
│   ├── model/                   # Structures partagées (request/response)
│   └── database/                # Service PostgreSQL
├── docker-compose.yml           # Service PostgreSQL
├── Makefile                     # Commandes courantes
└── .env                         # Variables d'environnement (non versionné)
```

## Prérequis

- Go 1.25+
- Docker & Docker Compose
- Une clé API OpenRouter (https://openrouter.ai/)

## Installation

```bash
git clone <repo-url>
cd job_match_api
go mod download
```

## Configuration

Créer un fichier `.env` à la racine :

```env
PORT=8080
APP_ENV=local

BLUEPRINT_DB_HOST=localhost
BLUEPRINT_DB_PORT=5432
BLUEPRINT_DB_DATABASE=blueprint
BLUEPRINT_DB_USERNAME=melkey
BLUEPRINT_DB_PASSWORD=password1234
BLUEPRINT_DB_SCHEMA=public

OPENROUTER_API_KEY=sk-or-v1-...
```

## Démarrage

Lancer la base de données :
```bash
make docker-run
```

Lancer l'API :
```bash
make run
```

Mode live reload (avec `air`) :
```bash
make watch
```

L'API écoute sur `http://localhost:8080`.

## Endpoints

### `GET /`
Renvoie un message de bienvenue.

### `GET /health`
État de santé de la base de données.

### `POST /v1/cv/analyze`
Analyse un CV face à une description de poste.

**Requête** :
```json
{
  "job_description": "Backend Go developer with PostgreSQL and Docker experience",
  "cv_text": "Jean Dupont, backend developer, Go, REST APIs, Docker..."
}
```

**Réponse `200`** :
```json
{
  "match_score": 78,
  "summary": "Backend Go developer with strong API and Docker experience.",
  "matched_skills": ["Go", "REST", "Docker"],
  "missing_skills": ["Kubernetes", "OpenTelemetry"],
  "recommendations": [
    "Ajoute un projet avec observabilité.",
    "Mentionne la CI/CD.",
    "Mets en avant ton expérience PostgreSQL."
  ]
}
```

**Codes d'erreur** :
- `400` — requête invalide (champ manquant ou vide)
- `422` — contenu impossible à analyser
- `500` — erreur serveur ou échec de l'appel LLM

**Format d'erreur** :
```json
{
  "error": {
    "code": "INVALID_REQUEST",
    "message": "job_description is required"
  }
}
```

## Exemple d'appel

```bash
curl -X POST http://localhost:8080/v1/cv/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "job_description": "Backend Go developer with PostgreSQL and Docker experience",
    "cv_text": "Jean Dupont, backend developer, Go, REST APIs, Docker..."
  }'
```

## Tests

Tous les tests :
```bash
make test
```

Tests d'intégration DB uniquement :
```bash
make itest
```

> Note : `TestAnalyze` appelle réellement OpenRouter et nécessite une clé API valide.

## Commandes Makefile

| Commande | Description |
|---|---|
| `make build` | Compile le binaire dans `./main` |
| `make run` | Lance l'API |
| `make watch` | Live reload via `air` |
| `make test` | Lance la suite de tests |
| `make itest` | Tests d'intégration de la DB |
| `make docker-run` | Démarre PostgreSQL via Docker Compose |
| `make docker-down` | Arrête PostgreSQL |
| `make clean` | Supprime le binaire compilé |

