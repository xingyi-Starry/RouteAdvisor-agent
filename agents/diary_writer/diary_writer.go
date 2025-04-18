package diary_writer

import (
	"context"
	"encoding/json"
	"fmt"
	"route-advisor-agent/models"
	"route-advisor-agent/utils/ai"

	"github.com/sirupsen/logrus"
)

const prompt = `你是一个游记写手，请根据以下标题和已有内容生成一篇游记，游记的内容应该包括景点、活动、饮食等方面的信息，并且要有一定的情感色彩。游记的长度应该在500字左右。若标题存在请确保生成的游记内容与标题相关，并且要有一定的逻辑性和连贯性，必要还可以优化标题内容；若不存在则根据内容补全标题。请使用JSON格式返回结果，包含"title"和"content"两个字段。确保输出的 JSON 格式正确，不要添加任何额外的文本或解释，不要用markdown的代码块包裹输出。
<输入示例>
{"title":"美丽的西湖之旅","content":"今天我来到了杭州，准备游览西湖。"}
<输入示例 />
<对应输出示例>
{"title":"美丽的西湖之旅", "今天我来到了杭州，开始了期待已久的西湖之旅。初见西湖，湖面如镜，微风拂面，湖光山色浑然一体，仿佛进入了一幅画。沿着苏堤漫步，春意盎然，桃红柳绿，脚下的青石板路延伸向远方，带来一种时光静止的美好感受。\n\n在断桥残雪处，迎面吹来的湖风带来了丝丝凉意，桥上游客熙熙攘攘，而我选择静静地坐在湖边石凳上，看着偶尔跳跃的鱼儿，感受西湖的宁静与深远。随后，乘一叶小船游湖，画舫轻轻划破水面，船夫的歌声悠扬，仿佛把我带回了古代诗人的浪漫时光。\n\n午餐时，我品尝了杭州著名的美食——西湖醋鱼，酸甜适口，鱼肉鲜嫩；还有龙井虾仁，鲜美的虾仁搭配清香的龙井茶，别有一番风味。西湖周边的小吃摊点也令人流连忘返，桂花糕、麻辣小龙虾，每一口都充满了地道的杭州味道。\n\n下午，我去了灵隐寺，那里香火缭绕，古木参天，带来一片心灵的宁静。夕阳西下时，站在宝石山顶俯瞰整个西湖，湖面泛起金色波光，城市的灯火逐渐亮起，夜幕中的西湖更加迷人。\n\n这一天的西湖之旅令我流连忘返，不仅感受到了自然的美丽，更享受了文化的底蕴与生活的惬意。西湖，真的是一个让人心醉神迷的地方。期待下一次的再访，继续探索这片美丽的土地。"}
<对应输出示例 />
`

type DiaryWriter struct {
	Model ai.LLM
}

func NewDiaryWriter(model ai.LLM) *DiaryWriter {
	return &DiaryWriter{
		Model: model,
	}
}

func (agent *DiaryWriter) CreateDiary(ctx context.Context, config models.DiaryConfig) (models.DiaryConfig, error) {
	q := fmt.Sprintf(`{"title":"%s","content":"%s"}`, config.Title, config.Content)

	question := ai.Messages{
		SystemPrompt: prompt,
		History:      []ai.HistoryEntry{},
		Question:     q,
	}

	// 调用 LLM 接口
	ch := ai.InvokeLLM(ctx, question, agent.Model)

	// 处理返回的结果
	result := ""
	for msg := range ch {
		result += msg
	}
	logrus.Info(result)

	var diaryResult models.DiaryConfig
	err := json.Unmarshal([]byte(result), &diaryResult)
	if err != nil {
		logrus.Error(err)
		return models.DiaryConfig{}, err
	}

	return diaryResult, nil
}
