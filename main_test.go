package main

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestWelcomeMessage(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	expected := "Hello world!"

	// than
	actual := welcomeMessage()

	// that
	g.Expect(actual).To(Equal(expected))
}
