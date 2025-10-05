package provider

type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) GetBalance(_, _ string) (float64, error) {
	return 0, nil
}
