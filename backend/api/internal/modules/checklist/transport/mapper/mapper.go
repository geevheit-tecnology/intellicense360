package mapper

import (
	"github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/transport/dto"
)

func ChecklistFromRequest(req dto.ChecklistRequest) domain.Checklist {
	return domain.Checklist{
		VehicleID:      req.VehicleID,
		Title:          req.Title,
		Description:    req.Description,
		Type:           req.Type,
		Status:         domain.ChecklistStatus(req.Status),
		DriverName:     req.DriverName,
		DriverDocument: req.DriverDocument,
	}
}

func ItemFromRequest(req dto.ChecklistItemRequest) domain.ChecklistItem {
	return domain.ChecklistItem{
		Title:         req.Title,
		Description:   req.Description,
		Category:      req.Category,
		Required:      req.Required,
		OrderIndex:    req.OrderIndex,
		AnswerType:    domain.AnswerType(req.AnswerType),
		ExpectedValue: req.ExpectedValue,
	}
}

func AnswerFromRequest(req dto.ChecklistAnswerRequest) domain.ChecklistAnswer {
	return domain.ChecklistAnswer{
		ChecklistItemID: req.ChecklistItemID,
		Answer:          req.Answer,
		Notes:           req.Notes,
		PhotoURL:        req.PhotoURL,
	}
}
