package services

type Controller struct {
	DiaryWriter  *DiaryWriterService
	TaskAssigner *TaskAssignerService
}

func NewController(diaryWriter *DiaryWriterService, taskAssigner *TaskAssignerService) *Controller {
	return &Controller{
		DiaryWriter:  diaryWriter,
		TaskAssigner: taskAssigner,
	}
}
