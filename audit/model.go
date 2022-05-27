package audit

import "time"

type Log struct {
	Timestamp time.Time
	Action    string
	Actor     string
	Data      interface{}
	Metadata  interface{}
}
