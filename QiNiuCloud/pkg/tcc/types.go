package tcc

type TransactionStatus string

const (
	StatusTrying     TransactionStatus = "TRYING"     // 尝试中
	StatusConfirming TransactionStatus = "CONFIRMING" // 确认中
	StatusCanceling  TransactionStatus = "CANCELING"  // 取消中
	StatusCompleted  TransactionStatus = "COMPLETED"  // 已完成
	StatusCanceled   TransactionStatus = "CANCELED"   // 已取消
)
