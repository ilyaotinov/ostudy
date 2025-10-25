package parser

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
)

const taskListMD = `## Task
- [ ] Task 1
- [x] Task 2`

const taskWithoutContentMD = `## Task
- [x]
- [x] Task 2`

const emptyTaskListMD = `## Task
- [ ]`

const invalidTaskListMD = `## Task
- [ ] Task 1
- Task 2`

var parser = goldmark.New(
	goldmark.WithExtensions(
		extension.TaskList, // ← Важно!
	),
).Parser()

func TestProcessMainHeader(t *testing.T) {
	// Arrange.
	markdown := []byte("# Main header")
	heading := createHeading(t, markdown)

	// Act.
	got := processHeader(heading, markdown)

	// Assert.
	assert.Equal(t, "Main header", got.text)
	assert.Equal(t, 1, got.level)
}

func TestProcessH3Header(t *testing.T) {
	// Arrange.
	markdown := []byte("### h3 header")
	heading := createHeading(t, markdown)

	// Act.
	got := processHeader(heading, markdown)

	// Assert.
	assert.Equal(t, "h3 header", got.text)
	assert.Equal(t, 3, got.level)
}

func TestProcessH2HeaderNotATask(t *testing.T) {
	// Arrange.
	markdown := []byte("### h3 header")
	heading := createHeading(t, markdown)
	// Act.
	got := processHeader(heading, markdown)
	// Assert.
	assert.Equal(t, "h3 header", got.text)
	assert.Equal(t, 3, got.level)
}

func TestProcessH2HeaderTask(t *testing.T) {
	markdown := []byte(taskListMD)
	heading := createHeading(t, markdown)

	got := processHeader(heading, markdown)
	assert.Equal(t, "Task", got.text)
	assert.Equal(t, 2, got.level)
}

func TestProcessTaskList(t *testing.T) {
	tests := []struct {
		name           string
		markdown       string
		expectTaskList []Task
		expectError    error
	}{
		{
			name:     "simple case",
			markdown: taskListMD,
			expectTaskList: []Task{
				{
					text:        "Task 1",
					isCompleted: false,
				},
				{
					text:        "Task 2",
					isCompleted: true,
				},
			},
		},
		{
			name:     "case without task content",
			markdown: taskWithoutContentMD,
			expectTaskList: []Task{
				{text: "Task 2", isCompleted: true},
			},
		},
		{
			name:           "empty task list",
			markdown:       emptyTaskListMD,
			expectTaskList: []Task{},
		},
		{
			name:        "task list contains invalid plain list elements",
			markdown:    invalidTaskListMD,
			expectError: errors.New("failed to process task list: task list expect contains only checkbox items"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			markdown := []byte(tt.markdown)
			list := createCheckboxList(t, markdown)

			// Act.
			got, err := processTaskList(list, markdown)

			if tt.expectError != nil {
				assert.EqualError(t, err, tt.expectError.Error())
				return
			}
			// Assert.
			for _, task := range tt.expectTaskList {
				assert.True(t, containsWithText(got, task.text, task.isCompleted))
			}
			assert.Equal(t, len(tt.expectTaskList), len(got))
		})
	}
}

func containsWithText(items []Task, text string, completed bool) bool {
	for _, item := range items {
		if item.IsCompleted() == completed && item.Text() == text {
			return true
		}
	}
	return false
}

func createCheckboxList(t *testing.T, markdown []byte) *ast.List {
	t.Helper()

	doc := parser.Parse(text.NewReader(markdown))

	var checkboxList *ast.List
	err := ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if list, ok := node.(*ast.List); ok {
			checkboxList = list
			return ast.WalkStop, nil
		}

		return ast.WalkContinue, nil
	})
	require.NoError(t, err)

	return checkboxList
}

func createHeading(t *testing.T, markdownHeader []byte) *ast.Heading {
	t.Helper()

	reader := text.NewReader(markdownHeader)
	doc := parser.Parse(reader)
	var heading *ast.Heading
	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && n.Kind() == ast.KindHeading {
			heading = n.(*ast.Heading)
			return ast.WalkStop, nil
		}
		return ast.WalkContinue, nil
	})

	require.NoError(t, err)

	return heading
}
