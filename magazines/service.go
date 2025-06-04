package magazines

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

type Service interface {
	LoadFromCSV(filePath string) error
	GetAll() []Magazine
	FindByISBN(isbn string) *Magazine
	FindByAuthorEmail(email string) []Magazine
}

type magazineService struct {
	magazines   []Magazine
	isbnIndex   map[string]*Magazine
	authorIndex map[string][]*Magazine
}

func NewMagazineService() Service {
	return &magazineService{
		magazines:   make([]Magazine, 0),
		isbnIndex:   make(map[string]*Magazine),
		authorIndex: make(map[string][]*Magazine),
	}
}

func (s *magazineService) LoadFromCSV(filePath string) error {
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

		publishedAtStr := strings.TrimSpace(record[3])
		publishedAt, err := time.Parse("02.01.2006", publishedAtStr)
		if err != nil {
			return fmt.Errorf("failed to parse publication date '%s': %w", publishedAtStr, err)
		}

		magazine := Magazine{
			Title:       strings.TrimSpace(record[0]),
			ISBN:        strings.TrimSpace(record[1]),
			Authors:     authors,
			PublishedAt: publishedAt,
		}

		s.magazines = append(s.magazines, magazine)

		magazinePtr := &s.magazines[len(s.magazines)-1]

		s.isbnIndex[magazine.ISBN] = magazinePtr

		for _, author := range magazine.Authors {
			s.authorIndex[author] = append(s.authorIndex[author], magazinePtr)
		}
	}

	return nil
}

func (s *magazineService) GetAll() []Magazine {
	return s.magazines
}

func (s *magazineService) FindByISBN(isbn string) *Magazine {
	return s.isbnIndex[isbn]
}

func (s *magazineService) FindByAuthorEmail(email string) []Magazine {
	magazinePtrs := s.authorIndex[email]
	result := make([]Magazine, len(magazinePtrs))
	for i, magazinePtr := range magazinePtrs {
		result[i] = *magazinePtr
	}
	return result
}
