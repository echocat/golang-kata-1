package publications

import (
	"reflect"
	"testing"
	"time"

	"github.com/echocat/golang-kata-1/authors"
	"github.com/echocat/golang-kata-1/books"
	"github.com/echocat/golang-kata-1/magazines"
)

// Mock services for testing
type mockBookService struct {
	books []books.Book
}

func (m *mockBookService) LoadFromCSV(filePath string) error {
	return nil
}

func (m *mockBookService) GetAll() []books.Book {
	return m.books
}

func (m *mockBookService) FindByISBN(isbn string) *books.Book {
	for i := range m.books {
		if m.books[i].ISBN == isbn {
			return &m.books[i]
		}
	}
	return nil
}

func (m *mockBookService) FindByAuthorEmail(email string) []books.Book {
	var result []books.Book
	for _, book := range m.books {
		for _, author := range book.Authors {
			if author == email {
				result = append(result, book)
				break
			}
		}
	}
	return result
}

type mockMagazineService struct {
	magazines []magazines.Magazine
}

func (m *mockMagazineService) LoadFromCSV(filePath string) error {
	return nil
}

func (m *mockMagazineService) GetAll() []magazines.Magazine {
	return m.magazines
}

func (m *mockMagazineService) FindByISBN(isbn string) *magazines.Magazine {
	for i := range m.magazines {
		if m.magazines[i].ISBN == isbn {
			return &m.magazines[i]
		}
	}
	return nil
}

func (m *mockMagazineService) FindByAuthorEmail(email string) []magazines.Magazine {
	var result []magazines.Magazine
	for _, magazine := range m.magazines {
		for _, author := range magazine.Authors {
			if author == email {
				result = append(result, magazine)
				break
			}
		}
	}
	return result
}

type mockAuthorService struct {
	authors []authors.Author
}

func (m *mockAuthorService) LoadFromCSV(filePath string) error {
	return nil
}

func (m *mockAuthorService) GetAll() []authors.Author {
	return m.authors
}

func (m *mockAuthorService) FindByEmail(email string) *authors.Author {
	for i := range m.authors {
		if m.authors[i].Email == email {
			return &m.authors[i]
		}
	}
	return nil
}

func createTestServices() (books.Service, magazines.Service, authors.Service) {
	publishedAt, _ := time.Parse("02.01.2006", "01.01.2020")

	bookService := &mockBookService{
		books: []books.Book{
			{
				Title:       "Test Book 1",
				ISBN:        "1111-1111-1111",
				Authors:     []string{"author1@test.com"},
				Description: "Test description 1",
			},
			{
				Title:       "Test Book 2",
				ISBN:        "2222-2222-2222",
				Authors:     []string{"author2@test.com", "author1@test.com"},
				Description: "Test description 2",
			},
		},
	}

	magazineService := &mockMagazineService{
		magazines: []magazines.Magazine{
			{
				Title:       "Test Magazine 1",
				ISBN:        "3333-3333-3333",
				Authors:     []string{"author1@test.com"},
				PublishedAt: publishedAt,
			},
			{
				Title:       "Test Magazine 2",
				ISBN:        "4444-4444-4444",
				Authors:     []string{"author3@test.com"},
				PublishedAt: publishedAt,
			},
		},
	}

	authorService := &mockAuthorService{
		authors: []authors.Author{
			{Email: "author1@test.com", FirstName: "John", LastName: "Doe"},
			{Email: "author2@test.com", FirstName: "Jane", LastName: "Smith"},
			{Email: "author3@test.com", FirstName: "Bob", LastName: "Johnson"},
		},
	}

	return bookService, magazineService, authorService
}

