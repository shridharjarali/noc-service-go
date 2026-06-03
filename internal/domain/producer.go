package domain

// Producer abstracts Kafka message publishing.
// Implementation is NOT in scope for Phase 3 — only the interface.
type Producer interface {
	Push(topic string, value interface{}) error
}
