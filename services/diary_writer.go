package services

import (
	"route-advisor-agent/agents/diary_writer"
	"route-advisor-agent/models"
	"route-advisor-agent/utils/quick_response"

	"github.com/gofiber/fiber/v2"
)

type DiaryWriterService struct {
	DiaryWriter *diary_writer.DiaryWriter
}

func NewDiaryWriterService(diaryWriter *diary_writer.DiaryWriter) *DiaryWriterService {
	return &DiaryWriterService{
		DiaryWriter: diaryWriter,
	}
}

func (s *DiaryWriterService) DiaryWriterHandler(c *fiber.Ctx) error {
	// 解析请求
	var diaryConfig models.DiaryConfig
	if err := c.BodyParser(&diaryConfig); err != nil {
		return quick_response.BadRequest(c, "Invalid request body")
	}

	// 生成游记
	diary, err := s.DiaryWriter.CreateDiary(c.Context(), diaryConfig)
	if err != nil {
		return quick_response.Internal(c, "Failed to generate diary")
	}

	// 返回结果
	result := models.DiaryResult{
		Code:    200,
		Message: "OK",
		Data:    diary,
	}
	return c.Status(200).JSON(result)
}
