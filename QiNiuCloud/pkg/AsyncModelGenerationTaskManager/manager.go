package AsyncModelGenerationTaskManager

import (
	"QiNiuCloud/QiNiuCloud/pkg/logger"
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

type SyncModelGenerationTaskManager interface {
	AddTask(ctx context.Context, gtid string, token string) error
	QueryTask(txid string, bizData interface{}) error
}

type Action interface {
	Try(ctx context.Context, tccID string, bizData interface{}) error
	Confirm(ctx context.Context, tccID string, jobID string) error
}

type TransactionStatus string

const (
	StatusTrying     TransactionStatus = "TRYING"     // 尝试中
	StatusConfirming TransactionStatus = "CONFIRMING" // 确认中
	StatusCanceling  TransactionStatus = "CANCELING"  // 取消中
	StatusCompleted  TransactionStatus = "COMPLETED"  // 已完成
	StatusCanceled   TransactionStatus = "CANCELED"   // 已取消
)

type ModelGenerateTransactionData struct {
	KeyWordsToken string
	TransactionId string
}
type TaskManager struct {
	l               logger.ZapLogger
	mu              sync.RWMutex
	redis           *redis.Client
	transactions    map[string]*Transaction
	timeout         time.Duration
	providerManager ModelAPIsProviderManager
}
type Transaction struct {
	ID        string
	Token     string
	mu        sync.RWMutex
	Status    TransactionStatus
	CreatedAt time.Time
	UpdatedAt time.Time
	BizData   interface{}
}

func (m *TaskManager) AddTask(ctx context.Context, gtid string, token string) error {
	err := m.providerManager.AddTask(gtid, token)
	if err != nil {
		tx := &Transaction{
			ID:        gtid,
			Token:     token,
			Status:    StatusTrying,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			BizData:   nil,
		}
		return m.registerTX(tx)
	}
	return err
}

func (m *TaskManager) registerTX(transaction *Transaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, exists := m.transactions[transaction.ID]
	if exists {
		return errors.New("transaction exists")
	}
	m.transactions[transaction.ID] = transaction
	return nil
}

func (m *TaskManager) QueryTask(txid string, bizData interface{}) (TransactionStatus, error) {
	return m.providerManager.QueryTask(txid, bizData)
}

const maxTransactions = 5000

func (m *TaskManager) StartTransactionsChecker(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			m.checkTransactions()
			if len(m.transactions) >= maxTransactions {
				m.cleanTransactions()
			}
		}
	}()
}
func (m *TaskManager) cleanTransactions() {
	m.mu.RLock()
	defer m.mu.RUnlock()
	newTransactions := make(map[string]*Transaction)
	for k, v := range m.transactions {
		if v.Status == StatusTrying {
			newTransactions[k] = v
		}
	}
	m.transactions = newTransactions
}
func (m *TaskManager) updateStatus(tx *Transaction, status TransactionStatus) {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	tx.Status = status
}
func (m *TaskManager) executeCompletePhase(tx *Transaction) {
	m.updateStatus(tx, StatusCompleted)
}
func (m *TaskManager) checkTransactions() {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, tx := range m.transactions {
		if tx.Status == StatusTrying {
			status, err := m.QueryTask(tx.ID, tx.BizData)
			if err != nil {
				if status == StatusCompleted {
					m.executeCompletePhase(tx)
				}
			}
		}
	}
}
