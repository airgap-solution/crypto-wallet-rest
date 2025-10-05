package provider

type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}
func (a *Adapter) GetBalance(symbol, address string) (float64, error) {
	return 0, nil
}
