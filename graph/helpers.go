package graph

import (
	"errors"
	"time"
)

var ErrUnauthorized = errors.New("unauthorized")

func coalesceStrPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func timePtrToStringPtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.UTC().Format(time.RFC3339Nano)
	return &s
}
