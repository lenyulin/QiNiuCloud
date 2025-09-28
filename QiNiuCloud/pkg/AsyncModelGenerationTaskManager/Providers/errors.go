package Providers

type ProviderResult []*Result

func (e *ProviderResult) Count() int {
	return len(*e)
}

type Result struct {
	Provider  string
	RequestId string
	JobId     string
	Msg       string
	Err       error
}
