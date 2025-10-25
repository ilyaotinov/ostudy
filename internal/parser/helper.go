package parser

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
)

func processHeader(header *ast.Heading, source []byte) Header {
	return Header{text: extractText(header, source), level: header.Level}
}

// processTaskList assume it is a list after ## Task header.
func processTaskList(list *ast.List, source []byte) ([]Task, error) {
	result := make([]Task, 0)
	err := ast.Walk(list, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if n.Kind() != ast.KindListItem {
			return ast.WalkContinue, nil
		}

		item := n.(*ast.ListItem)

		var isTask bool
		var isChecked bool

		err := ast.Walk(item, func(child ast.Node, childEntering bool) (ast.WalkStatus, error) {
			if !childEntering {
				return ast.WalkContinue, nil
			}

			if child.Kind() == extast.KindTaskCheckBox {
				isTask = true
				checkbox := child.(*extast.TaskCheckBox)
				isChecked = checkbox.IsChecked
				return ast.WalkStop, nil
			}

			return ast.WalkContinue, nil
		})
		if err != nil {
			return ast.WalkStop, fmt.Errorf("failed to process checkbox: %w", err)
		}

		if !isTask {
			return ast.WalkStop, errors.New("task list expect contains only checkbox items")
		}

		text := extractText(item, source)
		if len(strings.TrimSpace(text)) > 0 {
			result = append(result, Task{
				text:        text,
				isCompleted: isChecked,
			})
		}

		return ast.WalkContinue, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to process task list: %w", err)
	}

	return result, nil
}

func extractText(n ast.Node, source []byte) string {
	var buf bytes.Buffer

	err := ast.Walk(n, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && node.Kind() == ast.KindText {
			textNode := node.(*ast.Text)
			buf.Write(textNode.Segment.Value(source))
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		panic("it should not happen")
	}

	return buf.String()
}
