---
description: Wöchentliche Redaktionspflege für fehlerundstolzdrauf.de
user_invocable: true
---

# /redaktion – Wöchentliche Redaktionspflege

Dieses Skill dient der wöchentlichen Pflege der Website fehlerundstolzdrauf.de.

## Setup

Der Admin-Token wird aus der `.env`-Datei geladen:
```bash
source .env
```

Die Basis-URL ist `https://fehlerundstolzdrauf.de` (oder `http://localhost:8080` für lokale Tests).

## WICHTIG: UTF-8 Encoding

**Alle curl-Befehle mit Body MÜSSEN `charset=utf-8` im Content-Type haben**, sonst werden Umlaute (ä, ö, ü, ß) zerstört:
```
-H "Content-Type: application/json; charset=utf-8"
```

## Funktionen

Frage den Benutzer, welche Funktion er nutzen möchte:

### 1. Moderieren

Unmoderierte Einreichungen abrufen und einzeln freigeben oder verbergen.

```bash
# Unmoderierte anzeigen
curl -s -H "Authorization: Bearer $ADMIN_TOKEN" https://fehlerundstolzdrauf.de/admin/submissions?status=unmoderated | jq .

# Freigeben (approved)
curl -s -X PUT -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json; charset=utf-8" -d '{"status":"approved"}' https://fehlerundstolzdrauf.de/admin/submissions/ID

# Verbergen (hidden)
curl -s -X PUT -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json; charset=utf-8" -d '{"status":"hidden"}' https://fehlerundstolzdrauf.de/admin/submissions/ID
```

Zeige jede unmoderierte Einreichung dem Benutzer mit den Optionen:
1. Freigeben
2. Verbergen
3. Löschen
4. Überspringen (weiter zur nächsten Einreichung ohne Aktion)
5. Abbrechen (zurück zum Hauptmenü, restliche Einreichungen überspringen)

### 2. Zitate generieren

Generiere 3-5 neue Zitate über Fehler, Scheitern, Neubeginn. Stil: Bekannte Persönlichkeiten, Philosophen, Schriftsteller. Keine Privatpersonen.

```bash
curl -s -X POST -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json; charset=utf-8" \
  -d '{"text":"ZITAT_TEXT","attribution":"PERSON"}' \
  https://fehlerundstolzdrauf.de/admin/quotes
```

### 3. Historisch generieren

Generiere 1-2 neue historische Fehlschlag-/Zufallsentdeckungs-Stories. Stil: Hemingway – kurze Sätze, direkt, kein Pathos. 80-120 Wörter.

```bash
curl -s -X POST -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json; charset=utf-8" \
  -d '{"year":"JAHR","title":"TITEL","text":"TEXT"}' \
  https://fehlerundstolzdrauf.de/admin/historical
```

### 4. Eigene Stories schreiben

Entwickle interaktiv mit dem Benutzer neue Kurzgeschichten im Hemingway-Stil. Max 500 Zeichen. Kurze Sätze. Direkt. Kein Pathos.

**Titel-Variation**: Überschriften MÜSSEN im Muster variieren. NICHT alle nach dem Schema „Artikel + Substantiv" (z.B. „Die Bäckerei", „Der Laden", „Das Restaurant"). Stattdessen mischen zwischen:
- Substantiv ohne Artikel: „Barcelona", „Geduld", „Marathon"
- Adjektiv + Substantiv: „Vier Jahre Jura", „Achtzehn Absagen"
- Zeitangaben: „Zu früh", „Drei Semester Physik"
- Aktionen/Verben: „Blackout am Mikrofon", „Niemand brauchte sie"
- Einzelwörter: „Podcast", „Chinesisch"

Vor dem Generieren die bestehenden Titel prüfen und ein Muster wählen, das unterrepräsentiert ist.

```bash
curl -s -X POST -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json; charset=utf-8" \
  -d '{"year":"JAHR","title":"TITEL","text":"TEXT"}' \
  https://fehlerundstolzdrauf.de/admin/stories
```

### 5. Statistik

```bash
curl -s -H "Authorization: Bearer $ADMIN_TOKEN" https://fehlerundstolzdrauf.de/admin/stats | jq .
```

### 6. Neue Ausgabe

Neue Ausgabe erzwingen: Cache löschen und neuen Seed generieren, damit `/api/edition` sofort neue Inhalte liefert (ohne bis morgen zu warten).

```bash
curl -s -X POST -H "Authorization: Bearer $ADMIN_TOKEN" https://fehlerundstolzdrauf.de/admin/new-edition
```

### 7. Backup

Manuelles Backup aller Daten als SQL-Dump. Zeigt Statistik und Dateigröße.

```bash
curl -s -X POST -H "Authorization: Bearer $ADMIN_TOKEN" https://fehlerundstolzdrauf.de/admin/backup | jq .
```

### 8. Bearbeiten

Bestehende Einträge bearbeiten. Nur übergebene Felder werden geändert.

```bash
# Story bearbeiten (year, title, text)
curl -s -X PUT -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json; charset=utf-8" \
  -d '{"title":"NEUER_TITEL"}' \
  https://fehlerundstolzdrauf.de/admin/stories/ID

# Zitat bearbeiten (text, attribution)
curl -s -X PUT -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json; charset=utf-8" \
  -d '{"text":"NEUER_TEXT","attribution":"NEUE_PERSON"}' \
  https://fehlerundstolzdrauf.de/admin/quotes/ID

# Historisch bearbeiten (year, title, text)
curl -s -X PUT -H "Authorization: Bearer $ADMIN_TOKEN" -H "Content-Type: application/json; charset=utf-8" \
  -d '{"title":"NEUER_TITEL"}' \
  https://fehlerundstolzdrauf.de/admin/historical/ID
```

### 9. Löschen

```bash
# Typ: stories, quotes, historical
curl -s -X DELETE -H "Authorization: Bearer $ADMIN_TOKEN" https://fehlerundstolzdrauf.de/admin/TYP/ID
```

### 10. Beenden

Redaktion beenden. Keine weitere Aktion.

## Ablauf

1. Statistik anzeigen, damit der Benutzer den aktuellen Stand sieht
2. Fragen, welche Funktion gewünscht ist
3. Ausführen
4. Bei generierten Inhalten: Vorschau zeigen und Bestätigung abwarten, bevor sie eingefügt werden

## Duplikatvermeidung

Vor dem Generieren neuer Inhalte (Stories, Historisch, Zitate) IMMER zuerst die bestehenden Einträge von der Live-API abrufen:

```bash
# Alle Stories abrufen
curl -s -H "Authorization: Bearer $ADMIN_TOKEN" https://fehlerundstolzdrauf.de/admin/stories

# Alle Zitate abrufen
curl -s -H "Authorization: Bearer $ADMIN_TOKEN" https://fehlerundstolzdrauf.de/admin/quotes

# Alle historischen Einträge abrufen
curl -s -H "Authorization: Bearer $ADMIN_TOKEN" https://fehlerundstolzdrauf.de/admin/historical
```

Anhand der abgerufenen Daten sicherstellen: keine ähnlichen Themen, Personen, Titel oder Inhalte wie bereits vorhandene Einträge. NICHT auf seed.sql zugreifen – nur die echten Daten von der API verwenden.
