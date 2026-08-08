package domain

type RecommendationPriority string

const (
	PriorityLow      RecommendationPriority = "low"
	PriorityMedium   RecommendationPriority = "medium"
	PriorityHigh     RecommendationPriority = "high"
	PriorityCritical RecommendationPriority = "critical"
)

type Recommendation struct {
	ID          string                 `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Priority    RecommendationPriority `json:"priority"`
	ImpactArea  string                 `json:"impact_area"`
}
