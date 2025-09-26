package tcc

import (
	"context"
	"fmt"
)

type TCCManagerV1 struct {
	manager  *TCCManager
	adaptors []TccAction
}

func NewTCCMangaerV1(manager *TCCManager, adaptors []TccAction) *TCCManagerV1 {
	return &TCCManagerV1{
		manager:  manager,
		adaptors: adaptors,
	}
}

func (m *TCCManagerV1) Try(ctx context.Context, gtid string, data interface{}) error {
	for _, adaptor := range m.adaptors {
		_ = m.manager.RegisterAction(gtid, adaptor)
	}
	tx := m.manager.NewTransaction(gtid, data)
	if err := m.manager.RunTransaction(ctx, tx); err != nil {
		////处理事务执行失败
		fmt.Printf("Transaction failed: %v\n", err)
		return err
	} else {
		fmt.Println("Transaction completed successfully")
		return nil
	}
}

func (m *TCCManagerV1) Confirm(ctx context.Context, gtid string) error {
	status, err := m.manager.CheckTransactionStatus(ctx, gtid)
	if err != nil {
		return err
	}
	if status == StatusCompleted {
		// Generate Pay Order
	}
	if status == StatusCanceled {
		// Canceled
	}
	return nil
}

func (m *TCCManagerV1) Cancel(ctx context.Context, gtid string) error {
	panic("implement me")
}
