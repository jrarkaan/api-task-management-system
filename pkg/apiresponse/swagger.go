package apiresponse

import (
	"api-task-management-system/modules/accounts/v1/models/users"
	"api-task-management-system/modules/tasks/v1/models/tasks"
	"api-task-management-system/pkg/pagination"
)

// SwaggerErrorResponse represents a standard error response
type SwaggerErrorResponse struct {
	Message string `json:"message"`
	Status  uint16 `json:"status"`
	Error   XError `json:"error"`
}

// SwaggerSuccessResponse represents a basic success response
type SwaggerSuccessResponse struct {
	Meta    interface{} `json:"meta"`
	Message string      `json:"message"`
	Status  uint16      `json:"status"`
	Data    interface{} `json:"data"`
}

// SwaggerAuthLoginResponse represents response returned upon logging in
type SwaggerAuthLoginResponse struct {
	Meta    interface{}         `json:"meta"`
	Message string              `json:"message"`
	Status  uint16              `json:"status"`
	Data    users.LoginResponse `json:"data"`
}

// SwaggerTaskListResponse represents response mapping a list of tasks
type SwaggerTaskListResponse struct {
	Meta    pagination.Pagination `json:"meta"`
	Message string                `json:"message"`
	Status  uint16                `json:"status"`
	Data    []tasks.TaskResponse  `json:"data"`
}

// SwaggerTaskResponse represents response returned with single task data
type SwaggerTaskResponse struct {
	Meta    interface{}        `json:"meta"`
	Message string             `json:"message"`
	Status  uint16             `json:"status"`
	Data    tasks.TaskResponse `json:"data"`
}
