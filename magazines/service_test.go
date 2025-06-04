package magazines

import (
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestNewMagazineService(t *testing.T) {
	service := NewMagazineService()
	if service == nil {
		t.Fatal("NewMagazineService should return a non-nil service")
	}

	magazines := service.GetAll()
	if len(magazines) != 0 {
		t.Errorf("Expected empty magazines slice, got %d magazines", len(magazines))
	}
}

func TestMagazineService_LoadFromCSV(t *testing.T) {
	service := NewMagazineService()

	err := service.LoadFromCSV("../resources/magazines.csv")
	if err != nil {
		t.Fatalf("Failed to load magazines from CSV: %v", err)
	}

	magazines := service.GetAll()
	if len(magazines) == 0 {
		t.Fatal("Expected to load magazines, but got empty slice")
	}

	// Check if we loaded the expected number of magazines from the test data
	expectedMagazineCount := 6
	if len(magazines) != expectedMagazineCount {
		t.Errorf("Expected %d magazines, got %d", expectedMagazineCount, len(magazines))
	}

	// Test first magazine from CSV
	expectedDate, _ := time.Parse("02.01.2006", "21.05.2011")
	expectedFirstMagazine := Magazine{
		Title:       "Beautiful cooking",
		ISBN:        "5454-5587-3210",
		Authors:     []string{"null-walter@echocat.org"},
		PublishedAt: expectedDate,
	}

	if !reflect.DeepEqual(magazines[0], expectedFirstMagazine) {
		t.Errorf("First magazine doesn't match expected. Got: %+v, Expected: %+v", magazines[0], expectedFirstMagazine)
	}

	// Test magazine with multiple authors
	var multiAuthorMagazine *Magazine
	for _, magazine := range magazines {
		if magazine.ISBN == "2365-5632-7854" {
			multiAuthorMagazine = &magazine
			break
		}
	}

	if multiAuthorMagazine == nil {
		t.Error("Could not find magazine with multiple authors")
	} else {
		expectedAuthors := []string{"null-lieblich@echocat.org", "null-walter@echocat.org"}
		if !reflect.DeepEqual(multiAuthorMagazine.Authors, expectedAuthors) {
			t.Errorf("Expected authors %v, got %v", expectedAuthors, multiAuthorMagazine.Authors)
		}
	}
}

func TestMagazineService_LoadFromCSV_ErrorCases(t *testing.T) {
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
			name: "invalid date format",
			setupFile: func() (string, func()) {
				tempFile, err := os.CreateTemp("", "invalid_magazines_*.csv")
				if err != nil {
					t.Fatalf("Failed to create temp file: %v", err)
				}
				invalidDateCSV := `title;isbn;authors;publishedAt
Test Magazine;1234-5678-9012;author@example.com;invalid-date`
				if _, err := tempFile.WriteString(invalidDateCSV); err != nil {
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
			errorSubstr: "failed to parse publication date",
		},
		{
			name: "invalid CSV format",
			setupFile: func() (string, func()) {
				tempFile, err := os.CreateTemp("", "invalid_magazines_*.csv")
				if err != nil {
					t.Fatalf("Failed to create temp file: %v", err)
				}
				invalidCSV := `title;isbn;authors;publishedAt
"Unclosed quote magazine;1234-5678-9012;author@example.com;01.01.2020`
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

			service := NewMagazineService()
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

func TestMagazineService_FindByISBN(t *testing.T) {
	service := NewMagazineService()
	err := service.LoadFromCSV("../resources/magazines.csv")
	if err != nil {
		t.Fatalf("Failed to load magazines: %v", err)
	}

	tests := []struct {
		name            string
		isbn            string
		expectFound     bool
		expectedTitle   string
		expectedDateStr string
	}{
		{
			name:            "existing magazine",
			isbn:            "5454-5587-3210",
			expectFound:     true,
			expectedTitle:   "Beautiful cooking",
			expectedDateStr: "21.05.2011",
		},
		{
			name:            "non-existent magazine",
			isbn:            "0000-0000-0000",
			expectFound:     false,
			expectedTitle:   "",
			expectedDateStr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			magazine := service.FindByISBN(tt.isbn)

			if tt.expectFound && magazine == nil {
				t.Errorf("Expected to find magazine with ISBN %s, but got nil", tt.isbn)
				return
			}

			if !tt.expectFound && magazine != nil {
				t.Errorf("Expected to get nil for ISBN %s, but got %+v", tt.isbn, magazine)
				return
			}

			if tt.expectFound {
				if magazine.ISBN != tt.isbn {
					t.Errorf("Expected magazine ISBN to be %s, got %s", tt.isbn, magazine.ISBN)
				}
				if magazine.Title != tt.expectedTitle {
					t.Errorf("Expected magazine title to be %s, got %s", tt.expectedTitle, magazine.Title)
				}

				if tt.expectedDateStr != "" {
					expectedDate, _ := time.Parse("02.01.2006", tt.expectedDateStr)
					if !magazine.PublishedAt.Equal(expectedDate) {
						t.Errorf("Expected publication date to be %v, got %v", expectedDate, magazine.PublishedAt)
					}
				}
			}
		})
	}
}

func TestMagazineService_FindByAuthorEmail(t *testing.T) {
	service := NewMagazineService()
	err := service.LoadFromCSV("../resources/magazines.csv")
	if err != nil {
		t.Fatalf("Failed to load magazines: %v", err)
	}

	tests := []struct {
		name             string
		email            string
		expectedMinCount int
		expectEmpty      bool
	}{
		{
			name:             "existing author with multiple magazines",
			email:            "null-walter@echocat.org",
			expectedMinCount: 3,
			expectEmpty:      false,
		},
		{
			name:             "author with single magazine",
			email:            "null-mueller@echocat.org",
			expectedMinCount: 1,
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
			magazines := service.FindByAuthorEmail(tt.email)

			if tt.expectEmpty && len(magazines) != 0 {
				t.Errorf("Expected empty slice for %s, but got %d magazines", tt.email, len(magazines))
				return
			}

			if !tt.expectEmpty && len(magazines) == 0 {
				t.Errorf("Expected to find magazines for author %s, but got empty slice", tt.email)
				return
			}

			if len(magazines) < tt.expectedMinCount {
				t.Errorf("Expected at least %d magazines for author %s, got %d", tt.expectedMinCount, tt.email, len(magazines))
			}

			// Verify that all returned magazines contain the author
			for _, magazine := range magazines {
				found := false
				for _, author := range magazine.Authors {
					if author == tt.email {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Magazine %s does not contain author %s", magazine.Title, tt.email)
				}
			}
		})
	}
}

func TestMagazineService_GetAll(t *testing.T) {
	service := NewMagazineService()

	// Test empty service
	magazines := service.GetAll()
	if len(magazines) != 0 {
		t.Errorf("Expected empty slice for new service, got %d magazines", len(magazines))
	}

	// Load magazines and test again
	err := service.LoadFromCSV("../resources/magazines.csv")
	if err != nil {
		t.Fatalf("Failed to load magazines: %v", err)
	}

	magazines = service.GetAll()
	if len(magazines) == 0 {
		t.Error("Expected non-empty slice after loading magazines")
	}

	// Verify that modifying the returned slice doesn't affect the service
	originalLength := len(magazines)
	_ = append(magazines, Magazine{Title: "Test Magazine"})

	magazinesAfterModification := service.GetAll()
	if len(magazinesAfterModification) != originalLength {
		t.Error("Modifying returned slice should not affect the service's internal state")
	}
}

func TestMagazineService_DateParsing(t *testing.T) {
	tests := []struct {
		name          string
		csvContent    string
		expectError   bool
		expectedDates []string
	}{
		{
			name: "valid dates",
			csvContent: `title;isbn;authors;publishedAt
Test Magazine 1;1234-5678-9012;author@example.com;01.01.2020
Test Magazine 2;1234-5678-9013;author@example.com;31.12.2021`,
			expectError:   false,
			expectedDates: []string{"01.01.2020", "31.12.2021"},
		},
		{
			name: "single valid date",
			csvContent: `title;isbn;authors;publishedAt
Test Magazine;1234-5678-9012;author@example.com;15.06.2019`,
			expectError:   false,
			expectedDates: []string{"15.06.2019"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary CSV file
			tempFile, err := os.CreateTemp("", "date_test_magazines_*.csv")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer func() {
				if err := os.Remove(tempFile.Name()); err != nil {
					t.Errorf("Failed to remove temp file: %v", err)
				}
			}()

			if _, err := tempFile.WriteString(tt.csvContent); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			if err := tempFile.Close(); err != nil {
				t.Fatalf("Failed to close temp file: %v", err)
			}

			service := NewMagazineService()
			err = service.LoadFromCSV(tempFile.Name())

			if tt.expectError && err == nil {
				t.Errorf("Expected error for %s, but got nil", tt.name)
				return
			}

			if !tt.expectError && err != nil {
				t.Errorf("Expected no error for %s, but got: %v", tt.name, err)
				return
			}

			if !tt.expectError {
				magazines := service.GetAll()
				if len(magazines) != len(tt.expectedDates) {
					t.Errorf("Expected %d magazines, got %d", len(tt.expectedDates), len(magazines))
				}

				for i, expectedDateStr := range tt.expectedDates {
					if i < len(magazines) {
						expectedDate, _ := time.Parse("02.01.2006", expectedDateStr)
						if !magazines[i].PublishedAt.Equal(expectedDate) {
							t.Errorf("Expected magazine %d date to be %v, got %v", i, expectedDate, magazines[i].PublishedAt)
						}
					}
				}
			}
		})
	}
}

// Benchmark tests
func BenchmarkMagazineService_FindByISBN(b *testing.B) {
	service := NewMagazineService()
	err := service.LoadFromCSV("../resources/magazines.csv")
	if err != nil {
		b.Fatalf("Failed to load magazines: %v", err)
	}

	isbn := "5454-5587-3210"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.FindByISBN(isbn)
	}
}

func BenchmarkMagazineService_FindByAuthorEmail(b *testing.B) {
	service := NewMagazineService()
	err := service.LoadFromCSV("../resources/magazines.csv")
	if err != nil {
		b.Fatalf("Failed to load magazines: %v", err)
	}

	email := "null-walter@echocat.org"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.FindByAuthorEmail(email)
	}
}

func BenchmarkMagazineService_LoadFromCSV(b *testing.B) {
	for i := 0; i < b.N; i++ {
		service := NewMagazineService()
		err := service.LoadFromCSV("../resources/magazines.csv")
		if err != nil {
			b.Fatalf("Failed to load magazines: %v", err)
		}
	}
}
