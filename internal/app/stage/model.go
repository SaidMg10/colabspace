package stage

import (
	"encoding/json"
	"fmt"
	"time"
)

type Stage struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Position  int       `json:"position"`
	WipLimit  int       `json:"wip_limit"`
	StageType int       `json:"stage_type"`
	BoardId   int64     `json:"board_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StageType int

const (
	StageDefault StageType = iota
	StageDone
	StageArchive
)

func (t StageType) String() string {
	switch t {
	case StageDefault:
		return "default"
	case StageDone:
		return "done"
	case StageArchive:
		return "archive"
	default:
		return "unknown"
	}
}

func (t StageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *StageType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case "default":
		*t = StageDefault
	case "done":
		*t = StageDone
	case "archive":
		*t = StageArchive
	default:
		return fmt.Errorf("invalid stage type: %s", s)
	}
	return nil
}
