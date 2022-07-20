package book

type Books []*Book

type Book struct {
	*Fields
}

func NewBook() *Book {
	return &Book{NewFields()}
}
