package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/echocat/golang-kata-1/authors"
	"github.com/echocat/golang-kata-1/books"
	"github.com/echocat/golang-kata-1/magazines"
	pubService "github.com/echocat/golang-kata-1/publications"
	"github.com/echocat/golang-kata-1/tui"
)

func main() {
	var (
		printCmd     = flag.Bool("print", false, "Print all books and magazines")
		sortBy       = flag.String("sort-by", "", "Sort by field (only 'title' is supported)")
		isbnFilter   = flag.String("isbn", "", "Filter by ISBN")
		authorFilter = flag.String("author-email", "", "Filter by author email")
	)

	flag.Parse()

	if !*printCmd {
		fmt.Fprintf(os.Stderr, "Error: The 'print' command is mandatory\n")
		fmt.Fprintf(os.Stderr, "Usage: %s -print [options]\n", os.Args[0])
		os.Exit(1)
	}

	if *sortBy != "" && *sortBy != "title" {
		fmt.Fprintf(os.Stderr, "Error: sort-by flag only accepts 'title' as value\n")
		os.Exit(1)
	}

	authorService := authors.NewAuthorService()
	bookService := books.NewBookService()
	magazineService := magazines.NewMagazineService()

	resourcesPath := "resources"

	if err := authorService.LoadFromCSV(filepath.Join(resourcesPath, "authors.csv")); err != nil {
		log.Fatalf("Failed to load authors: %v", err)
	}

	if err := bookService.LoadFromCSV(filepath.Join(resourcesPath, "books.csv")); err != nil {
		log.Fatalf("Failed to load books: %v", err)
	}

	if err := magazineService.LoadFromCSV(filepath.Join(resourcesPath, "magazines.csv")); err != nil {
		log.Fatalf("Failed to load magazines: %v", err)
	}

	fmt.Printf("\nLoaded %d authors, %d books, and %d magazines\n",
		len(authorService.GetAll()),
		len(bookService.GetAll()),
		len(magazineService.GetAll()))

	publicationsList := pubService.CollectPublications(bookService, magazineService, authorService)

	if *isbnFilter != "" {
		publicationsList = pubService.FilterByISBN(publicationsList, *isbnFilter)
	}

	if *authorFilter != "" {
		publicationsList = pubService.FilterByAuthorEmail(publicationsList, *authorFilter)
	}

	if *sortBy == "title" {
		publicationsList = pubService.SortByTitle(publicationsList)
	}

	tui.PrintPublicationsTable(publicationsList)
}
