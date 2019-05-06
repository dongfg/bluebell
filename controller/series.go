package controller

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/dongfg/bluebell/config"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
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

type seriesEpisode struct {
	SeriesId string
	Name     string
	Season   int
	Episode  int
	Ed2k     string `json:",omitempty"`
	Magnet   string `json:",omitempty"`
}

type seriesSearchQuery struct {
	Keyword string `form:"keyword"`
	Details bool   `form:"details,default=false"`
}

type seriesSearchResp struct {
	Data []struct {
		ItemId string `json:"itemid"`
		Title  string
		Poster string
	}
}

type seriesPlayStatusResp struct {
	PlayStatus string `json:"play_status"`
}

type seriesEpisodesResp struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		Items []struct {
			Title  string `xml:"title"`
			Ed2k   string `xml:"ed2k"`
			Magnet string `xml:"magnet"`
		} `xml:"item"`
	} `xml:"channel"`
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
		c.JSON(http.StatusInternalServerError, failed(fmt.Sprintf("资源服务器不可用: %s", err.Error())))
		return
	}
	defer r.Body.Close()

	var jsonResp seriesSearchResp
	if err := json.NewDecoder(r.Body).Decode(&jsonResp); err != nil {
		c.JSON(http.StatusInternalServerError, failed(fmt.Sprintf("无法解析资源服务器返回的数据: %s", err.Error())))
		return
	}

	var seriesSearch = make([]series, len(jsonResp.Data))

	for i, item := range jsonResp.Data {
		series := series{}
		series.Id = item.ItemId
		series.CnName = item.Title
		// convert to large image
		series.Poster = strings.ReplaceAll(item.Poster, "s_", "")
		if query.Details {
			if err := seriesFill(&series); err != nil {
				continue
			}
		}

		seriesSearch[i] = series
	}

	c.JSON(http.StatusOK, data(seriesSearch))
}

func (controller *seriesController) seriesDetail(c *gin.Context) {
	seriesId := c.Param("seriesId")
	if seriesId == "" {
		c.JSON(http.StatusBadRequest, failed("missing path param 'seriesId'"))
		return
	}

	series := series{
		Id: seriesId,
	}

	if err := seriesFill(&series); err != nil {
		c.JSON(http.StatusInternalServerError, failed(err.Error()))
		return
	}

	c.JSON(http.StatusOK, series)
}

func (controller *seriesController) seriesEpisodes(c *gin.Context) {
	seriesId := c.Param("seriesId")
	if seriesId == "" {
		c.JSON(http.StatusBadRequest, failed("missing path param 'seriesId'"))
		return
	}

	series := series{}
	series.Id = seriesId
	if err := seriesFill(&series); err != nil || series.RssLink == "" {
		c.JSON(http.StatusInternalServerError, failed("rssLink not found"))
		return
	}

	r, err := http.Get(series.RssLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failed(fmt.Sprintf("rssLink无法访问: %s", err.Error())))
		return
	}
	defer r.Body.Close()

	var xmlResp seriesEpisodesResp

	if err := xml.NewDecoder(r.Body).Decode(&xmlResp); err != nil {
		c.JSON(http.StatusInternalServerError, failed(fmt.Sprintf("无法解析rss返回的数据: %s", err.Error())))
		return
	}

	var seriesEpisodes = make([]seriesEpisode, len(xmlResp.Channel.Items))

	seasonEpisodeParse := func(name string) (int, int) {
		re := regexp.MustCompile(`(?m)[Ss](\d{1,3})[Ee](\d{1,3})`)
		matches := re.FindStringSubmatch(name)
		if len(matches) != 3 {
			return -1, -1
		}
		season, err := strconv.Atoi(matches[1])
		if err != nil {
			season = -1
		}
		episode, err := strconv.Atoi(matches[2])
		if err != nil {
			episode = -1
		}
		return season, episode
	}

	for i, item := range xmlResp.Channel.Items {
		season, episode := seasonEpisodeParse(item.Title)
		seriesEpisodes[i] = seriesEpisode{
			SeriesId: seriesId,
			Name:     item.Title,
			Season:   season,
			Episode:  episode,
			Ed2k:     item.Ed2k,
			Magnet:   item.Magnet,
		}
	}

	c.JSON(http.StatusOK, data(seriesEpisodes))
}

func seriesFill(series *series) error {
	if series.Id == "" {
		return errors.New("series id must not empty")
	}
	playStatus := make(chan string, 1)
	// get playStatus from api
	go func() {
		url := fmt.Sprintf("http://%s/resource/index_json/rid/%s/channel/tv", config.Basic.Series.Domain, series.Id)

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
	url := fmt.Sprintf("http://%s/resource/%s", config.Basic.Series.Domain, series.Id)
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return err
	}
	cnNameHtml, err := doc.Find(".resource-tit h2").Html()
	if err != nil {
		return errors.New("无法获取中文名")
	}
	series.CnName = cnNameHtml[strings.Index(cnNameHtml, "《")+len("《") : strings.Index(cnNameHtml, "》")]
	series.EnName = doc.Find(".resource-con .fl-info li:nth-child(1) > strong").Text()
	series.RssLink = doc.Find(".resource-tit h2 a").AttrOr("href", "")
	series.Area = doc.Find(".resource-con .fl-info li:nth-child(2) > strong").Text()
	series.Category = doc.Find(".resource-con .fl-info li:nth-child(6) > strong").Text()
	if series.Poster == "" {
		series.Poster = doc.Find(".resource-con > div.fl-img > div.imglink > a").AttrOr("href", "")
	}
	series.Link = url
	return nil
}
