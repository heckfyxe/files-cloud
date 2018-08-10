package models

type File struct {
	Id   int64
	Name string
	Size string
}

type Files []*File
