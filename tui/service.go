package tui

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/echocat/golang-kata-1/publications"
)

// WrapText wraps text to specified width
func WrapText(text string, width int) []string {
	if len(text) <= width {
		return []string{text}
	}

	var lines []string
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	currentLine := words[0]

	for _, word := range words[1:] {
		if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

// PrintPublicationsTable prints a formatted table of publications
func PrintPublicationsTable(publicationsList []publications.Publication) {
	if len(publicationsList) == 0 {
		fmt.Println("\nNo publications found matching the criteria.")
		return
	}

	fmt.Println("\n--- Publications ---")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	_, _ = fmt.Fprintf(w, "Type\tTitle\tISBN\tAuthors\tDescription/Published At\n")
	_, _ = fmt.Fprintf(
		w, "========\t==========\t========\t==============\t======================================\n",
	)

	for i, pub := range publicationsList {
		if i > 0 {
			_, _ = fmt.Fprintf(
				w, "--------\t----------\t--------\t--------------\t--------------------------------------\n",
			)
		}

		authorsStr := strings.Join(pub.Authors, ", ")

		var lastColumn string
		if pub.Type == "Book" {
			lastColumn = pub.Description
		} else {
			lastColumn = pub.PublishedAt.Format("02.01.2006")
		}

		// Wrap long columns
		titleLines := WrapText(pub.Title, 30)
		authorLines := WrapText(authorsStr, 25)
		lastColumnLines := WrapText(lastColumn, 50)

		// Find the maximum number of lines needed
		maxLines := len(titleLines)
		if len(authorLines) > maxLines {
			maxLines = len(authorLines)
		}
		if len(lastColumnLines) > maxLines {
			maxLines = len(lastColumnLines)
		}

		// Pad all slices to the same length
		for len(titleLines) < maxLines {
			titleLines = append(titleLines, "")
		}
		for len(authorLines) < maxLines {
			authorLines = append(authorLines, "")
		}
		for len(lastColumnLines) < maxLines {
			lastColumnLines = append(lastColumnLines, "")
		}

		// Print each line
		for i := 0; i < maxLines; i++ {
			typeCol := ""
			isbnCol := ""
			if i == 0 {
				typeCol = pub.Type
				isbnCol = pub.ISBN
			}

			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				typeCol, titleLines[i], isbnCol, authorLines[i], lastColumnLines[i])
		}
	}

	if err := w.Flush(); err != nil {
		slog.Error("failed to flush writer", "error", err)
	}
	fmt.Printf("\nTotal publications: %d\n", len(publicationsList))
}
