package parser

import (
	"io"
)

type Content struct {
	Text string
}

type ParseResult struct {
	Content Content
	Error   error
}

type Parser interface {
	Invoke(reader io.Reader) (*ParseResult, error)
}
