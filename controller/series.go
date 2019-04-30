package controller

import (
	"fmt"
	"github.com/dongfg/bluebell/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type seriesController struct {
}

type series struct {
	Id     string
	CnName string
	EnName string
}

type seriesEpisodes struct {
	SeriesId string
}

func newSeriesController(g *gin.RouterGroup) {
	c := &seriesController{}
	g.GET("", c.seriesSearch)
	g.GET("/:seriesId", c.seriesDetail)
	g.GET("/:seriesId/episodes", c.seriesEpisodes)
}

func (controller *seriesController) seriesSearch(c *gin.Context) {
	keyword := c.Query("keyword")

	domain := config.Basic.Series.Domain
	url := fmt.Sprintf("http://%s/search/api?keyword=%s&type=resource", domain, keyword)
	r, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusOK, []series{})
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	c.JSON(http.StatusOK, []series{
		{
			CnName: keyword,
		},
	})
}

func (controller *seriesController) seriesDetail(c *gin.Context) {
	seriesId := c.Param("seriesId")

	c.JSON(http.StatusOK, series{
		Name: seriesId,
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
