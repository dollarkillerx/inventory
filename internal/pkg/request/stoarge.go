package request

import "github.com/dollarkillerx/RubiesCube/internal/pkg/models"

type Storage struct {
	Bucket  string         `json:"bucket" binding:"required"`
	Key     string         `json:"key" binding:"required"`
	Search  string         `json:"search"`
	Payload models.JSONMap `json:"payload"`
}

type Update struct {
	Bucket  string         `json:"bucket" binding:"required"`
	Key     string         `json:"key" binding:"required"`
	Payload models.JSONMap `json:"payload"`
}

type Delete struct {
	Bucket string `json:"bucket" binding:"required"`
	Key    string `json:"key" binding:"required"`
}

type Query struct {
	Bucket string   `json:"bucket" binding:"required"`
	Key    string   `json:"key"`
	Search string   `json:"search"`
	Limit  int      `json:"limit"`
	Offset int      `json:"offset"`
	Equals []Equals `json:"equals"`
}

type Equals struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
