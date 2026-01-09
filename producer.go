package audit

type Producer interface {
	ProduceAsync(topics []string, key, value []byte)
	EnsureTopics(topics []string) error
	Close() error
}

