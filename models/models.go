package models

import (
	"time"
)

type Config struct {
	Port     string
	DBConfig DBConfig
}

type DBConfig struct {
	Host         string
	Port         string
	UserName     string
	Password     string
	DatabaseName string
}

type ImageInfo struct {
	ID              string
	FileName        string
	ImageSize       int64 //in byte
	FileContentType string
	AcceptEncoding  string
	AcceptLanguage  string
	ContentType     string
	CreatedAt       time.Time
	CreatedBy       string
}
