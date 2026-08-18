package domain

import (
	"strings"
	"time"
)

func SeverityRank(value Severity) int {
	switch value {
	case SeverityCritical:
		return 5
	case SeverityHigh:
		return 4
	case SeverityMedium:
		return 3
	case SeverityLow:
		return 2
	default:
		return 1
	}
}

func PriorityRank(value Priority) int {
	switch value {
	case PriorityCritical:
		return 5
	case PriorityUrgent:
		return 4
	case PriorityHigh:
		return 3
	case PriorityNormal:
		return 2
	default:
		return 1
	}
}

func CanTransition(from CommandStatus, to CommandStatus) bool {
	switch from {
	case StatusOpen:
		return to == StatusAcknowledged || to == StatusDismissed || to == StatusResolved
	case StatusAcknowledged:
		return to == StatusInProgress || to == StatusDismissed || to == StatusResolved
	case StatusInProgress:
		return to == StatusResolved || to == StatusDismissed
	default:
		return false
	}
}

func RiskLevel(score float64) string {
	switch {
	case score <= 0.20:
		return "very_low"
	case score <= 0.40:
		return "low"
	case score <= 0.60:
		return "medium"
	case score <= 0.80:
		return "high"
	default:
		return "critical"
	}
}

func HealthLevel(score float64) string {
	switch {
	case score >= 90:
		return "excellent"
	case score >= 75:
		return "good"
	case score >= 60:
		return "attention"
	case score >= 40:
		return "warning"
	default:
		return "critical"
	}
}

func SLA(dueAt *time.Time, now time.Time) SLAStatus {
	if dueAt == nil {
		return SLANotApplicable
	}
	if now.After(*dueAt) {
		return SLABreached
	}
	if dueAt.Sub(now) <= 2*time.Hour {
		return SLAAtRisk
	}
	return SLAWithin
}

func Fingerprint(tenantID string, itemType CommandItemType, category Category, sourceType string, sourceID string, title string) string {
	parts := []string{tenantID, string(itemType), string(category), sourceType, sourceID, title}
	for i, part := range parts {
		parts[i] = strings.ToLower(strings.TrimSpace(part))
	}
	return strings.Join(parts, ":")
}
