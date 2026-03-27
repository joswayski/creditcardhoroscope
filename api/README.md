# API

For the required environment variables, please see [.env.example](.env.example)



### Routes

`GET /` Root. Maps to health check endpoint
`GET /api/v1/health` Healthcheck endpoint
`POST /api/v1/payment-intents` Creates a payment intent in stripe
`POST /api/v1/horoscopes` Creates a horoscope

Catch all --> 404