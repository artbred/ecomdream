package versions

type TrainVersionResponse struct {
	Code int `json:"code"`
	Message string `json:"message"`
	VersionID string `json:"version_id"`
}

type IsReadyResponse struct {
	Code int `json:"code"`
	IsReady bool `json:"is_ready"`
}
