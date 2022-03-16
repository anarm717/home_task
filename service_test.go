package main

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestTranslate(t *testing.T) {
	ctx := context.Background()
	rand.Seed(time.Now().UTC().UnixNano())
	s := NewService()
	resp, err := s.Translate(ctx, language.English, language.Japanese, "test")
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	request_param := language.English.String() + language.Japanese.String() + "test"
	_, found := s.cache_value.Get(request_param)
	if !found {
		t.Error("Cache not found")
	}
}
