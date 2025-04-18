package node_tag_generator

import (
	"context"
	"encoding/json"
	"fmt"

	"route-advisor-agent/utils/ai"

	"github.com/sirupsen/logrus"
)

const prompt = `你是一个地理信息系统专家，负责为给定的地理节点补全信息。给定的地理节点都位于首钢园。需要补全的信息有类型(type)、标签组(tags)、描述(description)、热度(heat)、评分(rate)。类型只能在以下几项中选择：["美食", "景点", "购物", "交通", "住宿", "娱乐", "厕所"]，标签可以是任意简短的符合特征的中文描述，如"咖啡"、"文化"、"户外"、"家庭"等，每个节点最少需要3个标签。描述是对于该节点的详细描述，不超过50字。热度是一个整数，范围在1到10000之间，表示该节点的热度。评分是一个浮点数，范围在0到5之间，保留1位小数，表示该节点的评分。

注意，除了新增的字段以外不要修改其它信息，请不要添加任何额外的文本或解释，只返回 JSON 格式的目标节点列表信息，并保持顺序不变。确保输出的 JSON 格式正确，并且包含所有必要的字段。请遵循以下示例格式：

<输入示例>
[{"id":130,"node":12064974127,"name":"全家"},{"id":131,"node":12064974128,"name":"五一剧场"},{"id":132,"node":12064974129,"name":"瑞幸咖啡"}]
<输入示例 />
<对应输出示例>
[{"id":130,"node":12064974127,"name":"全家","type":"购物","tags":["便利店","连锁","实惠"],"description":"景区内的便利店，提供日常生活用品和食品饮料。","heat":2514,"rate":4.2},
{"id":131,"node":12064974128,"name":"五一剧场","type":"娱乐","tags":["演出","文化","艺术"],"description":"历史悠久的剧场，经常举办各类文艺演出和活动。","heat":197,"rate":4.7},
{"id":132,"node":12064974129,"name":"瑞幸咖啡","type":"美食","tags":["咖啡","连锁","休闲"],"description":"提供各类咖啡饮品和轻食的连锁咖啡店。","heat":7649,"rate":4.5}]
<对应输出示例 />

请根据以上示例补全提供的地理节点信息，返回完整的JSON格式数据。
`

type BasicNode struct {
	Id   int    `json:"id"`
	Node int    `json:"node"`
	Name string `json:"name"`
}

type TargetNode struct {
	Id          int      `json:"id"`
	Node        int      `json:"node"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
	Heat        int      `json:"heat"`
	Rate        float32  `json:"rate"`
}

type NodeTagGenerator struct {
	model ai.LLM
}

func NewNodeTagGenerator(model ai.LLM) *NodeTagGenerator {
	return &NodeTagGenerator{
		model: model,
	}
}

func (agent *NodeTagGenerator) GenerateTag(ctx context.Context, nodes []BasicNode) ([]TargetNode, error) {
	src, err := json.Marshal(nodes)
	if err != nil {
		logrus.Error(err)
		return []TargetNode{}, err
	}

	question := ai.Messages{
		SystemPrompt: prompt,
		History:      []ai.HistoryEntry{},
		Question:     string(src),
	}

	// 调用 LLM 进行推理
	questionChan := ai.InvokeLLM(ctx, question, agent.model)

	answer := ""
	for msg := range questionChan {
		answer += msg
	}
	logrus.Debug(answer) // debug
	result := []TargetNode{}
	err = json.Unmarshal([]byte(answer), &result)
	if err != nil {
		logrus.Error(err)
		return []TargetNode{}, err
	}

	// validate the result
	err = validate(nodes, result)
	if err != nil {
		logrus.Error(err)
		return []TargetNode{}, err
	}

	return result, nil
}

func validate(src []BasicNode, result []TargetNode) error {
	if len(src) != len(result) {
		return fmt.Errorf("the length of src and result is not equal, src: %d, result: %d", len(src), len(result))
	}

	for i := range src {
		if src[i].Id != result[i].Id || src[i].Node != result[i].Node || src[i].Name != result[i].Name {
			logrus.Warnf("the src and result is mismatch, mismatch entries have been corrected.\nsrc: %v,\nresult: %v\n", src[i], result[i])
			result[i].Id = src[i].Id
			result[i].Name = src[i].Name
		}
		if result[i].Type == "" {
			logrus.Warnf("the type is empty, set it to undefined.\nsrc: %v,\nresult: %v\n", src[i], result[i])
			result[i].Type = "undefined"
		}
		if len(result[i].Tags) == 0 {
			logrus.Warnf("the tag is empty.")
		}
	}

	return nil
}
