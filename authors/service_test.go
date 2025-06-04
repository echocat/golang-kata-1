package authors

import (
	"os"
	"strings"
	"testing"
)

func TestNewAuthorService(t *testing.T) {
	service := NewAuthorService()
	if service == nil {
		t.Fatal("NewAuthorService should return a non-nil service")
	}

	authors := service.GetAll()
	if len(authors) != 0 {
		t.Errorf("Expected empty authors slice, got %d authors", len(authors))
	}
}

func TestAuthorService_LoadFromCSV(t *testing.T) {
	service := NewAuthorService()

	err := service.LoadFromCSV("../resources/authors.csv")
	if err != nil {
		t.Fatalf("Failed to load authors from CSV: %v", err)
	}

	authors := service.GetAll()
	if len(authors) == 0 {
		t.Fatal("Expected to load authors, but got empty slice")
	}

	// Check if we loaded the expected number of authors from the test data
	expectedAuthorCount := 6
	if len(authors) != expectedAuthorCount {
		t.Errorf("Expected %d authors, got %d", expectedAuthorCount, len(authors))
	}

	// Test first author from CSV
	expectedFirstAuthor := Author{
		Email:     "null-walter@echocat.org",
		FirstName: "Paul",
		LastName:  "Walter",
	}

	if authors[0] != expectedFirstAuthor {
		t.Errorf("First author doesn't match expected. Got: %+v, Expected: %+v", authors[0], expectedFirstAuthor)
	}
}

func TestAuthorService_LoadFromCSV_ErrorCases(t *testing.T) {
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
				tempFile, err := os.CreateTemp("", "invalid_authors_*.csv")
				if err != nil {
					t.Fatalf("Failed to create temp file: %v", err)
				}
				invalidCSV := `email;firstname;lastname
"unclosed-quote@example.com;John;Doe`
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
		{
			name: "incomplete record (less than 3 fields)",
			setupFile: func() (string, func()) {
				tempFile, err := os.CreateTemp("", "incomplete_authors_*.csv")
				if err != nil {
					t.Fatalf("Failed to create temp file: %v", err)
				}
				incompleteCSV := `email;firstname;lastname
test@example.com;John
complete@example.com;Jane;Doe`
				if _, err := tempFile.WriteString(incompleteCSV); err != nil {
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
			expectError: false, // Should not error, just skip incomplete records
			errorSubstr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath, cleanup := tt.setupFile()
			defer cleanup()

			service := NewAuthorService()
			err := service.LoadFromCSV(filePath)

			if tt.expectError && err == nil {
				t.Errorf("Expected error for %s, but got nil", tt.name)
			}

			if !tt.expectError && err != nil {
				t.Errorf("Expected no error for %s, but got: %v", tt.name, err)
			}

			if tt.errorSubstr != "" && err != nil && !strings.Contains(err.Error(), tt.errorSubstr) {
				t.Errorf("Expected error message to contain '%s', got: %v", tt.errorSubstr, err)
			}

			// For incomplete record test, verify that valid records are still loaded
			if tt.name == "incomplete record (less than 3 fields)" && err == nil {
				authors := service.GetAll()
				if len(authors) != 1 {
					t.Errorf("Expected 1 valid author to be loaded, got %d", len(authors))
				}
				if len(authors) > 0 && authors[0].Email != "complete@example.com" {
					t.Errorf("Expected valid author email to be 'complete@example.com', got '%s'", authors[0].Email)
				}
			}
		})
	}
}

func TestAuthorService_FindByEmail(t *testing.T) {
	service := NewAuthorService()
	err := service.LoadFromCSV("../resources/authors.csv")
	if err != nil {
		t.Fatalf("Failed to load authors: %v", err)
	}

	tests := []struct {
		name              string
		email             string
		expectFound       bool
		expectedFirstName string
		expectedLastName  string
	}{
		{
			name:              "existing author - Paul Walter",
			email:             "null-walter@echocat.org",
			expectFound:       true,
			expectedFirstName: "Paul",
			expectedLastName:  "Walter",
		},
		{
			name:              "existing author - Max Müller",
			email:             "null-mueller@echocat.org",
			expectFound:       true,
			expectedFirstName: "Max",
			expectedLastName:  "Müller",
		},
		{
			name:              "existing author - Franz Ferdinand",
			email:             "null-ferdinand@echocat.org",
			expectFound:       true,
			expectedFirstName: "Franz",
			expectedLastName:  "Ferdinand",
		},
		{
			name:              "non-existent author",
			email:             "non-existent@example.com",
			expectFound:       false,
			expectedFirstName: "",
			expectedLastName:  "",
		},
		{
			name:              "empty email",
			email:             "",
			expectFound:       false,
			expectedFirstName: "",
			expectedLastName:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			author := service.FindByEmail(tt.email)

			if tt.expectFound && author == nil {
				t.Errorf("Expected to find author with email %s, but got nil", tt.email)
				return
			}

			if !tt.expectFound && author != nil {
				t.Errorf("Expected to get nil for email %s, but got %+v", tt.email, author)
				return
			}

			if tt.expectFound {
				if author.Email != tt.email {
					t.Errorf("Expected author email to be %s, got %s", tt.email, author.Email)
				}
				if author.FirstName != tt.expectedFirstName {
					t.Errorf("Expected author first name to be %s, got %s", tt.expectedFirstName, author.FirstName)
				}
				if author.LastName != tt.expectedLastName {
					t.Errorf("Expected author last name to be %s, got %s", tt.expectedLastName, author.LastName)
				}
			}
		})
	}
}

