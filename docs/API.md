# API návrh

## GET /health
Vrací stav služby.

## GET /api/readings
Seznam uložených věšteb.

## GET /api/readings/{id}
Detail věštby.

## POST /api/readings
Vytvoří a uloží věštbu.

### Request
```json
{
  "question": "Jak se mám rozhodnout?",
  "method": "manual",
  "lines": [
    {"position":1,"value":7},
    {"position":2,"value":8},
    {"position":3,"value":9},
    {"position":4,"value":8},
    {"position":5,"value":7},
    {"position":6,"value":8}
  ]
}
```
