# fehlerundstolzdrauf.de

> Alle Infos zu diesem Projekt. Bei Änderungen immer aktualisieren.

## Konzept

"Fehler & Stolz drauf" – Eine fiktive Zeitung/Zeitschrift im Hemingway-Stil.
Menschen erzählen anonym ihre größten Misserfolge und was daraus wurde.
Kein Name. Keine Entschuldigung. Nur was passiert ist. Und was danach kam.

## Spalten / Sektionen

| Spalte | Inhalt | Stil |
|---|---|---|
| **Kurze Fälle** | Anonyme Kurzgeschichten über Scheitern (80-120 Wörter) | Hemingway: kurze Sätze, direkt, kein Pathos |
| **Diese Woche** | Featured Story, länger, mit Zitat-Block | Größere Schrift, Drop Cap, italic Titel |
| **Zahlen** | Statistiken + Einreichen-Hinweis | Große Zahlen (Playfair Display), kurze Erklärung |
| **Zitate** | Bekannte Zitate über Fehler/Scheitern (keine Privatpersonen) | Special Elite Font (Schreibmaschine) |
| **Historisch** | Berühmte historische Fehlschläge/Zufallsentdeckungen | Italic Titel, Hemingway-Stil |

## Design

- **Zeitungs-Layout**: 4 Spalten auf Desktop, 1 Spalte unter 900px
- **Papierfarbe**: `#f5f0e8` (warm off-white)
- **Hintergrund**: `#ffffff` (weiß, außerhalb der "Zeitung")
- **Schriftarten**:
  - `Playfair Display` – Überschriften, Fließtext, große Zahlen
  - `Courier Prime` – Labels, Jahreszahlen, Dateline, Footer (Schreibmaschine)
  - `Special Elite` – Zitate (Schreibmaschinen-Look)
- **Farben**: `--black: #1a1612`, `--red: #8b1a1a`, `--grey: #6b6560`, `--ink: #2a2520`
- **Effekte**: Papier-Grain (SVG Turbulence), Box-Shadow, col-label mit Doppellinie (2px + 1px ::after)
- **Header**: Zentriert, Ornamente (✦), Dateline oben
- **Footer**: Domain-Link, Trennzeichen, Motto
- **Drop Cap**: Erster Buchstabe der Featured Story in Rot

## Technischer Stack

- **Repo**: https://github.com/wolfspelz/fehlerundstolzdrauf
- **Technik**: Einzelne statische HTML-Datei (`public/index.html`), kein Build-System
- **Docker**: `nginx:alpine`, served aus `public/`
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

## CSS-Klassen

| Klasse | Verwendung |
|---|---|
| `.entry` | Container für eine Geschichte |
| `.entry-year` | Jahreszahl |
| `.entry-title` | Überschrift einer Geschichte |
| `.entry-text` | Fließtext |
| `.quote-entry` | Zitat-Container (Spalte Zitate) |
| `.quote-block` | Inline-Zitat innerhalb einer Featured Story |
| `.big-number` | Große Statistik-Zahlen |
| `.featured` | Featured-Spalte (größere Schrift, Drop Cap) |

## Print

- Print-Stylesheet: `@page { size: A3 landscape; }`

## Deployment

- Kein Staging-Environment – Push auf `main` geht direkt in Produktion
- Server läuft hinter einem Reverse Proxy mit Let's Encrypt (automatisches HTTPS)

## Workflow für Änderungen

1. HTML in `public/index.html` bearbeiten
2. Lokal prüfen:
   - Direkt im Browser: `public/index.html` öffnen
   - Oder per Docker: `docker build -t fehlerundstolzdrauf . && docker run -p 8080:80 fehlerundstolzdrauf`
3. `git add -A && git commit -m "Beschreibung" && git push`
4. CI/CD deployt automatisch auf Server
5. Live prüfen: https://fehlerundstolzdrauf.de