func TestCollectPublications(t *testing.T) {
	tests := []struct {
		name              string
		setupServices     func() (books.Service, magazines.Service, authors.Service)
		expectedTotal     int
		expectedBookCount int
		expectedMagCount  int
	}{
		{
			name:              "normal services with data",
			setupServices:     createTestServices,
			expectedTotal:     4,
			expectedBookCount: 2,
			expectedMagCount:  2,
		},
		{
			name: "empty services",
			setupServices: func() (books.Service, magazines.Service, authors.Service) {
				bookService := &mockBookService{books: []books.Book{}}
				magazineService := &mockMagazineService{magazines: []magazines.Magazine{}}
				authorService := &mockAuthorService{authors: []authors.Author{}}
				return bookService, magazineService, authorService
			},
			expectedTotal:     0,
			expectedBookCount: 0,
			expectedMagCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bookService, magazineService, authorService := tt.setupServices()
			publications := CollectPublications(bookService, magazineService, authorService)

			if len(publications) != tt.expectedTotal {
				t.Errorf("Expected %d publications, got %d", tt.expectedTotal, len(publications))
			}

			// Check that we have the expected counts of books and magazines
			bookCount := 0
			magazineCount := 0
			for _, pub := range publications {
				switch pub.Type {
				case "Book":
					bookCount++
				case "Magazine":
					magazineCount++
				}
			}

			if bookCount != tt.expectedBookCount {
				t.Errorf("Expected %d books, got %d", tt.expectedBookCount, bookCount)
			}

			if magazineCount != tt.expectedMagCount {
				t.Errorf("Expected %d magazines, got %d", tt.expectedMagCount, magazineCount)
			}

			// For non-empty tests, verify specific publication details
			if tt.expectedTotal > 0 {
				// Test first book publication
				var firstBook *Publication
				for _, pub := range publications {
					if pub.Type == "Book" && pub.ISBN == "1111-1111-1111" {
						firstBook = &pub
						break
					}
				}

				if firstBook == nil {
					t.Error("Could not find first book publication")
				} else {
					expectedBook := Publication{
						Title:       "Test Book 1",
						ISBN:        "1111-1111-1111",
						Authors:     []string{"author1@test.com"},
						Type:        "Book",
						Description: "Test description 1",
						PublishedAt: nil,
					}

					if !reflect.DeepEqual(*firstBook, expectedBook) {
						t.Errorf("First book publication doesn't match expected. Got: %+v, Expected: %+v", *firstBook, expectedBook)
					}
				}

				// Test first magazine publication
				var firstMagazine *Publication
				for _, pub := range publications {
					if pub.Type == "Magazine" && pub.ISBN == "3333-3333-3333" {
						firstMagazine = &pub
						break
					}
				}

				if firstMagazine == nil {
					t.Error("Could not find first magazine publication")
				} else {
					if firstMagazine.PublishedAt == nil {
						t.Error("Magazine publication should have a PublishedAt date")
					}

					if firstMagazine.Description != "" {
						t.Error("Magazine publication should have empty description")
					}
				}
			}
		})
	}
}

func TestFilterByISBN(t *testing.T) {
	publications := []Publication{
		{Title: "Test 1", ISBN: "1111-1111-1111", Type: "Book"},
		{Title: "Test 2", ISBN: "2222-2222-2222", Type: "Book"},
		{Title: "Test 3", ISBN: "3333-3333-3333", Type: "Magazine"},
		{Title: "Test 4", ISBN: "1111-1111-1111", Type: "Magazine"}, // Duplicate ISBN
	}

	tests := []struct {
		name          string
		publications  []Publication
		isbn          string
		expectedCount int
	}{
		{
			name:          "existing ISBN with duplicates",
			publications:  publications,
			isbn:          "1111-1111-1111",
			expectedCount: 2,
		},
		{
			name:          "existing ISBN single match",
			publications:  publications,
			isbn:          "2222-2222-2222",
			expectedCount: 1,
		},
		{
			name:          "non-existent ISBN",
			publications:  publications,
			isbn:          "9999-9999-9999",
			expectedCount: 0,
		},
		{
			name:          "empty publications slice",
			publications:  []Publication{},
			isbn:          "1111-1111-1111",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := FilterByISBN(tt.publications, tt.isbn)

			if len(filtered) != tt.expectedCount {
				t.Errorf("Expected %d publications with ISBN %s, got %d", tt.expectedCount, tt.isbn, len(filtered))
			}

			// Verify all filtered publications have the correct ISBN
			for _, pub := range filtered {
				if pub.ISBN != tt.isbn {
					t.Errorf("Filtered publication should have ISBN %s, got %s", tt.isbn, pub.ISBN)
				}
			}
		})
	}
}

func TestFilterByAuthorEmail(t *testing.T) {
	publications := []Publication{
		{Title: "Test 1", Authors: []string{"author1@test.com"}, Type: "Book"},
		{Title: "Test 2", Authors: []string{"author2@test.com", "author1@test.com"}, Type: "Book"},
		{Title: "Test 3", Authors: []string{"author3@test.com"}, Type: "Magazine"},
		{Title: "Test 4", Authors: []string{"author1@test.com", "author3@test.com"}, Type: "Magazine"},
	}

	tests := []struct {
		name          string
		publications  []Publication
		email         string
		expectedCount int
	}{
		{
			name:          "author with multiple publications",
			publications:  publications,
			email:         "author1@test.com",
			expectedCount: 3,
		},
		{
			name:          "author with single publication",
			publications:  publications,
			email:         "author2@test.com",
			expectedCount: 1,
		},
		{
			name:          "non-existent author",
			publications:  publications,
			email:         "nonexistent@test.com",
			expectedCount: 0,
		},
		{
			name:          "empty publications slice",
			publications:  []Publication{},
			email:         "author1@test.com",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := FilterByAuthorEmail(tt.publications, tt.email)

			if len(filtered) != tt.expectedCount {
				t.Errorf("Expected %d publications for %s, got %d", tt.expectedCount, tt.email, len(filtered))
			}

			// Verify all filtered publications contain the author
			for _, pub := range filtered {
				found := false
				for _, author := range pub.Authors {
					if author == tt.email {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Filtered publication should contain author %s, but doesn't: %+v", tt.email, pub)
				}
			}
		})
	}
}

