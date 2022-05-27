package main

import (
	"encoding/json"
	"github.com/gocolly/colly/v2"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type Item struct {
	// 图片
	Img string
	// 地址
	Link string
	// 标题
	Title string
	// 价格
	Price string
	// 描述
	Describe string
	// 赞
	TUp string
	// 踩
	TDown string
	// 收藏
	Collect string
	// 评论
	Comment string
	// 发布时间
	PublishTime string
	// 平台
	Platform string
	// 购买连接
	BuyUrl string
	// 元数据
	MetaData MetaData
}

type MetaData struct {
	Name        string `json:"name"`
	ID          string `json:"id"`
	Price       int    `json:"price"`
	Brand       string `json:"brand"`
	Category    string `json:"category"`
	Dimension12 string `json:"dimension12"`
	Metric1     int    `json:"metric1"`
	Dimension10 string `json:"dimension10"`
	Dimension9  string `json:"dimension9"`
	Dimension11 int    `json:"dimension11"`
	Dimension20 string `json:"dimension20"`
	Dimension64 string `json:"dimension64"`
	Quantity    int    `json:"quantity"`
	ChannelID   int    `json:"channel_id"`
	CateLevel1  string `json:"cate_level1"`
	Channel     string `json:"channel"`
}

type Query struct {
	Name     string `json:"name,omitempty" query:"name"`
	Page     string `json:"page,omitempty" query:"page"`
	MinPrice string `json:"min_price,omitempty" query:"min_price"`
	MaxPrice string `json:"max_price,omitempty" query:"max_price"`
}

func main() {
	web()
}

func web() {
	app := fiber.New()

	app.Get("/:name", func(ctx *fiber.Ctx) error {
		q := new(Query)
		if err := ctx.QueryParser(q); err != nil {
			return err
		}
		q.Name = ctx.Params("name")
		return ctx.JSON(find(q))
	})

	app.Listen(":3000")

	// https://wxpusher.dingliqc.com/docs/#/?id=%e5%8f%91%e9%80%81%e6%b6%88%e6%81%af-1
	// 消息推送

}

func find(q *Query) []Item {

	var list = make([]Item, 0)

	if q.Name == "" {
		return list
	}
	url := "https://search.smzdm.com/?c=faxian&v=b" + "&s=" + q.Name
	if q.Page != "" {
		url += "&p=" + q.Page
	}
	if q.MaxPrice != "" {
		url += "&max_price=" + q.MaxPrice
	}
	if q.MinPrice != "" {
		url += "&min_price=" + q.MinPrice
	}
	c := colly.NewCollector()
	// Find and visit all links
	c.OnHTML("#feed-main-list .feed-row-wide", func(e *colly.HTMLElement) {
		i := Item{}
		dataStr := e.ChildAttr(".feed-link-btn-inner a", "onclick")
		i.MetaData = parseMetaData(dataStr)
		// 图片
		i.Img = e.ChildAttr(".z-feed-img img", "src")
		// link
		i.Link = e.ChildAttr(".z-feed-img a", "href")

		// 标题
		i.Title = e.ChildText(".z-feed-content .feed-block-title .feed-nowrap")
		// 价格
		i.Price = e.ChildText(".z-feed-content .feed-block-title .z-highlight")
		// 描述
		i.Describe = e.ChildText(".z-feed-content .feed-block-descripe-top")

		// 点赞 踩
		thumbs := e.ChildTexts(".z-feed-foot-l .unvoted-wrap span")
		if thumbs != nil && len(thumbs) > 0 {
			i.TUp = thumbs[0]
			i.TDown = thumbs[1]
		}
		// 收藏
		i.Collect = e.ChildText(".z-feed-foot-l .feed-btn-fav span")
		// 评论
		i.Comment = e.ChildText(".z-feed-foot-l .feed-btn-comment")

		// 发布时间
		i.PublishTime = e.ChildText(".feed-block-extras")
		i.Platform = e.ChildText(".feed-block-extras span")
		i.BuyUrl = e.ChildAttr(".feed-link-btn-inner a", "href")

		list = append(list, i)
	})

	if err := c.Visit(url); err != nil {
		return list
	}
	c.Wait()
	return list
}

func parseMetaData(str string) MetaData {
	str = strings.Replace(str, ";gtmAddToCart", "", -1)
	str = strings.Replace(str, "'", "\"", -1)
	str = strings.Replace(str, "({", "{", -1)
	str = strings.Replace(str, "})", "}", -1)

	str = str[:strings.Index(str, "}")+1]
	//Json.u

	var list MetaData
	json.Unmarshal([]byte(str), &list)
	return list
}
