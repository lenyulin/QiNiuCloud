package AsyncModelGenerationTaskManager

import (
	"QiNiuCloud/QiNiuCloud/pkg/AsyncModelGenerationTaskManager/Providers"
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type ModelAPIsProviderManager interface {
	AddTask(txid string, bizData interface{}) error
	QueryTask(txid string, bizData interface{}) (TransactionStatus, error)
}
type ModelGenerateStatus string

var (
	StatusAdd     ModelGenerateStatus = "ADD"
	StatusSucceed ModelGenerateStatus = "SUCCEED"
	StatusFailed  ModelGenerateStatus = "FAILED"
)

type generatorManager struct {
	l             logger.LoggerV1
	client        *http.Client
	mu            sync.RWMutex
	pTransactions map[string]*ProviderTransaction
	providers     []Providers.ProviderSpecificGenerator
}
type ProviderTransaction struct {
	mu        sync.RWMutex
	Status    ModelGenerateStatus
	RequestID string
}

var (
	ErrAllProviders = errors.New("all providers failed")
)

func (g *generatorManager) AddTask(txid string, bizData interface{}) error {
	var wg sync.WaitGroup
	wg.Add(len(g.providers))
	var res []*Providers.Result
	for _, provider := range g.providers {
		go func(provider Providers.ProviderSpecificGenerator) {
			re := provider.SubmitTask(g.client, bizData)
			if re != nil && re.Err != nil {
				g.mu.Lock()
				defer g.mu.Unlock()
				g.pTransactions[fmt.Sprintf("%s:%s", txid, re.Provider)] = &ProviderTransaction{
					Status:    StatusAdd,
					RequestID: re.RequestId,
				}
				res = append(res, re)
			}
		}(provider)
	}
	wg.Wait()
	if len(res) == len(g.providers) {
		return ErrAllProviders
	}
	return nil
}

func (g *generatorManager) QueryTask(txid string, bizData interface{}) (TransactionStatus, error) {
	var wg sync.WaitGroup
	wg.Add(len(g.providers))
	var res []*Providers.Result
	for _, provider := range g.providers {
		go func(provider Providers.ProviderSpecificGenerator) {
			er := provider.QueryTask(g.client, txid, provider.GetProviderName(), bizData.(string))
			if er != nil {
				res = append(res, er)
			}
		}(provider)
	}
	wg.Wait()
	//可以进一步处理
	for _, r := range res {
		if r.Err != nil {
			if r.Msg == "StatusTrying" {
				return StatusTrying, nil
			}
		}
	}
	return StatusCompleted, nil
}

func NewModelAPIsProviderManager(l logger.LoggerV1, providers []Providers.ProviderSpecificGenerator) ModelAPIsProviderManager {
	client := createHTTPClientWithConnectPool()
	defer func() {
		if transport, ok := client.Transport.(*http.Transport); ok {
			transport.CloseIdleConnections()
		}
	}()
	return &generatorManager{
		l:         l,
		client:    client,
		providers: providers,
	}
}
func createHTTPClientWithConnectPool() *http.Client {
	transport := &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 5,
		IdleConnTimeout:     30 * time.Second,
		DisableCompression:  false,
	}
	return &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}
}
