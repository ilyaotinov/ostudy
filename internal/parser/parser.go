package parser

import (
	"errors"
	"fmt"
	"io"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	internalparser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

const h1 = 1
const h2 = 2

type Parser struct {
	contentReader  io.Reader
	internalParser internalparser.Parser
}

func New(contentReader io.Reader) *Parser {
	internalParser := goldmark.New(
		goldmark.WithExtensions(
			extension.TaskList,
		)).Parser()
	return &Parser{
		contentReader:  contentReader,
		internalParser: internalParser,
	}
}

func (p *Parser) Parse() (Note, error) {
	content, err := io.ReadAll(p.contentReader)
	if err != nil {
		return Note{}, fmt.Errorf("failed to read content: %w", err)
	}
	reader := text.NewReader(content)
	doc := p.internalParser.Parse(reader)
	note := &Note{}
	var lastSection string
	var hasTaskSection bool

	err = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		switch n.Kind() {
		case ast.KindHeading:
			header := n.(*ast.Heading)
			h := processHeader(header, content)
			switch h.level {
			case h1:
				note.title = h.text
			case h2:
				lastSection = h.text
				if h.text == "Task" {
					hasTaskSection = true
				}
			}
		case ast.KindList:
			if lastSection == "Task" {
				list := n.(*ast.List)
				taskList, err := processTaskList(list, content)
				if err != nil {
					return ast.WalkStop, err
				}
				note.taskList = taskList
			}
		}
		return ast.WalkContinue, nil
	})

	if err != nil {
		return Note{}, fmt.Errorf("failed to parse markdown: %w", err)
	}

	if !hasTaskSection {
		return Note{}, errors.New("task section is missing")
	}

	if note.title == "" {
		return Note{}, errors.New("note title is missing")
	}

	if len(note.taskList) == 0 {
		return Note{}, errors.New("task section is missing")
	}

	return *note, nil
}
