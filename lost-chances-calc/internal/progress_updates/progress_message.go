package progressupdates

type ProgressMessage struct {
	RequestID string `json:"requestID"`
	Progress  int    `json:"progress"`
}
