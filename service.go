package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/patrickmn/go-cache"
	"golang.org/x/text/language"
)

// Service is a Translator user.
type Service struct {
	translator  Translator
	cache_value *cache.Cache
}

func NewService() *Service {
	t := newRandomTranslator(
		100*time.Millisecond,
		500*time.Millisecond,
		0.1,
	)
	c := cache.New(5*time.Minute, 10*time.Minute)
	fmt.Println("cache created")
	return &Service{
		translator:  t,
		cache_value: c,
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

	request_param := from.String() + to.String() + data
	fmt.Println("requested params: " + request_param)

	resp1, found := s.cache_value.Get(request_param)
	if found {
		fmt.Println("found cache result:" + resp1.(string))
		return resp1.(string), nil
	}

	for _, backoff := range backoffSchedule {
		resp, err = s.translator.Translate(ctx, from, to, data)

		if err == nil {
			s.cache_value.Set(request_param, resp, cache.DefaultExpiration)
			fmt.Println("params inserted to cache")
			break
		}
		fmt.Fprintf(os.Stderr, "Retrying in %v\n", backoff)
	}
	if err != nil {
		return "", err
	}

	return resp, err
}
