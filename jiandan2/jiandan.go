package pholcus_lib

// 基础包
import (
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	"github.com/henrylee2cn/pholcus/common/goquery"         //DOM解析
	"strconv"
	"fmt"
	"log"
	"net/http"
)

func init() {
	Jiandan2.Register()
}

var Jiandan2 = &Spider{

	Name:        "jiandan",
	Description: "煎蛋图片下载",
	//Pausetime: 300,
	Keyin:           KEYIN,
	Limit:           LIMIT,
	EnableCookie:    false,
	NotDefaultField: true,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.Aid(map[string]interface{}{"loop": [2]int{0, ctx.GetLimit()}, "Rule": "生成请求"}, "生成请求")
		},

		Trunk: map[string]*Rule{
			"生成请求": {
				ItemFields: []string{
					"img",
					"tag",
				},
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					var url string
					log.Println(url, "xx")
					for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {
						url = fmt.Sprintf("http://jandan.net/ooxx/page-%s#comments", strconv.Itoa(loop[0]+1))

						ctx.AddQueue(&request.Request{
							Url:          url,
							DownloaderID: 1,
							Rule:         aid["Rule"].(string),
							Header: http.Header{
								"Content-Type": []string{"text/html; charset=utf-8"},
								"User-Agent": []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.186 Safari/537.36"},
								},
						})
					}
					return nil
				},
				ParseFunc: func(ctx *Context) {
					doc := ctx.GetDom()

					log.Println(doc.Html())
					return

					// Find the review items
					doc.Find(".commentlist li").Each(func(i int, s *goquery.Selection) {

						like, _ := s.Find("a.comment-like + span").Html()
						unlike, _ := s.Find("a.comment-unlike + span").Html()

						if like != "" {
							numLike, _ := strconv.Atoi(like)
							numUnLike, _ := strconv.Atoi(unlike)
							if numUnLike == 0 {
								numUnLike = 1
							}

							if numLike > numUnLike && numLike/numUnLike*100 >= 61 {
								s.Find("a.view_img_link").Each(func(i int, img *goquery.Selection) {
									href, _ := img.Attr("href")


									log.Println("https:" + href)

									src := "https:" + href
									ctx.AddQueue(&request.Request{
										Url:          src,
										Rule:         "下载文件",
										ConnTimeout:  -1,
										DownloaderID: 0,
									})

									ctx.Output(map[int]interface{}{
										0: src,
										1: "jiandan",
									}, "下载文件")
								})
							}
						}
					})

					// 用指定规则解析响应流
					ctx.Parse("下载文件")
				},
			},
			"下载文件": {
				ParseFunc: func(ctx *Context) {
					ctx.FileOutput()
				},
			},
		},
	},
}
