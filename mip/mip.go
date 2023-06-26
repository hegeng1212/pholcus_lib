package pholcus_lib

// 基础包
import (
	"github.com/hegeng1212/pholcus/app/downloader/request" //必需
	. "github.com/hegeng1212/pholcus/app/spider"           //必需
	//"github.com/hegeng1212/pholcus/common/goquery"         //DOM解析
	"github.com/hegeng1212/pholcus/logs"                   //信息输出
	"github.com/hegeng1212/pholcus/app/downloader/surfer/agent"
	// . "github.com/hegeng1212/pholcus/app/spider/common"          //选用

	// net包
	// "net/http" //设置http.Header
	// "net/url"

	// 编码包
	// "encoding/xml"
	// "encoding/json"

	// 字符串处理包

	"strconv"
	"strings"

	// 其他包
	"fmt"
	"math/rand"

	"net/http"
	"time"
)

func init() {
	Mip.Register()
}

var Mip = &Spider{
	Name:        "百度健康问答",
	Description: "百度健康问答 [www.baidu.com]",
	// Pausetime: 300,
	Keyin:        KEYIN,
	Limit:        LIMIT,
	EnableCookie: true,
	// 禁止输出默认字段 Url/ParentUrl/DownloadTime
	NotDefaultField: true,
	// 命名空间相对于数据库名，不依赖具体数据内容，可选
	Namespace: nil,
	// 子命名空间相对于表名，可依赖具体数据内容，可选
	SubNamespace: nil,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.Aid(map[string]interface{}{"loop": [2]int{0, 1}, "Rule": "生成请求"}, "生成请求")
		},

		Trunk: map[string]*Rule{

			"生成请求": {
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					var duplicatable bool
					for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {
						if loop[0] == 0 {
							duplicatable = true
						} else {
							duplicatable = false
						}
						header := make(http.Header)
						l := len(agent.UserAgents["common"])
						r := rand.New(rand.NewSource(time.Now().UnixNano()))
						header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")//agent.UserAgents["common"][r.Intn(l)])
						header.Add("Connection", "keep-alive")
						header.Add("Accept-Encoding", "gzip, deflate, br")
						header.Add("Accept", "*/*")
						header.Add("Host", "www.baidu.com")
						fmt.Println(fmt.Printf("baidu header %#v %#v %#v ", header, l, r))
						ctx.AddQueue(&request.Request{
							Url:        "https://wapask-mip.39.net/bdsshz/question/"+ctx.GetKeyin()+".html?v=" + strconv.Itoa(50*loop[0]),
							Rule:       aid["Rule"].(string),
							Reloadable: duplicatable,
							Header: header,
						})
					}
					return nil
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					title := query.Find(".question-box >h2").Text()
					if title == "" {
						logs.Log.Critical("[消息提示：| 任务：%v | KEYIN：%v | 规则：%v] 没有抓取到任何数据！!!\n", ctx.GetName(), ctx.GetKeyin(), ctx.GetRuleName())
						return
					}
					// 调用指定规则下辅助函数
					ctx.Aid(map[string]interface{}{"loop": [2]int{1, 1}, "Rule": "搜索结果"})
					// 用指定规则解析响应流
					ctx.Parse("搜索结果")
				},
			},

			"搜索结果": {
				//注意：有无字段语义和是否输出数据必须保持一致
				ItemFields: []string{
					"标题",
					"内容",
					"跳转",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					title := query.Find(".question-box >h2").Text()
					content := query.Find(".doctor-replay-text >p").Text()
					href := ctx.Request.Url
					fmt.Println(fmt.Printf("mip 搜索结果 %s %#v %#v ", query, title, content))
					// 结果存入Response中转
					ctx.Output(map[int]interface{}{
						0: strings.Replace(strings.Trim(title, " \t\n"), ",", "，", 0),
						1: strings.Replace(strings.Trim(content, " \t\n"), ",", "，", 0),
						2: href,
					})
				},
			},
		},
	},
}
