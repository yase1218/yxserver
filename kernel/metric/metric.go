package metric

import "sync/atomic"

type ProcessorMetrics struct {
	PushAttempts  atomic.Int64
	PushSuccess   atomic.Int64
	PushRejects   atomic.Int64
	ProcessCount  atomic.Int64
	ProcessErrors atomic.Int64
}

func (m *ProcessorMetrics) Info() map[string]interface{} {
	attemps := m.PushAttempts.Load()
	if attemps > 0 {
		return map[string]interface{}{
			"push_attempts":  m.PushAttempts.Load(),
			"push_success":   m.PushSuccess.Load(),
			"push_rejects":   m.PushRejects.Load(),
			"process_count":  m.ProcessCount.Load(),
			"process_errors": m.ProcessErrors.Load(),
			"reject_rate":    float64(m.PushRejects.Load()) / float64(m.PushAttempts.Load()),
			"success_rate":   float64(m.PushSuccess.Load()) / float64(m.PushAttempts.Load()),
		}
	} else {
		return map[string]interface{}{
			"push_attempts":  0,
			"push_success":   0,
			"push_rejects":   0,
			"process_count":  0,
			"process_errors": 0,
			"reject_rate":    0,
			"success_rate":   0,
		}
	}
}
