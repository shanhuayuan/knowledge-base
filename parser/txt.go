package parser

import (
	"errors"
	"io"
	"unicode/utf8"
)

type TxtParser struct{}

func NewTxtParser() *TxtParser {
	return &TxtParser{}
}

func (p *TxtParser) Invoke(reader io.Reader) (*ParseResult, error) {
	//读取文件
	if reader == nil {
		return nil, errors.New("reader is nil")
	}

	data, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	if !utf8.Valid(data) {
		return nil, errors.New("invalid UTF-8")
	}

	content := Content{
		Text: string(data),
	}

	return &ParseResult{Content: content}, nil
}
