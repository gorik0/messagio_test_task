package service

type Service struct {
	St Storager
}

type Storager interface {
	CloseDB() error
	AddRecord(string, string) (int64, error)
	MessageConsumed(int64) error                  // need to work on
	NumberOfProcessedMessages() (int, int, error) // need to work on
	Produce(string, []byte) error
	Consume() (string, []byte, error)
}
