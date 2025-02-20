package main

import "github.com/markponce/snippetbox/internal/models"

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}
