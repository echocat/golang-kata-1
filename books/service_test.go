package books

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestNewBookService(t *testing.T) {
	service := NewBookService()
	if service == nil {
		t.Fatal("NewBookService should return a non-nil service")
	}

	books := service.GetAll()
	if len(books) != 0 {
		t.Errorf("Expected empty books slice, got %d books", len(books))
	}
}

func TestBookService_LoadFromCSV(t *testing.T) {
	service := NewBookService()

	err := service.LoadFromCSV("../resources/books.csv")
	if err != nil {
		t.Fatalf("Failed to load books from CSV: %v", err)
	}

	books := service.GetAll()
	if len(books) == 0 {
		t.Fatal("Expected to load books, but got empty slice")
	}

	// Check if we loaded the expected number of books from the test data
	expectedBookCount := 8
	if len(books) != expectedBookCount {
		t.Errorf("Expected %d books, got %d", expectedBookCount, len(books))
	}

	// Test first book from CSV
	expectedFirstBook := Book{
		Title:       "Ich helfe dir kochen. Das erfolgreiche Universalkochbuch mit großem Backteil",
		ISBN:        "5554-5545-4518",
		Authors:     []string{"null-walter@echocat.org"},
		Description: "Auf der Suche nach einem Basiskochbuch steht man heutzutage vor einer Fülle von Alternativen. Es fällt schwer, daraus die für sich passende Mixtur aus Grundlagenwerk und Rezeptesammlung zu finden. Man sollte sich darüber im Klaren sein, welchen Schwerpunkt man setzen möchte oder von welchen Koch- und Backkenntnissen man bereits ausgehen kann.",
	}

	if !reflect.DeepEqual(books[0], expectedFirstBook) {
		t.Errorf("First book doesn't match expected. Got: %+v", books[0])
	}

	// Test book with multiple authors
	var multiAuthorBook *Book
	for _, book := range books {
		if book.ISBN == "2145-8548-3325" {
			multiAuthorBook = &book
			break
		}
	}

	if multiAuthorBook == nil {
		t.Error("Could not find book with multiple authors")
	} else {
		expectedAuthors := []string{"null-ferdinand@echocat.org", "null-lieblich@echocat.org"}
		if !reflect.DeepEqual(multiAuthorBook.Authors, expectedAuthors) {
			t.Errorf("Expected authors %v, got %v", expectedAuthors, multiAuthorBook.Authors)
		}
	}
}

func TestBookService_LoadFromCSV_ErrorCases(t *testing.T) {
	tests := []struct {
		name        string
		setupFile   func() (string, func())
		expectError bool
		errorSubstr string
	}{
		{
			name: "non-existent file",
			setupFile: func() (string, func()) {
				return "non_existent_file.csv", func() {}
			},
			expectError: true,
			errorSubstr: "failed to open CSV file",
		},
		{
			name: "invalid CSV format",
			setupFile: func() (string, func()) {
				tempFile, err := os.CreateTemp("", "invalid_books_*.csv")
				if err != nil {
					t.Fatalf("Failed to create temp file: %v", err)
				}
				invalidCSV := `title;isbn;authors;description
"Unclosed quote book;1234-5678-9012;author@example.com;Description`
				if _, err := tempFile.WriteString(invalidCSV); err != nil {
					t.Fatalf("Failed to write to temp file: %v", err)
				}
				if err := tempFile.Close(); err != nil {
					t.Fatalf("Failed to close temp file: %v", err)
				}
				return tempFile.Name(), func() {
					if err := os.Remove(tempFile.Name()); err != nil {
						t.Errorf("Failed to remove temp file: %v", err)
					}
				}
			},
			expectError: true,
			errorSubstr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath, cleanup := tt.setupFile()
			defer cleanup()

			service := NewBookService()
			err := service.LoadFromCSV(filePath)

			if tt.expectError && err == nil {
				t.Errorf("Expected error for %s, but got nil", tt.name)
			}

			if !tt.expectError && err != nil {
				t.Errorf("Expected no error for %s, but got: %v", tt.name, err)
			}

			if tt.errorSubstr != "" && !strings.Contains(err.Error(), tt.errorSubstr) {
				t.Errorf("Expected error message to contain '%s', got: %v", tt.errorSubstr, err)
			}
		})
	}
}

