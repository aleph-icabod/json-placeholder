package entity

import "encoding/json"

// Photo contains the info available on the json-placeholder API for a Photo entity
type Photo struct {
	ID           int    `json:"id"`
	AlbumID      int    `json:"albumId"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnailUrl"`
}

func (p *Photo) ToJsonString() string {
	data, _ := json.Marshal(p)
	return string(data)
}
