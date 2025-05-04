package events

import "context"

type Producer interface {
	ProduceInconsistentEvent(ctx context.Context, evt InconsistentEvent) error
}

type InconsistentEvent struct {
	ID        int64
	Direction string
	Type      string
}

const (
	InconsistentEventTypeTargetMissing = "target_missing"
	InconsistentEventTypeBaseMissing   = "base_missing"
	InconsistentEventTypeNEQ           = "neq"
)
