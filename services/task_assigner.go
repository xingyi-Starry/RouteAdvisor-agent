package services

import (
	"route-advisor-agent/agents/task_assigner"
	"route-advisor-agent/models"
	"route-advisor-agent/utils/quick_response"

	"github.com/gofiber/fiber/v2"
)

type TaskAssignerService struct {
	TaskAssigner *task_assigner.TaskAssigner
}

func NewTaskAssignerService(taskAssigner *task_assigner.TaskAssigner) *TaskAssignerService {
	return &TaskAssignerService{
		TaskAssigner: taskAssigner,
	}
}

func (s *TaskAssignerService) AssignTaskHandler(c *fiber.Ctx) error {
	// 解析请求
	var taskConfig models.TaskAssignConfig
	if err := c.BodyParser(&taskConfig); err != nil {
		return quick_response.BadRequest(c, "Invalid request body")
	}

	// 分配任务
	result, err := s.TaskAssigner.AssignTask(c.Context(), taskConfig)
	if err != nil {
		return quick_response.Internal(c, "Failed to assign task")
	}

	// 返回结果
	vo := models.TaskAssignVO{
		Code:    200,
		Message: "OK",
		Data:    result,
	}
	return c.JSON(vo)
}
