package domain

import "net/url"

type Media struct {
	BaseObject
	Name    string
	Tags    []*Tag
	FileUrl url.URL
}
