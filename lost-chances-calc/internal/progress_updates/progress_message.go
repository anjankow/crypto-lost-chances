package progressupdates

type progressMessage struct {
	RequestID string `json:"requestID"`
	Progress  int    `json:"progress"`
}
