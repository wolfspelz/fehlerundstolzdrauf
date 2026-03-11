---
description: Wöchentliche Redaktionspflege für fehlerundstolzdrauf.de
user_invocable: true
---

# /redaktion – Wöchentliche Redaktionspflege

Dieses Skill dient der wöchentlichen Pflege der Website fehlerundstolzdrauf.de.

## Setup

Der Admin-Token wird aus dem GitHub Secret `ADMIN_TOKEN` gelesen. Für lokale Nutzung:
```bash
export ADMIN_TOKEN="dein-token-hier"
```

Die Basis-URL ist `https://fehlerundstolzdrauf.de` (oder `http://localhost:8080` für lokale Tests).

## Funktionen

Frage den Benutzer, welche Funktion er nutzen möchte:

### 1. Moderieren

Unmoderierte Einreichungen abrufen und einzeln freigeben oder verbergen.

```bash
# Unmoderierte anzeigen
curl -s -H "Authorization: Bearer $ADMIN_TOKEN" https://fehlerundstolzdrauf.de/admin/submissions?status=unmoderated | jq .

# Freigeben (approved)
curl -s -X PUT -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json" -d '{"status":"approved"}' https://fehlerundstolzdrauf.de/admin/submissions/ID

# Verbergen (hidden)
curl -s -X PUT -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json" -d '{"status":"hidden"}' https://fehlerundstolzdrauf.de/admin/submissions/ID
```

Zeige jede unmoderierte Einreichung dem Benutzer mit den Optionen:
1. Freigeben
2. Verbergen
3. Löschen

### 2. Zitate generieren

Generiere 3-5 neue Zitate über Fehler, Scheitern, Neubeginn. Stil: Bekannte Persönlichkeiten, Philosophen, Schriftsteller. Keine Privatpersonen.

```bash
curl -s -X POST -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json" \
  -d '{"text":"ZITAT_TEXT","attribution":"PERSON"}' \
  https://fehlerundstolzdrauf.de/admin/quotes
```

### 3. Historisch generieren

Generiere 1-2 neue historische Fehlschlag-/Zufallsentdeckungs-Stories. Stil: Hemingway – kurze Sätze, direkt, kein Pathos. 80-120 Wörter.

```bash
curl -s -X POST -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json" \
  -d '{"year":"JAHR","title":"TITEL","text":"TEXT"}' \
  https://fehlerundstolzdrauf.de/admin/historical
```

### 4. Eigene Stories schreiben

Entwickle interaktiv mit dem Benutzer neue Kurzgeschichten im Hemingway-Stil. Max 500 Zeichen. Kurze Sätze. Direkt. Kein Pathos.

```bash
curl -s -X POST -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json" \
  -d '{"year":"JAHR","title":"TITEL","text":"TEXT"}' \
  https://fehlerundstolzdrauf.de/admin/stories
```

### 5. Featured Story schreiben

Entwickle eine längere Featured Story mit Intro, Zitat-Block und Outro.

```bash
curl -s -X POST -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json" \
  -d '{"year_range":"ZEITRAUM","title":"TITEL","intro":"INTRO","quote":"ZITAT","outro":"OUTRO"}' \
  https://fehlerundstolzdrauf.de/admin/featured
```

### 6. Statistik

```bash
curl -s -H "Authorization: Bearer $ADMIN_TOKEN" https://fehlerundstolzdrauf.de/admin/stats | jq .
```

### 7. Edition zurücksetzen

Cache löschen und neuen Seed erzwingen, damit `/api/edition` sofort neue Inhalte liefert (ohne bis morgen zu warten).

```bash
curl -s -X POST -H "Authorization: Bearer $ADMIN_TOKEN" https://fehlerundstolzdrauf.de/admin/reset-edition
```

### 8. Beenden

Redaktion beenden. Keine weitere Aktion.

### 9. Löschen

```bash
# Typ: stories, featured, quotes, historical
curl -s -X DELETE -H "Authorization: Bearer $ADMIN_TOKEN" https://fehlerundstolzdrauf.de/admin/TYP/ID
```

## Ablauf

1. Statistik anzeigen, damit der Benutzer den aktuellen Stand sieht
2. Fragen, welche Funktion gewünscht ist
3. Ausführen
4. Bei generierten Inhalten: Vorschau zeigen und Bestätigung abwarten, bevor sie eingefügt werden
