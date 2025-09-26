package tcc

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Transaction struct {
	ID        string
	Status    TransactionStatus
	Actions   []TccAction
	BizData   interface{}
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TccAction interface {
	Try(ctx context.Context, tccID string, bizData interface{}) error
	Confirm(ctx context.Context, tccID string) error
	Cancel(ctx context.Context, tccID string) error
}

type TCCManager struct {
	transactions map[string]*Transaction
	mu           sync.RWMutex
	timeout      time.Duration
}

func NewTCCManager(timeout time.Duration) *TCCManager {
	return &TCCManager{
		transactions: make(map[string]*Transaction),
		timeout:      timeout,
	}
}

func (m *TCCManager) NewTransaction(gtid string, bizData interface{}) *Transaction {
	return &Transaction{
		ID:        gtid,
		Status:    StatusTrying,
		Actions:   make([]TccAction, 0),
		BizData:   bizData,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (m *TCCManager) RegisterAction(gtid string, action TccAction) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	tx, exists := m.transactions[gtid]
	if exists {
		return errors.New("transaction exists")
	}

	tx.Actions = append(tx.Actions, action)
	tx.UpdatedAt = time.Now()
	return nil
}

func (m *TCCManager) RunTransaction(ctx context.Context, tx *Transaction) error {
	// 保存事务到管理器
	m.mu.Lock()
	m.transactions[tx.ID] = tx
	m.mu.Unlock()

	if err := m.executeTryPhase(ctx, tx); err != nil {
		if cancelErr := m.executeCancelPhase(ctx, tx); cancelErr != nil {
			return fmt.Errorf("try failed and cancel error: %v, cancel err: %v", err, cancelErr)
		}
		return fmt.Errorf("try phase failed: %v", err)
	}

	if err := m.executeConfirmPhase(ctx, tx); err != nil {
		return fmt.Errorf("confirm phase failed: %v", err)
	}

	return nil
}

func (m *TCCManager) executeTryPhase(ctx context.Context, tx *Transaction) error {
	m.updateStatus(tx.ID, StatusTrying)
	for i, action := range tx.Actions {
		if err := action.Try(ctx, tx.ID, tx.BizData); err != nil {
			return fmt.Errorf("action %d try failed: %v", i, err)
		}
	}
	return nil
}

func (m *TCCManager) executeConfirmPhase(ctx context.Context, tx *Transaction) error {
	m.updateStatus(tx.ID, StatusConfirming)
	for _, action := range tx.Actions {
		if err := action.Confirm(ctx, tx.ID); err != nil {
			return fmt.Errorf("confirm failed: %v", err)
		}
	}
	m.updateStatus(tx.ID, StatusCompleted)
	return nil
}

func (m *TCCManager) executeCancelPhase(ctx context.Context, tx *Transaction) error {
	m.updateStatus(tx.ID, StatusCanceling)
	for i := len(tx.Actions) - 1; i >= 0; i-- {
		action := tx.Actions[i]
		if err := action.Cancel(ctx, tx.ID); err != nil {
			fmt.Printf("warning: action %d cancel failed: %v\n", i, err)
		}
	}
	m.updateStatus(tx.ID, StatusCanceled)
	return nil
}

func (m *TCCManager) StartTimeoutChecker(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			m.checkTimeoutTransactions()
		}
	}()
}

func (m *TCCManager) checkTimeoutTransactions() {
	m.mu.RLock()
	defer m.mu.RUnlock()
	now := time.Now()
	for _, tx := range m.transactions {
		if tx.Status == StatusTrying && now.Sub(tx.UpdatedAt) > m.timeout {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			go func(tx *Transaction) {
				defer cancel()
				_ = m.executeCancelPhase(ctx, tx)
			}(tx)
		}
	}
}

func (m *TCCManager) updateStatus(gtid string, status TransactionStatus) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if tx, exists := m.transactions[gtid]; exists {
		tx.Status = status
		tx.UpdatedAt = time.Now()
	}
}

func (m *TCCManager) CheckTransactionStatus(ctx context.Context, gtid string) (TransactionStatus, error) {
	tx, exists := m.transactions[gtid]
	if !exists {
		return "", errors.New("transaction does not exist")
	}
	return tx.Status, nil
}
