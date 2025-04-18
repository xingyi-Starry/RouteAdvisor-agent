package services

type Controller struct {
	DiaryWriter *DiaryWriterService
}

func NewController(diaryWriter *DiaryWriterService) *Controller {
	return &Controller{
		DiaryWriter: diaryWriter,
	}
}