func TestAuthorService_GetAll(t *testing.T) {
	tests := []struct {
		name          string
		setupService  func() Service
		expectedCount int
		shouldLoadCSV bool
	}{
		{
			name: "empty service",
			setupService: func() Service {
				return NewAuthorService()
			},
			expectedCount: 0,
			shouldLoadCSV: false,
		},
		{
			name: "service with loaded authors",
			setupService: func() Service {
				service := NewAuthorService()
				err := service.LoadFromCSV("../resources/authors.csv")
				if err != nil {
					t.Fatalf("Failed to load authors: %v", err)
				}
				return service
			},
			expectedCount: 6,
			shouldLoadCSV: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := tt.setupService()
			authors := service.GetAll()

			if len(authors) != tt.expectedCount {
				t.Errorf("Expected %d authors, got %d", tt.expectedCount, len(authors))
			}

			// Verify that modifying the returned slice doesn't affect the service's internal state
			if len(authors) > 0 {
				originalLength := len(authors)
				_ = append(authors, Author{Email: "test@example.com", FirstName: "Test", LastName: "User"})

				authorsAfterModification := service.GetAll()
				if len(authorsAfterModification) != originalLength {
					t.Error("Modifying returned slice should not affect the service's internal state")
				}
			}
		})
	}
}

func TestAuthorService_DataIntegrity(t *testing.T) {
	// Test that data is properly trimmed and handled
	tempFile, err := os.CreateTemp("", "data_integrity_authors_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() {
		if err := os.Remove(tempFile.Name()); err != nil {
			t.Fatalf("Failed to remove temp file: %v", err)
		}
	}()

	// CSV with whitespace that should be trimmed
	csvWithWhitespace := `email;firstname;lastname
  spaced@example.com  ;  John  ;  Doe
normal@example.com;Jane;Smith`

	if _, err := tempFile.WriteString(csvWithWhitespace); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	service := NewAuthorService()
	err = service.LoadFromCSV(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to load authors: %v", err)
	}

	authors := service.GetAll()
	if len(authors) != 2 {
		t.Errorf("Expected 2 authors, got %d", len(authors))
	}

	// Test that whitespace is properly trimmed
	spacedAuthor := service.FindByEmail("spaced@example.com")
	if spacedAuthor == nil {
		t.Error("Expected to find author with trimmed email")
	} else {
		if spacedAuthor.FirstName != "John" {
			t.Errorf("Expected first name to be 'John' (trimmed), got '%s'", spacedAuthor.FirstName)
		}
		if spacedAuthor.LastName != "Doe" {
			t.Errorf("Expected last name to be 'Doe' (trimmed), got '%s'", spacedAuthor.LastName)
		}
	}
}

func TestAuthorService_EmailIndex(t *testing.T) {
	// Test that the email index is properly maintained
	service := NewAuthorService()
	err := service.LoadFromCSV("../resources/authors.csv")
	if err != nil {
		t.Fatalf("Failed to load authors: %v", err)
	}

	authors := service.GetAll()

	// Verify that each author can be found by their email
	for _, expectedAuthor := range authors {
		foundAuthor := service.FindByEmail(expectedAuthor.Email)
		if foundAuthor == nil {
			t.Errorf("Expected to find author with email %s", expectedAuthor.Email)
		} else if *foundAuthor != expectedAuthor {
			t.Errorf("Found author doesn't match expected. Got: %+v, Expected: %+v", *foundAuthor, expectedAuthor)
		}
	}
}

func BenchmarkAuthorService_FindByEmail(b *testing.B) {
	service := NewAuthorService()

	err := service.LoadFromCSV("../resources/authors.csv")
	if err != nil {
		b.Fatalf("Failed to load authors: %v", err)
	}

	email := "null-walter@echocat.org"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.FindByEmail(email)
	}
}

func BenchmarkAuthorService_LoadFromCSV(b *testing.B) {
	for i := 0; i < b.N; i++ {
		service := NewAuthorService()
		err := service.LoadFromCSV("../resources/authors.csv")
		if err != nil {
			b.Fatalf("Failed to load authors: %v", err)
		}
	}
}

func BenchmarkAuthorService_GetAll(b *testing.B) {
	service := NewAuthorService()
	err := service.LoadFromCSV("../resources/authors.csv")
	if err != nil {
		b.Fatalf("Failed to load authors: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.GetAll()
	}
}
