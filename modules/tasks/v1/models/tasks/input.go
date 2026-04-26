package tasks

type CreateTaskInput struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Status      string `json:"status" validate:"omitempty,oneof=pending in-progress done"`
	Deadline    string `json:"deadline" validate:"omitempty,datetime=2006-01-02"`
}

type UpdateTaskInput struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Status      string `json:"status" validate:"omitempty,oneof=pending in-progress done"`
	Deadline    string `json:"deadline" validate:"omitempty,datetime=2006-01-02"`
}

type ListTaskQuery struct {
	Status string `form:"status" validate:"omitempty,oneof=pending in-progress done"`
}
