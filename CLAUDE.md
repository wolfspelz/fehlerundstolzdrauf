# fehlerundstolzdrauf.de

> Alle Infos zu diesem Projekt. Bei Änderungen immer aktualisieren.

## Konzept

"Fehler & Stolz drauf" – Eine fiktive Zeitung/Zeitschrift im Hemingway-Stil.
Menschen erzählen anonym ihre größten Misserfolge und was daraus wurde.
Kein Name. Keine Entschuldigung. Nur was passiert ist. Und was danach kam.

Interaktiv: Besucher reichen Stories ein, Inhalte rotieren täglich, Pflege via Claude Code Skill `/redaktion`.

## Spalten / Sektionen

| Spalte | Inhalt | Stil |
|---|---|---|
| **Kurze Fälle** | Anonyme Kurzgeschichten über Scheitern (80-120 Wörter). Überschriften variieren im Muster | Hemingway: kurze Sätze, direkt, kein Pathos |
| **Diese Woche** | Featured Story, länger, mit Zitat-Block | Größere Schrift, Drop Cap, italic Titel |
| **Zahlen** | Statistiken + Einreichungsformular | Große Zahlen (Playfair Display), Live-Zähler |
| **Zitate** | Bekannte Zitate über Fehler/Scheitern | Special Elite Font (Schreibmaschine) |
| **Historisch** | Berühmte historische Fehlschläge/Zufallsentdeckungen | Italic Titel, Hemingway-Stil |

## Architektur

### Go-Backend (Single Container)

Go-Binary ersetzt nginx. Serviert statische Dateien UND API.

```
Browser  →  nginx-proxy (besteht)  →  Go-Server (:80)
                                        ├── GET /           → public/index.html
                                        ├── GET /api/edition → Tages-Ausgabe JSON
                                        ├── POST /api/submit → Besucher-Einreichung
                                        └── Admin-Endpunkte (Token-geschützt)
```

### Tages-Ausgabe (Rotation)

- `GET /api/edition` liefert täglich wechselnde Inhalte
- Date-seeded: Alle Besucher sehen am selben Tag dasselbe
- Bevorzugt Content mit niedrigstem `shown_count`
- Pro Tag: 5 Stories, 1 Featured, 3 Quotes, 2 Historical
- Cache in `edition_cache` Tabelle

### Einreichung

- `POST /api/submit` mit `{year, title, text}`
- Rate-Limit: 1 pro IP pro Stunde
- Status `unmoderated` → sofort sichtbar
- Moderation nachträglich via Admin-API

### Admin-API

Header: `Authorization: Bearer <ADMIN_TOKEN>`

| Method | Path | Zweck |
|---|---|---|
| GET | `/admin/submissions?status=unmoderated` | Einreichungen filtern |
| PUT | `/admin/submissions/:id` | Status ändern (approved/hidden) |
| POST | `/admin/stories` | Neue Story (status=approved) |
| POST | `/admin/featured` | Neue Featured Story |
| POST | `/admin/quotes` | Neues Zitat |
| POST | `/admin/historical` | Neuer Historisch-Eintrag |
| GET | `/admin/stats` | Übersicht |
| DELETE | `/admin/:type/:id` | Löschen |

### SQLite

- Datenbank: `/data/fehlerundstolzdrauf.db` (Docker Volume)
- Tabellen: `stories`, `featured`, `quotes`, `historical`, `edition_cache`
- Seed-Daten in `data/seed.sql`

## Design

- **Zeitungs-Layout**: 4 Spalten auf Desktop, 1 Spalte unter 900px
- **Papierfarbe**: `#f5f0e8` (warm off-white)
- **Hintergrund**: `#ffffff` (weiß, außerhalb der "Zeitung")
- **Schriftarten**:
  - `Playfair Display` – Überschriften, Fließtext, große Zahlen
  - `Courier Prime` – Labels, Jahreszahlen, Dateline, Footer, Formular
  - `Special Elite` – Zitate (Schreibmaschinen-Look)
- **Farben**: `--black: #1a1612`, `--red: #8b1a1a`, `--grey: #6b6560`, `--ink: #2a2520`
- **Effekte**: Papier-Grain (SVG Turbulence), Box-Shadow, col-label mit Doppellinie

## Projektstruktur

```
cmd/server/main.go           -- Entry Point, Routes, Server
internal/
  db/db.go                   -- SQLite Init, Schema, Seed
  api/public.go              -- GET /api/edition, POST /api/submit
  api/admin.go               -- Admin-Endpunkte mit Token-Auth
  rotation/rotation.go       -- Tages-Ausgabe Auswahl-Logik
public/index.html            -- Frontend mit JS fetch + Formular
data/seed.sql                -- Initiale Daten
Dockerfile                   -- Multi-Stage Go Build
.github/workflows/cicd.yml   -- CI/CD Pipeline
.claude/skills/redaktion.md  -- Redaktions-Skill
```

## Technischer Stack

- **Repo**: https://github.com/wolfspelz/fehlerundstolzdrauf
- **Backend**: Go (net/http + mattn/go-sqlite3)
- **Storage**: SQLite mit Docker Volume
- **Frontend**: Einzelne HTML-Datei mit Inline-JS, kein Build-System
- **Docker**: Multi-Stage Go Build → Alpine
- **Image**: `wolfspelz/fehlerundstolzdrauf` auf Docker Hub

## CI/CD

- **Workflow**: `.github/workflows/cicd.yml`
- **Trigger**: Push auf `main`
- **Ablauf**: Checkout → Docker Login → Build & Push Image → SSH Deploy auf Server
- **GitHub Secrets** (im Repo konfiguriert):
  - `DOCKERHUB_USERNAME`
  - `DOCKERHUB_TOKEN`
  - `SERVER_HOSTNAME`
  - `SERVER_USERNAME`
  - `SERVER_PASSWORD`
  - `ADMIN_TOKEN` – Token für Admin-API

## Redaktions-Skill `/redaktion`

Wöchentliche Pflege via Claude Code:
1. **Moderieren**: Einreichungen freigeben/verbergen
2. **Zitate generieren**: Neue Zitate einfügen
3. **Historisch generieren**: Neue historische Stories
4. **Stories schreiben**: Kurzgeschichten im Hemingway-Stil
5. **Featured schreiben**: Längere Featured Stories
6. **Statistik**: Übersicht Content-Mengen

## Deployment

- Kein Staging-Environment – Push auf `main` geht direkt in Produktion
- Server läuft hinter einem Reverse Proxy mit Let's Encrypt (automatisches HTTPS)
- SQLite-Datenbank:
  - **Lokal**: `data/fehlerundstolzdrauf.db` (im Projektordner, via `DB_PATH=data/fehlerundstolzdrauf.db`)
  - **Server**: `/home/wolf/fehlerundstolzdrauf/data/fehlerundstolzdrauf.db` (bind mount nach `/data` im Container)

## Workflow für Änderungen

1. Code bearbeiten
2. Lokal prüfen: `docker build -t fsd . && docker run -p 8080:80 -e ADMIN_TOKEN=test -v ./data:/data fsd`
3. `git add && git commit && git push`
4. CI/CD deployt automatisch auf Server
5. Live prüfen: https://fehlerundstolzdrauf.de
