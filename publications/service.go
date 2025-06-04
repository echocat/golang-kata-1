package publications

import (
	"sort"

	"github.com/echocat/golang-kata-1/authors"
	"github.com/echocat/golang-kata-1/books"
	"github.com/echocat/golang-kata-1/magazines"
)

func CollectPublications(
	bookService books.Service,
	magazineService magazines.Service,
	authorService authors.Service,
) []Publication {

	var publications []Publication

	for _, book := range bookService.GetAll() {
		pub := Publication{
			Title:       book.Title,
			ISBN:        book.ISBN,
			Authors:     book.Authors,
			Type:        "Book",
			Description: book.Description,
		}
		publications = append(publications, pub)
	}

	for _, magazine := range magazineService.GetAll() {
		pub := Publication{
			Title:       magazine.Title,
			ISBN:        magazine.ISBN,
			Authors:     magazine.Authors,
			Type:        "Magazine",
			PublishedAt: &magazine.PublishedAt,
		}
		publications = append(publications, pub)
	}

	return publications
}

func FilterByISBN(publications []Publication, isbn string) []Publication {
	var filtered []Publication
	for _, pub := range publications {
		if pub.ISBN == isbn {
			filtered = append(filtered, pub)
		}
	}
	return filtered
}

func FilterByAuthorEmail(publications []Publication, email string) []Publication {
	var filtered []Publication
	for _, pub := range publications {
		for _, author := range pub.Authors {
			if author == email {
				filtered = append(filtered, pub)
				break
			}
		}
	}
	return filtered
}

func SortByTitle(publications []Publication) []Publication {
	// Create a copy to avoid modifying the original slice
	sorted := make([]Publication, len(publications))
	copy(sorted, publications)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Title < sorted[j].Title
	})

	return sorted
}
