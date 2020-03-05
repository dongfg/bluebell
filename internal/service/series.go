package service

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/dongfg/bluebell/internal/config"
	"github.com/dongfg/bluebell/internal/payload"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// SeriesService ref
type SeriesService struct {
	opts *SeriesServiceOptions
}

// SeriesServiceOptions service dependency
type SeriesServiceOptions struct {
	Conf *config.Config
}

type seriesSearchResp struct {
	Data []struct {
		ItemID string `json:"itemid"`
		Title  string `json:"title"`
		Poster string `json:"poster"`
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

// NewSeriesService instance
func NewSeriesService(opts *SeriesServiceOptions) *SeriesService {
	return &SeriesService{opts}
}

// Search series
func (svc *SeriesService) Search(query payload.SeriesSearchQuery) ([]payload.Series, error) {
	domain := svc.opts.Conf.Series.Domain
	url := fmt.Sprintf("http://%s/search/api?keyword=%s&type=resource", domain, query.Keyword)

	r, err := http.Get(url)
	if err != nil {
		log.Printf("无法访问资源服务器: %s, %s\n", domain, err.Error())
		return nil, fmt.Errorf("无法访问资源服务器: %s, %s", domain, err.Error())
	}
	defer r.Body.Close()

	var jsonResp seriesSearchResp
	if err := json.NewDecoder(r.Body).Decode(&jsonResp); err != nil {
		log.Printf("无法解析资源服务器返回的数据: %s\n", err.Error())
		return nil, fmt.Errorf("无法解析资源服务器返回的数据: %s", err.Error())
	}

	var seriesSearch = make([]payload.Series, len(jsonResp.Data))

	for i, item := range jsonResp.Data {
		series := payload.Series{}
		series.ID = item.ItemID
		series.CnName = item.Title
		// convert to large image
		series.Poster = strings.ReplaceAll(item.Poster, "s_", "")
		if query.Details {
			if err := svc.seriesFill(&series); err != nil {
				continue
			}
		}

		seriesSearch[i] = series
	}
	return seriesSearch, nil
}

// Detail of series
func (svc *SeriesService) Detail(seriesID string) (payload.Series, error) {
	series := payload.Series{
		ID: seriesID,
	}
	err := svc.seriesFill(&series)
	return series, err
}

// Episodes of series
func (svc *SeriesService) Episodes(seriesID string) ([]payload.SeriesEpisode, error) {
	series := payload.Series{
		ID: seriesID,
	}
	if err := svc.seriesFill(&series); err != nil {
		log.Printf("无法获取剧集详情: %s, %s\n", seriesID, err.Error())
		return nil, err
	}
	if series.RssLink == "" {
		log.Printf("rssLink not found: %s\n", seriesID)
		return nil, errors.New("rssLink not found")
	}

	r, err := http.Get(series.RssLink)
	if err != nil {
		log.Printf("rssLink 无法访问: %s\n", series.RssLink)
		return nil, fmt.Errorf("rssLink 无法访问: %s", series.RssLink)
	}
	defer r.Body.Close()

	var xmlResp seriesEpisodesResp

	if err := xml.NewDecoder(r.Body).Decode(&xmlResp); err != nil {
		log.Printf("无法解析rss返回的数据: %s, %s\n", series.RssLink, err.Error())
		return nil, fmt.Errorf("无法解析rss返回的数据: %s, %s", series.RssLink, err.Error())
	}

	var seriesEpisodes = make([]payload.SeriesEpisode, len(xmlResp.Channel.Items))

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
		seriesEpisodes[i] = payload.SeriesEpisode{
			SeriesID: seriesID,
			Name:     item.Title,
			Season:   season,
			Episode:  episode,
			Ed2k:     item.Ed2k,
			Magnet:   item.Magnet,
		}
	}

	return seriesEpisodes, nil
}

func (svc *SeriesService) seriesFill(series *payload.Series) error {
	domain := svc.opts.Conf.Series.Domain
	if series.ID == "" {
		return errors.New("series id must not empty")
	}
	playStatus := make(chan string, 1)
	// get playStatus from api
	go func() {
		url := fmt.Sprintf("http://%s/resource/index_json/rid/%s/channel/tv", domain, series.ID)

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
	url := fmt.Sprintf("http://%s/resource/%s", domain, series.ID)
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return err
	}
	cnNameHTML, err := doc.Find(".resource-tit h2").Html()
	if err != nil {
		return errors.New("无法获取中文名")
	}
	series.CnName = cnNameHTML[strings.Index(cnNameHTML, "《")+len("《") : strings.Index(cnNameHTML, "》")]
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
