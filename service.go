package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"golang.org/x/text/language"
)

// Service is a Translator user.
type Service struct {
	translator Translator
}

func NewService() *Service {
	t := newRandomTranslator(
		100*time.Millisecond,
		500*time.Millisecond,
		0.1,
	)

	return &Service{
		translator: t,
	}
}

var backoffSchedule = []time.Duration{
	1 * time.Second,
	3 * time.Second,
	10 * time.Second,
}

func (s *Service) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	var resp string
	var err error
	for _, backoff := range backoffSchedule {
		resp, err = s.translator.Translate(ctx, from, to, data)

		if err == nil {
			break
		}
		fmt.Fprintf(os.Stderr, "Retrying in %v\n", backoff)
	}
	if err != nil {
		return "", err
	}

	return resp, err
}
