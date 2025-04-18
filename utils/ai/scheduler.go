package ai

const (
	STRATEGY_ROUND_ROBIN = iota
)

type Strategy interface {
	Next() int
}

type Scheduler struct {
	Models   []*ModelLLM
	total    int
	Strategy Strategy
}

func NewScheduler(strategy int, models ...*ModelLLM) *Scheduler {
	balancer := &Scheduler{
		Models: models,
		total:  len(models),
	}
	switch strategy {
	case STRATEGY_ROUND_ROBIN:
		balancer.Strategy = NewRoundRobin(&balancer.total)
	default:
		balancer.Strategy = NewRoundRobin(&balancer.total)
	}
	return balancer
}

func (b *Scheduler) Run() {
	for _, model := range b.Models {
		model.Run()
	}
}

func (b *Scheduler) AddModel(model *ModelLLM) {
	b.Models = append(b.Models, model)
	b.total++
}

func (b *Scheduler) Enqueue(t task) {
	b.Models[b.Strategy.Next()].Enqueue(t)
}

// 轮询
type RoundRobin struct {
	total *int
	index int
}

func NewRoundRobin(total *int) *RoundRobin {
	return &RoundRobin{
		total: total,
		index: 0,
	}
}

func (r *RoundRobin) Next() int {
	index := r.index % *r.total
	r.index++
	return index
}
