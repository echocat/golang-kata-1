package authors

type Author struct {
	Email     string `csv:"email"`
	FirstName string `csv:"firstname"`
	LastName  string `csv:"lastname"`
}
