package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/dongfg/bluebell/config"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
)

type seriesController struct {
}

type series struct {
	Id         string
	CnName     string
	Poster     string
	EnName     string `json:",omitempty"`
	Link       string `json:",omitempty"`
	RssLink    string `json:",omitempty"`
	PlayStatus string `json:",omitempty"`
	Area       string `json:",omitempty"`
	Category   string `json:",omitempty"`
}

type seriesEpisodes struct {
	SeriesId string
}

type seriesSearchQuery struct {
	Keyword string `form:"keyword"`
	Details bool   `form:"details,default=false"`
}

type seriesSearchResp struct {
	Data [] struct {
		ItemId string `json:"itemid"`
		Title  string
		Poster string
	}
}

type seriesPlayStatusResp struct {
	PlayStatus string `json:"play_status"`
}

func newSeriesController(g *gin.RouterGroup) {
	c := &seriesController{}
	g.GET("", c.seriesSearch)
	g.GET("/:seriesId", c.seriesDetail)
	g.GET("/:seriesId/episodes", c.seriesEpisodes)
}

func (controller *seriesController) seriesSearch(c *gin.Context) {
	var query seriesSearchQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, failed(err.Error()))
	}

	url := fmt.Sprintf("http://%s/search/api?keyword=%s&type=resource", config.Basic.Series.Domain, query.Keyword)

	r, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("资源服务器不可用: %s", err.Error()))
		return
	}
	defer r.Body.Close()

	var jsonResp seriesSearchResp
	if err := json.NewDecoder(r.Body).Decode(&jsonResp); err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("无法解析资源服务器返回的数据: %s", err.Error()))
		return
	}

	var seriesSearch = make([]series, len(jsonResp.Data))

	for i, item := range jsonResp.Data {
		series := series{}
		series.CnName = item.Title
		if query.Details {
			series, err = seriesById(item.ItemId)
			if err != nil {
				continue
			}
		}
		series.Id = item.ItemId
		// convert to large image
		series.Poster = strings.ReplaceAll(item.Poster, "s_", "")

		seriesSearch[i] = series
	}

	c.JSON(http.StatusOK, data(seriesSearch))
}

func (controller *seriesController) seriesDetail(c *gin.Context) {
	seriesId := c.Param("seriesId")

	c.JSON(http.StatusOK, series{
		CnName: seriesId,
	})
}

func (controller *seriesController) seriesEpisodes(c *gin.Context) {
	seriesId := c.Param("seriesId")

	c.JSON(http.StatusOK, []seriesEpisodes{
		{
			SeriesId: seriesId,
		},
	})
}

func seriesById(id string) (series, error) {
	playStatus := make(chan string, 1)
	// get playStatus from api
	go func() {
		url := fmt.Sprintf("http://%s/resource/index_json/rid/%s/channel/tv", config.Basic.Series.Domain, id)

		r, err := http.Get(url)
		if err != nil {
			playStatus <- "无法获取连载状态"
			return
		}
		defer r.Body.Close()

		rawResp, err := ioutil.ReadAll(r.Body)
		if err != nil {
			playStatus <- "无法获取连载状态"
			return
		}

		// remove 'var index_info='
		stringResp := string(rawResp)[len("var index_info="):]
		var playStatusResp seriesPlayStatusResp
		err = json.Unmarshal([]byte(stringResp), &playStatusResp)
		if err != nil {
			playStatus <- "无法获取连载状态"
			return
		}
		playStatus <- playStatusResp.PlayStatus
	}()

	// get detail from parse page
	url := fmt.Sprintf("http://%s/resource/%s", config.Basic.Series.Domain, id)
	r, err := http.Get(url)
	if err != nil {
		return series{}, err
	}
	defer r.Body.Close()

	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return series{}, err
	}
	cnNameHtml, err := doc.Find(".resource-tit h2").Html()
	if err != nil {
		return series{}, errors.New("无法获取中文名")
	}
	cnName := cnNameHtml[strings.Index(cnNameHtml, "《")+len("《") : strings.Index(cnNameHtml, "》")]
	enName := doc.Find(".resource-con .fl-info li:nth-child(1) > strong").Text()
	rssLink := doc.Find(".resource-tit h2 a").AttrOr("href", "")
	area := doc.Find(".resource-con .fl-info li:nth-child(2) > strong").Text()
	category := doc.Find(".resource-con .fl-info li:nth-child(6) > strong").Text()

	return series{
		Id:         id,
		CnName:     cnName,
		EnName:     enName,
		Link:       url,
		RssLink:    rssLink,
		PlayStatus: <-playStatus,
		Area:       area,
		Category:   category,
	}, nil
}