func TestBookService_FindByISBN(t *testing.T) {
	service := NewBookService()
	err := service.LoadFromCSV("../resources/books.csv")
	if err != nil {
		t.Fatalf("Failed to load books: %v", err)
	}

	tests := []struct {
		name          string
		isbn          string
		expectFound   bool
		expectedTitle string
	}{
		{
			name:          "existing book",
			isbn:          "5554-5545-4518",
			expectFound:   true,
			expectedTitle: "Ich helfe dir kochen. Das erfolgreiche Universalkochbuch mit großem Backteil",
		},
		{
			name:          "non-existent book",
			isbn:          "0000-0000-0000",
			expectFound:   false,
			expectedTitle: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			book := service.FindByISBN(tt.isbn)

			if tt.expectFound && book == nil {
				t.Errorf("Expected to find book with ISBN %s, but got nil", tt.isbn)
				return
			}

			if !tt.expectFound && book != nil {
				t.Errorf("Expected to get nil for ISBN %s, but got %+v", tt.isbn, book)
				return
			}

			if tt.expectFound {
				if book.ISBN != tt.isbn {
					t.Errorf("Expected book ISBN to be %s, got %s", tt.isbn, book.ISBN)
				}
				if book.Title != tt.expectedTitle {
					t.Errorf("Expected book title to be %s, got %s", tt.expectedTitle, book.Title)
				}
			}
		})
	}
}

func TestBookService_FindByAuthorEmail(t *testing.T) {
	service := NewBookService()
	err := service.LoadFromCSV("../resources/books.csv")
	if err != nil {
		t.Fatalf("Failed to load books: %v", err)
	}

	tests := []struct {
		name             string
		email            string
		expectedMinCount int
		expectEmpty      bool
	}{
		{
			name:             "existing author with books",
			email:            "null-walter@echocat.org",
			expectedMinCount: 1,
			expectEmpty:      false,
		},
		{
			name:             "author with multiple books",
			email:            "null-lieblich@echocat.org",
			expectedMinCount: 2,
			expectEmpty:      false,
		},
		{
			name:             "non-existent author",
			email:            "non-existent@example.com",
			expectedMinCount: 0,
			expectEmpty:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			books := service.FindByAuthorEmail(tt.email)

			if tt.expectEmpty && len(books) != 0 {
				t.Errorf("Expected empty slice for %s, but got %d books", tt.email, len(books))
				return
			}

			if !tt.expectEmpty && len(books) == 0 {
				t.Errorf("Expected to find books for author %s, but got empty slice", tt.email)
				return
			}

			if len(books) < tt.expectedMinCount {
				t.Errorf("Expected at least %d books for author %s, got %d", tt.expectedMinCount, tt.email, len(books))
			}

			// Verify that all returned books contain the author
			for _, book := range books {
				found := false
				for _, author := range book.Authors {
					if author == tt.email {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Book %s does not contain author %s", book.Title, tt.email)
				}
			}
		})
	}
}

func TestBookService_GetAll(t *testing.T) {
	service := NewBookService()

	// Test empty service
	books := service.GetAll()
	if len(books) != 0 {
		t.Errorf("Expected empty slice for new service, got %d books", len(books))
	}

	// Load books and test again
	err := service.LoadFromCSV("../resources/books.csv")
	if err != nil {
		t.Fatalf("Failed to load books: %v", err)
	}

	books = service.GetAll()
	if len(books) == 0 {
		t.Error("Expected non-empty slice after loading books")
	}

	// Verify that modifying the returned slice doesn't affect the service
	originalLength := len(books)
	_ = append(books, Book{Title: "Test Book"})

	booksAfterModification := service.GetAll()
	if len(booksAfterModification) != originalLength {
		t.Error("Modifying returned slice should not affect the service's internal state")
	}
}

// Benchmark tests
func BenchmarkBookService_FindByISBN(b *testing.B) {
	service := NewBookService()
	err := service.LoadFromCSV("../resources/books.csv")
	if err != nil {
		b.Fatalf("Failed to load books: %v", err)
	}

	isbn := "5554-5545-4518"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.FindByISBN(isbn)
	}
}

func BenchmarkBookService_FindByAuthorEmail(b *testing.B) {
	service := NewBookService()
	err := service.LoadFromCSV("../resources/books.csv")
	if err != nil {
		b.Fatalf("Failed to load books: %v", err)
	}

	email := "null-walter@echocat.org"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.FindByAuthorEmail(email)
	}
}

func BenchmarkBookService_LoadFromCSV(b *testing.B) {
	for i := 0; i < b.N; i++ {
		service := NewBookService()
		err := service.LoadFromCSV("../resources/books.csv")
		if err != nil {
			b.Fatalf("Failed to load books: %v", err)
		}
	}
}
