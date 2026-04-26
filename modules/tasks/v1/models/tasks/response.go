package tasks

type TaskResponse struct {
	UUID        string  `json:"uuid"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Status      string  `json:"status"`
	Deadline    *string `json:"deadline,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

func NewTaskResponse(task *Task) TaskResponse {
	var deadline *string
	if task.Deadline != nil {
		formattedDeadline := task.Deadline.Format("2006-01-02")
		deadline = &formattedDeadline
	}

	return TaskResponse{
		UUID:        task.UUID.String(),
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Deadline:    deadline,
		CreatedAt:   task.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   task.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func NewTaskResponses(taskList []Task) []TaskResponse {
	responses := make([]TaskResponse, 0, len(taskList))
	for i := range taskList {
		responses = append(responses, NewTaskResponse(&taskList[i]))
	}

	return responses
}
