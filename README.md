[![Go Report Card](https://goreportcard.com/badge/github.com/Luqqk/goexrates)](https://goreportcard.com/report/github.com/Luqqk/goexrates)

## 💰 goexrates

A foreign exchange rates and currency conversion API that exposes data published by European Central Bank.

The rates are updated daily around 4PM CET.

### **Usage**

Get the latest foreign exchange reference rates in JSON format.

```http
GET /v1/latest
Host: localhost:3000
```

Get historical rates for any day since 1999-01-04.

```http
GET /v1/historical/2008-03-18
Host: localhost:3000
```

Rates are quoted against the Euro by default. Quote against a different currency by setting the base parameter in your request.

```http
GET /v1/latest?base=USD
Host: localhost:3000
```

Request specific exchange rates by setting the codes parameter.

```http
GET /v1/latest?codes=USD,GBP
Host: localhost:3000
```

Response format.

```json
{
    "base": "EUR",
    "date": "2021-05-05",
    "rates": {
        "AUD": 1.4832,
        "PLN": 4.2173,
        "MYR": 4.7543,
        "USD": 1.0961,
        [41 world currencies],
    }
}
```

### **Run**

```bash
# Build and start containers:
docker compose up -d
# Enter api container:
docker exec -it api bash
# Populate database with rates published by ECB (check goexrates-cli --help):
goexrates-cli load historical
# Run the API
goexrates-api
```

API's endpoints can be accessed at `localhost:3000` (container's port is published).
