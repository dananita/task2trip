package handlers

import (
	"github.com/go-openapi/strfmt"
	"github.com/itimofeev/task2trip/backend"
	"github.com/itimofeev/task2trip/rest/models"
)

func convertUser(user *backend.User) *models.User {
	return &models.User{
		ID:   &user.ID,
		Name: &user.Email,
	}
}

func convertTasksPage(tasks_ []*backend.Task, total int64) *models.TaskPage {
	payload := make([]*models.Task, 0, len(tasks_))
	for _, task := range tasks_ {
		payload = append(payload, convertTask(task))
	}
	return &models.TaskPage{
		Total:   total,
		Payload: payload,
	}
}

func convertCategories(categs []*backend.Category) (res []*models.Category) {
	for _, category := range categs {
		res = append(res, &models.Category{
			ID:           &category.ID,
			Key:          &category.Key,
			DefaultValue: category.DefaultValue,
		})
	}
	return res
}

func convertTask(task *backend.Task) *models.Task {
	return &models.Task{
		ID:             &task.ID,
		Name:           &task.Name,
		Category:       convertCategory(task.Category),
		CreateTime:     strfmt.DateTime(task.CreateTime),
		BudgetEstimate: &task.BudgetEstimate,
		Description:    &task.Description,
	}
}

func convertCategory(category *backend.Category) *models.Category {
	return &models.Category{
		ID:           &category.ID,
		Key:          &category.Key,
		DefaultValue: category.DefaultValue,
	}
}

func convertOffer(offer *backend.Offer) *models.Offer {
	return &models.Offer{
		ID:      &offer.ID,
		Price:   &offer.Price,
		Comment: offer.Comment,
		User:    convertUser(offer.User),
	}
}

func convertOffers(offers_ []*backend.Offer) (res []*models.Offer) {
	for _, offer := range offers_ {
		res = append(res, convertOffer(offer))
	}
	return res
}
