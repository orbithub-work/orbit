package services

import (
	"context"
	"sync"
	"time"
)

const (
	scanPolicyCheckInterval    = 30 * time.Second
	scanPolicyQueueThreshold   = 8
	scanPolicyPendingThreshold = 3000
)

type scanSystemPolicy struct {
	mu               sync.RWMutex
	checkInterval    time.Duration
	queueThreshold   int
	pendingThreshold int
	heavyTasksMuted  bool
	heavyTasksReason string
}

func newScanSystemPolicy() *scanSystemPolicy {
	return &scanSystemPolicy{
		checkInterval:    scanPolicyCheckInterval,
		queueThreshold:   scanPolicyQueueThreshold,
		pendingThreshold: scanPolicyPendingThreshold,
	}
}

func (p *scanSystemPolicy) interval() time.Duration {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.checkInterval <= 0 {
		return scanPolicyCheckInterval
	}
	return p.checkInterval
}

func (p *scanSystemPolicy) evaluate(running bool, queueSize int, pendingCount int) (bool, string) {
	p.mu.RLock()
	queueThreshold := p.queueThreshold
	pendingThreshold := p.pendingThreshold
	p.mu.RUnlock()

	shouldMute := running || queueSize >= queueThreshold || pendingCount >= pendingThreshold
	reason := ""
	if running {
		reason = "scan_running"
	} else if queueSize >= queueThreshold {
		reason = "scan_queue_busy"
	} else if pendingCount >= pendingThreshold {
		reason = "asset_pipeline_busy"
	}
	return shouldMute, reason
}

func (p *scanSystemPolicy) apply(shouldMute bool, reason string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	changed := p.heavyTasksMuted != shouldMute || p.heavyTasksReason != reason
	p.heavyTasksMuted = shouldMute
	p.heavyTasksReason = reason
	return changed
}

func (p *scanSystemPolicy) shouldRunHeavyTasks() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return !p.heavyTasksMuted
}

func (s *ScanService) runPerformancePolicy() {
	interval := scanPolicyCheckInterval
	if s.policy != nil {
		interval = s.policy.interval()
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			pendingCount, _ := s.assetService.assets.GetPendingCount(context.Background())
			queueSize := len(s.scanQueue)

			s.mu.Lock()
			running := s.running
			s.mu.Unlock()

			shouldMute := false
			reason := ""
			if s.policy != nil {
				shouldMute, reason = s.policy.evaluate(running, queueSize, pendingCount)
			}
			changed := false
			if s.policy != nil {
				changed = s.policy.apply(shouldMute, reason)
			}

			if changed && s.eventHub != nil {
				mode := "heavy"
				if shouldMute {
					mode = "silent"
				}
				s.eventHub.Broadcast(map[string]any{
					"type": "scan_mode_changed",
					"data": map[string]any{
						"mode":          mode,
						"reason":        reason,
						"queue_size":    queueSize,
						"pending_count": pendingCount,
					},
				})
			}
		case <-s.stopChan:
			return
		}
	}
}

func (s *ScanService) shouldRunHeavyTasks() bool {
	if s.policy == nil {
		return true
	}
	return s.policy.shouldRunHeavyTasks()
}
