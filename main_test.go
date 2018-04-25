package main

import (
	"testing"
)

func TestWelcomeMessage(t *testing.T) {
	// given
	expected := "Hello world!"

	// than
	actual := welcomeMessage()

	// that
	if actual != expected {
		t.Errorf("Message returned by welcomeMessage is '%s' - but expected was '%s'", actual, expected)
	}
}
