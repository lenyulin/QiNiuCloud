package ModelGnerationResultHelper

type ModelsInfo struct {
	Token                     string `json:"token"`
	Hash                      string `json:"hash"`
	Thumbnail                 string `json:"thumbnail"`
	Url                       string `json:"url"`
	DownloadCount             int    `json:"download_count"`
	LikeCount                 int    `json:"like_count"`
	CloseAfterDownloadedCount int    `json:"close_after_downloaded_count"`
}
type ModelGenerationTaskResult struct {
	JobId string `json:"job_id"`
	Token string `json:"token"`
	Url   string `json:"url"`
	Thumb string `json:"thumb"`
}
