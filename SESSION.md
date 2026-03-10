# Session-Stand: fehlerundstolzdrauf.de interaktiv

> Diese Datei dient zur Fortsetzung der Arbeit auf einem anderen Rechner.
> Nach Abschluss kann sie gelöscht werden.

## Was wurde gemacht

Die statische Zeitungs-Website wurde zu einer interaktiven Go-Anwendung umgebaut:

### Neue Dateien erstellt

| Datei | Inhalt |
|-------|--------|
| `go.mod` / `go.sum` | Go-Modul mit `mattn/go-sqlite3` Dependency |
| `cmd/server/main.go` | Entry Point, Routing, Static-File-Serving |
| `internal/db/db.go` | SQLite Init, Schema (5 Tabellen), Seed-Logik |
| `internal/api/public.go` | `GET /api/edition`, `POST /api/submit` mit Rate-Limiting |
| `internal/api/admin.go` | Admin-Endpunkte mit Bearer-Token-Auth |
| `internal/rotation/rotation.go` | Tages-Ausgabe: date-seeded Rotation, shown_count |
| `data/seed.sql` | 15 Stories, 3 Featured, 28 Quotes, 12 Historical |
| `.claude/skills/redaktion.md` | `/redaktion` Skill für wöchentliche Pflege |
| `.dockerignore` | Build-Context-Filter |

### Geänderte Dateien

| Datei | Änderung |
|-------|----------|
| `public/index.html` | Dynamische Inhalte via `fetch('/api/edition')`, Einreichungsformular mit Zeichenzähler, dynamisches Datum |
| `Dockerfile` | Multi-Stage Go Build (golang:1.22-alpine → alpine) mit CGO für sqlite3 |
| `.github/workflows/cicd.yml` | Trigger auf Branch `deployment` (nicht main), ADMIN_TOKEN env, bind mount `/home/wolf/fehlerundstolzdrauf/data:/data` |
| `info.md` | Vollständig aktualisiert mit neuer Architektur |

## Was noch offen ist

1. **Docker Build & Test** — Docker war auf dem lokalen Rechner nicht verfügbar. Erster Build+Test steht noch aus:
   ```bash
   docker build -t fsd . && docker run --rm -p 8080:80 -e ADMIN_TOKEN=test -v ./data:/data fsd
   ```
   Dann testen:
   - http://localhost:8080 — Seite laden
   - `curl http://localhost:8080/api/edition` — JSON
   - `curl -X POST -H "Content-Type: application/json" -d '{"year":"2024","title":"Test","text":"Eine Testgeschichte."}' http://localhost:8080/api/submit`
   - `curl -H "Authorization: Bearer test" http://localhost:8080/admin/stats`

2. **go.sum verifizieren** — Die go.sum wurde manuell geschrieben. Der Docker-Build läuft `go mod download && go mod verify`, sollte also ggf. fehlschlagen wenn der Hash nicht stimmt. Fix: im Container oder auf Rechner mit Go `go mod tidy` laufen lassen.

3. **GitHub Secret** — `ADMIN_TOKEN` muss als GitHub Secret im Repo angelegt werden.

4. **Branch `deployment` erstellen** — CI/CD triggert jetzt auf `deployment` statt `main`.

5. **Erster Deploy** — Nach erfolgreichem lokalen Test: Push auf `deployment` Branch.

## Architektur-Kurzübersicht

```
Browser → nginx-proxy → Go-Server (:80)
                          ├── GET /            → public/index.html (statisch)
                          ├── GET /api/edition → Tages-Ausgabe JSON (cached pro Tag)
                          ├── POST /api/submit → Einreichung (rate-limited, 1/h/IP)
                          └── /admin/*         → Token-geschützt (Bearer ADMIN_TOKEN)
```

- **DB**: SQLite in `/data/fehlerundstolzdrauf.db`
  - Lokal: `./data/` im Projektordner (bind mount)
  - Server: `/home/wolf/fehlerundstolzdrauf/data/` (bind mount)
- **Rotation**: SHA256(datum) als Seed, bevorzugt niedrigen shown_count
- **Moderation**: Einreichungen sind sofort `unmoderated` und sichtbar, werden nachträglich `approved` oder `hidden`
- **SQL Injection**: Alle Queries parametrisiert (`?`), Tabellennamen gegen Whitelist geprüft

## Hinweise

- Kein Framework — reines `net/http` + `mattn/go-sqlite3`
- CGO_ENABLED=1 nötig für sqlite3 (gcc + musl-dev im Dockerfile)
- Seed-Daten werden nur geladen wenn `stories`-Tabelle leer ist
- Edition-Cache wird bei Submit/Admin-Änderungen invalidiert
