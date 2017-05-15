## 💰 goexrates

A foreign exchange rates and currency conversion API. Golang implementation of [fixer.io](http://fixer.io) (Ruby). Data from European Central Bank API.

The rates are updated daily around 3PM CET.

### **Usage**

Get the latest foreign exchange reference rates in JSON format.

```http
GET /latest
Host: localhost:3000
```

Get historical rates for any day since 1999.

```http
GET /2008-03-18
Host: localhost:3000
```

Rates are quoted against the Euro by default. Quote against a different currency by setting the base parameter in your request.

```http
GET /latest?base=USD
Host: localhost:3000
```

Request specific exchange rates by setting the symbols or currencies parameter.

```http
GET /latest?symbols=USD,GBP
Host: localhost:3000
```

Response format.

```json
{
    "base": "EUR",
    "date": "2017-05-05",
    "rates": {
        "AUD": 1.4832,
        "PLN": 4.2173,
        "MYR": 4.7543,
        "USD": 1.0961,
        "...": "and so on...",
    }
}
```

### **Run**

```bash
go run goexrates.go
```

### TODO

* CORS requests support
* Script for daily database update (newest data)

### **Important note**

Currently this API isn't available at any domain (it's just proof of concept). I created it for **.go** learning purposes. I will try to publish it later with up-to-date database.
