package publications

import (
	"time"
)

type Publication struct {
	Title       string
	ISBN        string
	Authors     []string
	Type        string
	Description string
	PublishedAt *time.Time
}
