# API

For the required environment variables, please see [.env.example](.env.example)


### Migrations
They are ran at API start. See [migrations.go](/api/internal/database/migrations.go) and [/migrations](/api/internal/database/migrations)

### Routes

`GET /` - Root. Maps to health check endpoint

`GET /api/v1/health` - Healthcheck endpoint

`POST /api/v1/payment-intents` - Creates a payment intent in Stripe. If we fail to save it to the DB, we attempt to cancel it.

`POST /api/v1/horoscopes` - Creates a horoscope

Catch all --> 404
