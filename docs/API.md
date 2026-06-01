# API návrh (Fiber)

## GET /health
Health check.

## GET /api/readings
Seznam věšteb.

## GET /api/readings/:id
Detail jedné věštby.

## POST /api/readings
Vytvoření a uložení věštby.

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
