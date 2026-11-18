package task

type TypeStatus string

const (
	TypePending   TypeStatus = "pending"
	TypeRunning   TypeStatus = "running"
	TypeSucceeded TypeStatus = "succeeded"
	TypeFailed    TypeStatus = "failed"
	TypeCanceled  TypeStatus = "canceled"
)

type Task struct {
	Id string `json:"id"`
}
