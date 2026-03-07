package util

const (
	USD = "USD"
	EUR = "EUR"
	RUB = "RUB"
	EGP = "EGP"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, RUB, EGP:
		return true
	}
	return false
}
