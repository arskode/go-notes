package model_test

import (
	"testing"

	"github.com/arskode/go-notes/api/model"
	"github.com/stretchr/testify/assert"
)

func TestValidation(t *testing.T) {
	testCases := []struct {
		Title        string
		Description  string
		errorMessage string
	}{
		{
			Title:        "Foo",
			Description:  "",
			errorMessage: "Required Description",
		},
		{
			Title:        "",
			Description:  "Bar",
			errorMessage: "Required Title",
		},
	}

	for _, v := range testCases {
		note := model.Note{Title: v.Title, Description: v.Description}
		err := note.Validate()
		assert.Equal(t, err.Error(), v.errorMessage)
	}
}
