package entity

import "net/url"

type ModelEntity interface {
	GetPK() string
}

type Image interface {
	GetPath() *url.URL
}
