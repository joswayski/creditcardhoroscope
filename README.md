# Credit Card Horoscope

Public website for ✨ https://creditcardhoroscope.com ✨

What does your credit card say about you?

---

### Running
```bash
cp api/.env.example api/.env
cp web/.env.example web/.env

docker compose up -d
```

| Service  | URL                    | Description |
|----------|------------------------| ---------- |
| web      | http://localhost:3000   | React Router |
| api      | http://localhost:8080   | Go API |
| postgres | localhost:5432          | Postgres |


### Database
There's a `payment_intents` table and a `generations` table. That's pretty much it :)

### Misc

Stripe for payments and OpenRouter for the LLM