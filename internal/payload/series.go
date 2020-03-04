package payload

// Series detail
type Series struct {
	ID         string `json:"id"`
	CnName     string `json:"cnName"`
	Poster     string `json:"poster"`
	EnName     string `json:"enName,omitempty"`
	Link       string `json:"link,omitempty"`
	RssLink    string `json:"rssLink,omitempty"`
	PlayStatus string `json:"playStatus,omitempty"`
	Area       string `json:"area,omitempty"`
	Category   string `json:"category,omitempty"`
}

// SeriesEpisode detail
type SeriesEpisode struct {
	SeriesID string `json:"seriesId"`
	Name     string `json:"name"`
	Season   int    `json:"season"`
	Episode  int    `json:"episode"`
	Ed2k     string `json:"ed2k,omitempty"`
	Magnet   string `json:"magnet,omitempty"`
}

// SeriesSearchQuery parameter
type SeriesSearchQuery struct {
	Keyword string `form:"keyword"`
	Details bool   `form:"details,default=false"`
}
