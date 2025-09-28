package Providers

import (
	"net/http"
)

type ProviderSpecificGenerator interface {
	SubmitTask(client *http.Client, bizData interface{}) *Result
	QueryTask(client *http.Client, txId string, jobId string, token string) *Result
	GetProviderName() string
}

type ModelGenerationTaskResult struct {
	jobiD string
	token string
	url   string
	thumb string
}
