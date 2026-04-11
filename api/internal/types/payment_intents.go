package types

import "time"

type PaymentIntent struct {
	ID              int64     `db:"id"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
	PaymentIntentID string    `db:"payment_intent_id"`
	Amount          int       `db:"amount"`
	Currency        string    `db:"currency"`
	Status          string    `db:"status"`
	CardBrand       *string   `db:"card_brand"`
	CardExpMonth    *string   `db:"card_exp_month"`
	CardExpYear     *string   `db:"card_exp_year"`
	CardLast4       *string   `db:"card_last_4"`
	CardCountry     *string   `db:"card_country"`
	CardPostal      *string   `db:"card_postal"`
}

const maxGenerationCount = 3

// Pending gets let through because we're awaiting a generation
// Paid gets let through because we'll allow multiple generations (TODO)
func (pi *PaymentIntent) AllowsGenerations(currentGenerationCount int) bool {
	return (pi.Status == "pending" || pi.Status == "paid") && currentGenerationCount < maxGenerationCount
}
