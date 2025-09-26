package dao

type Models struct {
	Token   string      `bson:"token" json:"token"`
	Ctime   int64       `bson:"ctime" json:"ctime"`
	Context []ModelInfo `bson:"context" json:"context"`
}

type ModelInfo struct {
	Hash                string `bson:"hash" json:"hash"`
	Score               int64  `bson:"score" json:"score"`
	NCloseAfterDownload int64  `bson:"n_close_after_download" json:"n_close_after_download"`
	DownloadCount       int    `json:"download_count" bson:"download_count"`
	LikeCount           int    `bson:"like_count" json:"like_count"`
	Ctime               int64  `bson:"ctime" json:"ctime"`
	LastDownloadTime    int64  `bson:"last_download_time" json:"last_download_time"`
	Thumbnail           string `json:"thumbnail"`
	Url                 string `bson:"url" json:"url"`
}
