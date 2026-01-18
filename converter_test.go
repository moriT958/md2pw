package main

import (
	"testing"
)

func TestConvert_H1(t *testing.T) {
	input := []byte("# Heading 1")
	expected := "* Heading 1"

	result, err := Convert(input)
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestConvert_H2(t *testing.T) {
	input := []byte("## Heading 2")
	expected := "** Heading 2"

	result, err := Convert(input)
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestConvert_H3(t *testing.T) {
	input := []byte("### Heading 3")
	expected := "*** Heading 3"

	result, err := Convert(input)
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestConvert_MultipleHeadings(t *testing.T) {
	input := []byte("# H1\n\n## H2\n\n### H3")
	expected := "* H1\n\n** H2\n\n*** H3"

	result, err := Convert(input)
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestConvert_H4Unchanged(t *testing.T) {
	input := []byte("# H1\n\n#### H4\n\n## H2")
	expected := "* H1\n\n#### H4\n\n** H2"

	result, err := Convert(input)
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestConvert_PreserveContent(t *testing.T) {
	input := []byte("# Title\n\nSome paragraph text.\n\n## Section\n\n- list item")
	expected := "* Title\n\nSome paragraph text.\n\n** Section\n\n- list item"

	result, err := Convert(input)
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}
