# Plan: fehlerundstolzdrauf.de interaktiv machen

## Kontext

Die statische Zeitungs-Website soll interaktiv werden: Besucher reichen Stories ein, Inhalte rotieren täglich, und der Betreiber pflegt alles wöchentlich per Claude Code Skill (Moderation, neue Zitate/Historisch, eigene Beiträge schreiben). Kein Claude API im Backend.

## Entscheidungen

- **Backend**: Go (single binary, kein Runtime)
- **Storage**: SQLite mit Host-Dateisystem (bind mount, kein Docker Volume)
- **Rotation**: Tages-Ausgabe (date-seeded, alle sehen dasselbe)
- **Formular**: Immer offen (Rate-Limiting gegen Spam)
- **Moderation**: Nachträglich via Claude Code Skill. Beiträge sind sofort live, werden bei Moderation freigegeben oder verborgen
- **Content-Generierung**: Manuell via Claude Code Skill
- **Deployment**: Push auf Branch `deployment` (nicht main)
- **SQL Injection**: Alle Queries parametrisiert, keine Datenbank-Escaping-Abhängigkeit

---

## Architektur

### Single Container

Go-Binary ersetzt nginx. Serviert statische Dateien UND API.

```
Browser  →  nginx-proxy (besteht)  →  Go-Server (:80)
                                        ├── GET /           → public/index.html
                                        ├── GET /api/edition → Tages-Ausgabe JSON
                                        ├── POST /api/submit → Besucher-Einreichung
                                        └── Admin-Endpunkte (Token-geschützt)
```

### Admin-API (Token-geschützt)

Header: `Authorization: Bearer <ADMIN_TOKEN>`

| Method | Path | Zweck |
|--------|------|-------|
| GET | `/admin/submissions?status=unmoderated` | Unmoderierte Einreichungen anzeigen |
| PUT | `/admin/submissions/:id` | Status ändern (approved/hidden) |
| POST | `/admin/stories` | Neue Kurzgeschichte anlegen |
| POST | `/admin/featured` | Neue Featured Story anlegen |
| POST | `/admin/quotes` | Neues Zitat anlegen |
| POST | `/admin/historical` | Neuen Historisch-Eintrag anlegen |
| GET | `/admin/stats` | Übersicht: Anzahl Stories, Pending, etc. |
| DELETE | `/admin/:type/:id` | Eintrag löschen |

`ADMIN_TOKEN` wird als Env-Variable übergeben und als GitHub Secret gespeichert.

### SQLite Schema

```sql
CREATE TABLE stories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    year TEXT NOT NULL,
    title TEXT NOT NULL,
    text TEXT NOT NULL,
    created_at TEXT DEFAULT (datetime('now')),
    status TEXT DEFAULT 'unmoderated',
    shown_count INTEGER DEFAULT 0
);

CREATE TABLE featured (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    year_range TEXT NOT NULL,
    title TEXT NOT NULL,
    intro TEXT NOT NULL,
    quote TEXT,
    outro TEXT,
    created_at TEXT DEFAULT (datetime('now')),
    shown_count INTEGER DEFAULT 0
);

CREATE TABLE quotes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    text TEXT NOT NULL,
    attribution TEXT NOT NULL,
    shown_count INTEGER DEFAULT 0
);

CREATE TABLE historical (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    year TEXT NOT NULL,
    title TEXT NOT NULL,
    text TEXT NOT NULL,
    shown_count INTEGER DEFAULT 0
);

CREATE TABLE edition_cache (
    date TEXT PRIMARY KEY,
    content_json TEXT NOT NULL
);
```

### Tages-Ausgabe Logik

1. `GET /api/edition` prüft `edition_cache` für heute
2. Falls nicht vorhanden: wähle per date-seed aus approved Content:
   - 5 Stories mit status `unmoderated` oder `approved` (niedrigster `shown_count` bevorzugt)
   - 1 Featured
   - 3 Quotes
   - 2 Historical
3. Erhöhe `shown_count`, speichere in Cache
4. Gibt JSON zurück

Fallback: Falls nicht genug approved Stories vorhanden, fülle mit Seed-Daten auf.

### Besucher-Einreichung

- `POST /api/submit` mit `{year: "2009", title: "Die Bäckerei", text: "..."}`
- Validierung: year/title/text nicht leer, text max 500 Zeichen, kein HTML
- Rate-Limit: 1 pro IP pro Stunde (in-memory Map)
- Status: `unmoderated` → sofort sichtbar für Besucher
- Bei nachträglicher Moderation: `approved` (bleibt sichtbar) oder `hidden` (nie mehr angezeigt)
- Rotation nutzt sowohl `unmoderated` als auch `approved` Stories
- Antwort: Danke-Nachricht

---

## Projektstruktur

