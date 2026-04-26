package tasks

type CreateTaskInput struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Status      string `json:"status" validate:"omitempty,oneof=pending in-progress done"`
	Deadline    string `json:"deadline" validate:"omitempty,datetime=2006-01-02"`
}

type UpdateTaskInput struct {
	Title       *string `json:"title" validate:"omitempty,min=1"`
	Description *string `json:"description" validate:"omitempty"`
	Status      *string `json:"status" validate:"omitempty,oneof=pending in-progress done"`
	Deadline    *string `json:"deadline" validate:"omitempty,datetime=2006-01-02"`
}

func (input UpdateTaskInput) HasUpdateFields() bool {
	return input.Title != nil ||
		input.Description != nil ||
		input.Status != nil ||
		input.Deadline != nil
}

type ListTaskQuery struct {
	Status string `form:"status" validate:"omitempty,oneof=pending in-progress done"`
	Page   int    `form:"page" validate:"omitempty,min=1"`
	Limit  int    `form:"limit" validate:"omitempty,min=1,max=100"`
}
