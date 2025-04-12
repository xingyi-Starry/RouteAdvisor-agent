package ai

const (
	STRATEGY_ROUND_ROBIN = iota
)

type Strategy interface {
	Next() int
}

type LoadBalancer struct {
	Models   []*ModelLLM
	total    int
	Strategy Strategy
}

func NewLoadBalancer(strategy int, models ...*ModelLLM) *LoadBalancer {
	balancer := &LoadBalancer{
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

func (b *LoadBalancer) Run() {
	for _, model := range b.Models {
		model.Run()
	}
}

func (b *LoadBalancer) AddModel(model *ModelLLM) {
	model.Run()
	b.Models = append(b.Models, model)
	b.total++
}

func (b *LoadBalancer) Enqueue(t task) {
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