```
cmd/server/main.go           -- Entry Point, Routes, Server
internal/
  db/db.go                   -- SQLite Init, Migrations, Queries
  api/public.go              -- GET /api/edition, POST /api/submit
  api/admin.go               -- Admin-Endpunkte
  rotation/rotation.go       -- Tages-Ausgabe Auswahl-Logik
public/index.html            -- Modifiziert: JS fetch + Formular
data/seed.sql                -- Initiale Quotes + Historical + Stories
Dockerfile
go.mod
go.sum
.github/workflows/cicd.yml
.dockerignore
info.md
.claude/skills/redaktion.md
```

## Frontend-Änderungen (public/index.html)

1. **Content-Bereiche** werden zu leeren Containern mit IDs (`id="kurze-faelle"`, etc.)
2. **Einreichungs-Formular** ersetzt den "Einreichen"-Text in Spalte 3:
   - Jahr-Feld (kurzes Input)
   - Titel-Feld (kurzes Input)
   - `<textarea>` für Text (Courier Prime, max 500 Zeichen)
   - Zeichenzähler
   - Submit-Button im Zeitungs-Stil
   - Erfolgs-/Fehlermeldung
3. **Inline JavaScript** (~80 Zeilen, kein Build):
   - `fetch('/api/edition')` → DOM populieren
   - Formular-Submit → `fetch('/api/submit')`
   - Zeichenzähler
4. **Dateline** dynamisch: aktuelles Datum statt hardcoded "März 2026"
5. **"Zahlen"** Sektion: Gesamtzahl Stories wird dynamisch aus `/api/edition` geladen
6. Alle CSS bleiben unverändert

## Claude Code Skill: `/redaktion`

Ein Skill in `.claude/skills/redaktion.md` für die wöchentliche Pflege:

**Funktionen:**
1. **Moderieren**: Unmoderierte Einreichungen abrufen, jeweils freigeben (approved) oder verbergen (hidden)
2. **Zitate generieren**: 3-5 neue Zitate im Stil der bestehenden generieren, per Admin-API einfügen
3. **Historisch generieren**: 1-2 neue historische Fehlschlag-Stories generieren, einfügen
4. **Eigene Stories schreiben**: Interaktiv neue Kurzgeschichten im Hemingway-Stil entwickeln, einfügen
5. **Featured Story schreiben**: Längere Featured Story mit Intro, Zitat, Outro
6. **Statistik**: Übersicht über Anzahl Content pro Typ, Pending, etc.

Der Skill nutzt `curl` via Bash-Tool gegen die Admin-API mit dem Token.

## Dockerfile

```dockerfile
FROM golang:1.22-alpine AS build
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o server ./cmd/server

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=build /app/server /server
COPY --from=build /app/public /public
COPY --from=build /app/data /data
EXPOSE 80
CMD ["/server"]
```

## CI/CD

Trigger: Push auf Branch `deployment` (nicht main).

Deploy-Kommando in `cicd.yml`:
```yaml
mkdir -p /home/wolf/fehlerundstolzdrauf/data
docker run -d --name fehlerundstolzdrauf-prod \
  --restart unless-stopped \
  -e LETSENCRYPT_HOST=fehlerundstolzdrauf.de \
  -e VIRTUAL_HOST=fehlerundstolzdrauf.de \
  -e VIRTUAL_PORT=80 \
  -e ADMIN_TOKEN=${{ secrets.ADMIN_TOKEN }} \
  -v /home/wolf/fehlerundstolzdrauf/data:/data \
  --expose=80 \
  --network=web \
  -it wolfspelz/fehlerundstolzdrauf
```

Neue GitHub Secrets: `ADMIN_TOKEN`

## Datenbank-Pfade

- **Lokal**: `./data/fehlerundstolzdrauf.db` im Projektordner (bind mount `./data:/data`)
- **Server**: `/home/wolf/fehlerundstolzdrauf/data/fehlerundstolzdrauf.db` (bind mount nach `/data` im Container)
- Default im Go-Code: `/data/fehlerundstolzdrauf.db`, überschreibbar via `DB_PATH` env

## Seed-Daten

`data/seed.sql` enthält:
- 15 Kurzgeschichten (5 Original + 10 neue, alle status=approved)
- 3 Featured Stories
- 28 Zitate
- 12 Historisch-Einträge

## Verifizierung

1. `docker build -t fsd .` — kompiliert ohne Fehler
2. `docker run --rm -p 8080:80 -e ADMIN_TOKEN=test -v ./data:/data fsd` — Container startet
3. Browser: `http://localhost:8080` — Seite lädt, Inhalte erscheinen
4. `curl http://localhost:8080/api/edition` — JSON mit Tages-Ausgabe
5. `curl -X POST -H "Content-Type: application/json" -d '{"year":"2024","title":"Test","text":"Eine Testgeschichte."}' http://localhost:8080/api/submit` — Einreichung funktioniert
6. `curl -H "Authorization: Bearer test" http://localhost:8080/admin/stats` — Admin-API funktioniert
7. `/redaktion` Skill in Claude Code testen
