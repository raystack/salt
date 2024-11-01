package audit

import "time"

type Log struct {
	Timestamp time.Time   `json:"timestamp"`
	Action    string      `json:"action"`
	Actor     string      `json:"actor"`
	Data      interface{} `json:"data"`
	Metadata  interface{} `json:"metadata"`
}
