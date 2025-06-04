package magazines

import (
	"time"
)

type Magazine struct {
	Title       string    `csv:"title"`
	ISBN        string    `csv:"isbn"`
	Authors     []string  `csv:"authors"`
	PublishedAt time.Time `csv:"publishedAt"`
}
