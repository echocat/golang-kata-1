package books

type Book struct {
	Title       string   `csv:"title"`
	ISBN        string   `csv:"isbn"`
	Authors     []string `csv:"authors"`
	Description string   `csv:"description"`
}
