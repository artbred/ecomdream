package versions

type TrainVersionResponse struct {
	Code int `json:"code"`
	Message string `json:"message"`
	VersionID string `json:"version_id"`
}

type VersionInfoResponse struct {
	Code int `json:"code"`
	IsReady bool `json:"is_ready"`
	TimeTraining string `json:"time_training"`
	Info interface{}`json:"info"`
}
