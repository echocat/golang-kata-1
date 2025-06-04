package books

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type Service interface {
	LoadFromCSV(filePath string) error
	GetAll() []Book
	FindByISBN(isbn string) *Book
	FindByAuthorEmail(email string) []Book
}

type bookService struct {
	books       []Book
	isbnIndex   map[string]*Book
	authorIndex map[string][]*Book
}

func NewBookService() Service {
	return &bookService{
		books:       make([]Book, 0),
		isbnIndex:   make(map[string]*Book),
		authorIndex: make(map[string][]*Book),
	}
}

func (s *bookService) LoadFromCSV(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			slog.Error("failed to close file", "filePath", filePath, "error", err)
		}
	}()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV records: %w", err)
	}

	for i, record := range records {
		if i == 0 {
			continue
		}

		if len(record) < 4 {
			continue
		}

		authorsStr := strings.TrimSpace(record[2])
		var authors []string
		if authorsStr != "" {
			for _, author := range strings.Split(authorsStr, ",") {
				authors = append(authors, strings.TrimSpace(author))
			}
		}

		book := Book{
			Title:       strings.TrimSpace(record[0]),
			ISBN:        strings.TrimSpace(record[1]),
			Authors:     authors,
			Description: strings.TrimSpace(record[3]),
		}

		s.books = append(s.books, book)

		bookPtr := &s.books[len(s.books)-1]

		s.isbnIndex[book.ISBN] = bookPtr

		for _, author := range book.Authors {
			s.authorIndex[author] = append(s.authorIndex[author], bookPtr)
		}
	}

	return nil
}

func (s *bookService) GetAll() []Book {
	return s.books
}

func (s *bookService) FindByISBN(isbn string) *Book {
	return s.isbnIndex[isbn]
}

func (s *bookService) FindByAuthorEmail(email string) []Book {
	bookPtrs := s.authorIndex[email]
	result := make([]Book, len(bookPtrs))
	for i, bookPtr := range bookPtrs {
		result[i] = *bookPtr
	}
	return result
}
