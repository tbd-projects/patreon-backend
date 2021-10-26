package models

type AccessCounter struct {
	Counter int64
	Limit   int64
}

func NewAccessCounter(limit int64) *AccessCounter {
	return &AccessCounter{
		Counter: 0,
		Limit:   limit,
	}
}
func (ac *AccessCounter) Overflow() bool {
	return ac.Counter >= ac.Limit
}
