package authors

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type Service interface {
	LoadFromCSV(filePath string) error
	GetAll() []Author
	FindByEmail(email string) *Author
}

type authorService struct {
	authors    []Author
	emailIndex map[string]*Author
}

func NewAuthorService() Service {
	return &authorService{
		authors:    make([]Author, 0),
		emailIndex: make(map[string]*Author),
	}
}

func (s *authorService) LoadFromCSV(filePath string) error {
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

		if len(record) < 3 {
			continue
		}

		author := Author{
			Email:     strings.TrimSpace(record[0]),
			FirstName: strings.TrimSpace(record[1]),
			LastName:  strings.TrimSpace(record[2]),
		}

		s.authors = append(s.authors, author)

		authorPtr := &s.authors[len(s.authors)-1]
		s.emailIndex[author.Email] = authorPtr
	}

	return nil
}

func (s *authorService) GetAll() []Author {
	return s.authors
}

func (s *authorService) FindByEmail(email string) *Author {
	return s.emailIndex[email]
}
