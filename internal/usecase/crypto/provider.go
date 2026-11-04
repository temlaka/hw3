package crypto

type PriceProvider interface {
	GetPrice(symbol string) (name string, price float64, err error)
}
