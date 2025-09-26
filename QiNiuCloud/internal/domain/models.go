package domain

type ModelsInfo struct {
	Thumbnail     string `json:"thumbnail"`
	Url           string `json:"url"`
	DownloadCount int    `json:"download_count"`
	LikeCount     int    `json:"like_count"`
}
