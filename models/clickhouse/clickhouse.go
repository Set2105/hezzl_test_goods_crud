package clickhouse

import (
	"encoding/json"
	"time"
)

type GoodsLog struct {
	Id          int64     `json:"id"`
	ProjectId   int64     `json:"project_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int64     `json:"priority"`
	Removed     bool      `json:"removed"`
	EventTime   time.Time `json:"event_time"`
}

func (gl *GoodsLog) MarshalJSON() ([]byte, error) {
	type Alias GoodsLog
	return json.Marshal(&struct {
		*Alias
		EventTime string `json:"event_time"`
	}{
		Alias:     (*Alias)(gl),
		EventTime: gl.EventTime.Format("2006-01-02 15:04:05"),
	})
}
