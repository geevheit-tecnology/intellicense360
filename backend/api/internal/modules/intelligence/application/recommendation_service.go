package application

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/intelligence/domain"
)

type OperationalSummary struct {
	CurrentOperation string                  `json:"current_operation"`
	CurrentRisks     []string                `json:"current_risks"`
	CurrentSavings   string                  `json:"current_savings"`
	Priorities       []string                `json:"priorities"`
	Recommendations  []domain.Recommendation `json:"recommendations"`
}

type RecommendationService struct{}

func NewRecommendationService() RecommendationService {
	return RecommendationService{}
}

func (RecommendationService) CurrentOperationalSummary(ctx context.Context) OperationalSummary {
	_ = ctx

	return OperationalSummary{
		CurrentOperation: "Operacao em modo de implantacao",
		CurrentRisks: []string{
			"Sem telemetria integrada",
			"Sem historico operacional consolidado",
		},
		CurrentSavings: "Aguardando eventos reais para estimativa",
		Priorities: []string{
			"Conectar identidade e tenants",
			"Registrar eventos operacionais",
			"Classificar custos por veiculo",
		},
		Recommendations: []domain.Recommendation{
			{
				ID:          "bootstrap-intelligence-engine",
				Title:       "Iniciar motor de inteligencia",
				Description: "Capturar eventos operacionais antes de expandir telas administrativas.",
				Priority:    domain.PriorityHigh,
				ImpactArea:  "mission-control",
			},
		},
	}
}
