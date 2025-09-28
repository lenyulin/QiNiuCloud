package domain

type Interactive struct {
	Token                     string `json:"token"`
	Hash                      string `json:"hash"`
	DownloadCount             int    `json:"download_count"`
	LikeCount                 int    `json:"like_count"`
	CloseAfterDownloadedCount int    `json:"close_after_downloaded_count"`
}
