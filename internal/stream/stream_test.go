package stream_test

import (
	"testing"

	"github.com/GlebYaltchik/sc-keybind-extract/internal/stream"
)

func TestStream_ReadCString(t *testing.T) {
	t.Parallel()

	data := []byte("\x00test1\x00\x00test2")

	s := stream.New(data)

	if got, want := s.ReadCString(), ""; got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	if got, want := s.ReadCString(), "test1"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	if got, want := s.ReadCString(), ""; got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	if got, want := s.ReadCString(), "test2"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	if got, want := s.ReadCString(), ""; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestStream_ReadFString(t *testing.T) {
	t.Parallel()

	data := []byte("\x00test1\x00\x00test2")

	s := stream.New(data)

	if got, want := s.ReadFString(0), ""; got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	if got, want := s.ReadFString(1), ""; got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	if got, want := s.ReadFString(2), "te"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	if got, want := s.ReadFString(6), "st1"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	if got, want := s.ReadFString(10), "est2"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
