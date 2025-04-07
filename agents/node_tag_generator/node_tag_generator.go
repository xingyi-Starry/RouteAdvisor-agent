package node_tag_generator

import (
	"context"
	"fmt"
	"route-advisor-agent/utils/ai"

	"encoding/json"

	"github.com/sirupsen/logrus"
)

const prompt = `你是一个地理信息系统专家，负责为给定的地理节点生成标签。请根据以下节点列表中节点的名称生成包含标签的目标节点信息。标签只能在以下几项中选择：["美食", "景点", "购物", "交通", "住宿", "娱乐", "文化", "教育", "医疗", "体育", "default"]，其中"default"仅在节点名字为空或者无其它合适选项时使用。请确保生成的标签与节点名称相关，并且一个节点只包含一个标签。

注意，请不要添加任何额外的文本或解释，除了新增的标签以外不要修改其它信息，只返回 JSON 格式的目标节点列表信息，并保持顺序不变。确保输出的 JSON 格式正确，并且包含所有必要的字段。请遵循以下示例格式：

<输入示例>
[{"id":6146512085,"name":"麦当劳","lat":39.9622611,"lon":116.3513657},{"id":6041100036,"name":"楼上楼茶餐厅","lat":39.961991,"lon":116.352893},{"id":8810354632,"name":"音乐喷泉","lat":39.9598892,"lon":116.3516116},{"id":3511264386,"name":"南区超市","lat":39.9581307,"lon":116.3510136}]
<输入示例 />
<对应输出示例>
[{"id":6146512085,"name":"麦当劳","lat":39.9622611,"lon":116.3513657,"tag":"美食"},{"id":6041100036,"name":"楼上楼茶餐厅","lat":39.961991,"lon":116.352893,"tag":"美食"},{"id":8810354632,"name":"音乐喷泉","lat":39.9598892,"lon":116.3516116,"tag":"娱乐"},{"id":3511264386,"name":"南区超市","lat":39.9581307,"lon":116.3510136,"tag":"购物"}]
<对应输出示例 />
`

type BasicNode struct {
	Id   int     `json:"id"`
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

type TargetNode struct {
	Id   int     `json:"id"`
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
	Tag  string  `json:"tag"`
}

type NodeTagGenerator struct {
	modelConfig *ai.ModelConfig
}

func NewNodeTagGenerator(modelConfig *ai.ModelConfig) *NodeTagGenerator {
	return &NodeTagGenerator{
		modelConfig: modelConfig,
	}
}

func (agent *NodeTagGenerator) GenerateTag(ctx context.Context, nodes []BasicNode) ([]TargetNode, error) {
	src, err := json.Marshal(nodes)
	if err != nil {
		logrus.Error(err)
		return []TargetNode{}, err
	}

	question := ai.AiQuestion{
		SystemPrompt: prompt,
		History:      []ai.HistoryEntry{},
		Question:     string(src),
	}

	// 调用 LLM 进行推理
	questionChan, err := ai.InvokeLLM(ctx, question, *agent.modelConfig)
	if err != nil {
		return []TargetNode{}, err
	}

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
		if src[i].Id != result[i].Id || src[i].Name != result[i].Name || src[i].Lat != result[i].Lat || src[i].Lon != result[i].Lon {
			logrus.Warnf("the src and result is mismatch, mismatch entries have been corrected.\nsrc: %v,\nresult: %v\n", src[i], result[i])
			result[i].Id = src[i].Id
			result[i].Name = src[i].Name
			result[i].Lat = src[i].Lat
			result[i].Lon = src[i].Lon
		}
		if result[i].Tag == "" {
			logrus.Warnf("the tag is empty, set it to default.\nsrc: %v,\nresult: %v\n", src[i], result[i])
			result[i].Tag = "default"
		}
	}

	return nil
}
