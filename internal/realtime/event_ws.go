package realtime

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
