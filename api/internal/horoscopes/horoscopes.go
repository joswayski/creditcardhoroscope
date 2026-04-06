package horoscopes

import (
	"fmt"
	"time"

	"github.com/joswayski/creditcardhoroscope/api/internal/config"
	"github.com/joswayski/creditcardhoroscope/api/internal/types"
)

func GetSystemPrompt(cfg *config.Config) string {
	return fmt.Sprintf(cfg.AISystemPrompt, time.Now().UTC())
}

func FormatUserMessage(dbPaymentIntent *types.PaymentIntent) string {
	var brand, expMonth, expYear, last4, country, postalCode string

	if dbPaymentIntent.CardBrand != nil {
		brand = *dbPaymentIntent.CardBrand
	} else {
		brand = getMissingMessage("BRAND")
	}

	if dbPaymentIntent.CardExpMonth != nil {
		expMonth = *dbPaymentIntent.CardExpMonth
	} else {
		expMonth = getMissingMessage("EXPIRATION MONTH")
	}

	if dbPaymentIntent.CardExpYear != nil {
		expYear = *dbPaymentIntent.CardExpYear
	} else {
		expYear = getMissingMessage("EXPIRATION YEAR")
	}

	if dbPaymentIntent.CardLast4 != nil {
		last4 = *dbPaymentIntent.CardLast4
	} else {
		last4 = getMissingMessage("LAST 4 DIGITS")
	}

	if dbPaymentIntent.CardCountry != nil {
		country = *dbPaymentIntent.CardCountry
	} else {
		country = getMissingMessage("COUNTRY")
	}

	if dbPaymentIntent.CardPostal != nil {
		postalCode = *dbPaymentIntent.CardPostal
	} else {
		postalCode = getMissingMessage("POSTAL CODE / ZIP CODE")
	}

	return fmt.Sprintf(`
	Create a fun horoscope!

Card Details:
- Card Brand: {%v}
- Expiration Month: {%v}
- Expiration Year: {%v}
- Last 4 digits: {%v}
- Country: {%v}
- Postal Code / Zip: {%v}
`, brand, expMonth, expYear, last4, country, postalCode)
}

func getMissingMessage(field string) string {
	return fmt.Sprintf("{{ADMIN MESSAGE: NO CARD %s PROVIDED - DO NOT INCLUDE THIS IN THE RESPONSE}}", field)
}
