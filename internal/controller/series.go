package controller

import (
	"github.com/dongfg/bluebell/internal/payload"
	"github.com/dongfg/bluebell/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type seriesController struct {
	opts *seriesControllerOptions
}

type seriesControllerOptions struct {
	routerGroup   *gin.RouterGroup
	seriesService *service.SeriesService
}

func newSeriesController(opts *seriesControllerOptions) {
	ctrl := &seriesController{
		opts,
	}
	routerGroup := ctrl.opts.routerGroup

	routerGroup.GET("", ctrl.seriesSearch)
	routerGroup.GET("/:seriesId", ctrl.seriesDetail)
	routerGroup.GET("/:seriesId/episodes", ctrl.seriesEpisodes)
}

func (ctrl *seriesController) seriesSearch(c *gin.Context) {
	seriesService := ctrl.opts.seriesService
	var query payload.SeriesSearchQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, failed(err.Error()))
	}

	seriesSearch, err := seriesService.Search(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failed(err.Error()))
		return
	}

	c.JSON(http.StatusOK, data(seriesSearch))
}

func (ctrl *seriesController) seriesDetail(c *gin.Context) {
	seriesService := ctrl.opts.seriesService
	seriesID := c.Param("seriesId")
	if seriesID == "" {
		c.JSON(http.StatusBadRequest, failed("missing path param 'seriesId'"))
		return
	}

	series, err := seriesService.Detail(seriesID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failed(err.Error()))
		return
	}

	c.JSON(http.StatusOK, series)
}

func (ctrl *seriesController) seriesEpisodes(c *gin.Context) {
	seriesService := ctrl.opts.seriesService
	seriesID := c.Param("seriesId")
	if seriesID == "" {
		c.JSON(http.StatusBadRequest, failed("missing path param 'seriesId'"))
		return
	}

	seriesEpisodes, err := seriesService.Episodes(seriesID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failed(err.Error()))
		return
	}

	c.JSON(http.StatusOK, data(seriesEpisodes))
}