func TestSortByTitle(t *testing.T) {
	tests := []struct {
		name            string
		publications    []Publication
		expectedTitles  []string
		shouldNotModify bool
	}{
		{
			name: "normal sorting",
			publications: []Publication{
				{Title: "Zebra Book", ISBN: "1111"},
				{Title: "Apple Magazine", ISBN: "2222"},
				{Title: "Beta Publication", ISBN: "3333"},
				{Title: "Alpha Article", ISBN: "4444"},
			},
			expectedTitles:  []string{"Alpha Article", "Apple Magazine", "Beta Publication", "Zebra Book"},
			shouldNotModify: true,
		},
		{
			name:            "empty slice",
			publications:    []Publication{},
			expectedTitles:  []string{},
			shouldNotModify: false,
		},
		{
			name: "single publication",
			publications: []Publication{
				{Title: "Single Book", ISBN: "1111"},
			},
			expectedTitles:  []string{"Single Book"},
			shouldNotModify: false,
		},
		{
			name: "case-sensitive sorting",
			publications: []Publication{
				{Title: "zebra book", ISBN: "1111"},
				{Title: "Apple Magazine", ISBN: "2222"},
				{Title: "BETA publication", ISBN: "3333"},
			},
			expectedTitles:  []string{"Apple Magazine", "BETA publication", "zebra book"},
			shouldNotModify: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalFirstTitle := ""
			if len(tt.publications) > 0 {
				originalFirstTitle = tt.publications[0].Title
			}

			sorted := SortByTitle(tt.publications)

			if len(sorted) != len(tt.expectedTitles) {
				t.Errorf("Expected %d sorted publications, got %d", len(tt.expectedTitles), len(sorted))
			}

			for i, expectedTitle := range tt.expectedTitles {
				if i < len(sorted) && sorted[i].Title != expectedTitle {
					t.Errorf("Expected title at position %d to be '%s', got '%s'", i, expectedTitle, sorted[i].Title)
				}
			}

			// Verify original slice is not modified if expected
			if tt.shouldNotModify && len(tt.publications) > 0 && tt.publications[0].Title != originalFirstTitle {
				t.Error("Original publications slice should not be modified by SortByTitle")
			}
		})
	}
}

// Integration test using real services
func TestCollectPublications_Integration(t *testing.T) {
	// This test requires the actual CSV files to be present
	bookService := books.NewBookService()
	magazineService := magazines.NewMagazineService()
	authorService := authors.NewAuthorService()

	// Load from actual CSV files
	err := bookService.LoadFromCSV("../resources/books.csv")
	if err != nil {
		t.Skipf("Skipping integration test - could not load books.csv: %v", err)
	}

	err = magazineService.LoadFromCSV("../resources/magazines.csv")
	if err != nil {
		t.Skipf("Skipping integration test - could not load magazines.csv: %v", err)
	}

	err = authorService.LoadFromCSV("../resources/authors.csv")
	if err != nil {
		t.Skipf("Skipping integration test - could not load authors.csv: %v", err)
	}

	publications := CollectPublications(bookService, magazineService, authorService)

	if len(publications) == 0 {
		t.Error("Expected to collect publications from CSV files, but got empty slice")
	}

	// Should have both books and magazines
	hasBooks := false
	hasMagazines := false
	for _, pub := range publications {
		switch pub.Type {
		case "Book":
			hasBooks = true
		case "Magazine":
			hasMagazines = true
		}
	}

	if !hasBooks {
		t.Error("Expected to find books in publications")
	}
	if !hasMagazines {
		t.Error("Expected to find magazines in publications")
	}
}

// Benchmark tests
func BenchmarkCollectPublications(b *testing.B) {
	bookService, magazineService, authorService := createTestServices()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CollectPublications(bookService, magazineService, authorService)
	}
}

func BenchmarkFilterByISBN(b *testing.B) {
	publications := make([]Publication, 1000)
	for i := 0; i < 1000; i++ {
		publications[i] = Publication{
			Title: "Test Publication",
			ISBN:  "1111-1111-1111",
			Type:  "Book",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FilterByISBN(publications, "1111-1111-1111")
	}
}

func BenchmarkFilterByAuthorEmail(b *testing.B) {
	publications := make([]Publication, 1000)
	for i := 0; i < 1000; i++ {
		publications[i] = Publication{
			Title:   "Test Publication",
			Authors: []string{"author1@test.com", "author2@test.com"},
			Type:    "Book",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FilterByAuthorEmail(publications, "author1@test.com")
	}
}

func BenchmarkSortByTitle(b *testing.B) {
	publications := make([]Publication, 1000)
	for i := 0; i < 1000; i++ {
		publications[i] = Publication{
			Title: "Test Publication",
			Type:  "Book",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SortByTitle(publications)
	}
}
