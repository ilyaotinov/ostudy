package parser_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/ilyaotinov/ostudy/internal/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const simpleMDData = `# Progress note
## Some title
some interesting description and maybe another thing
## Task
- [x] Task 1
- [x] Task 2
- [x] Task 3
- [ ] Task 4
- [ ] Task 5
- [ ] Task 6
- [ ] Task 7`

const withoutTaskSectionMD = `# Progress note
## Some title
some interesting description and maybe another thing
- [x] Task 1
- [x] Task 2
- [x] Task 3
- [ ] Task 4
- [ ] Task 5
- [ ] Task 6
- [ ] Task 7`

const missingNoteTitleMD = `## Some title
some interesting description and maybe another thing
## Task
- [x] Task 1
- [x] Task 2
- [x] Task 3
- [ ] Task 4
- [ ] Task 5
- [ ] Task 6
- [ ] Task 7`

const taskHeaderInWrongPlace = `# Progress note
## Task
Some text
## Some title
some interesting description and maybe another thing
- [x] Task 1
- [x] Task 2
- [x] Task 3
- [ ] Task 4
- [ ] Task 5
- [ ] Task 6
- [ ] Task 7`

func TestMDParser_Parse_Success(t *testing.T) {
	reader := strings.NewReader(simpleMDData)
	p := parser.New(reader)
	got, err := p.Parse()
	require.NoError(t, err)
	assert.Equal(t, got.Title(), "Progress note")
	assert.InDelta(t, got.PercentOfCompleteness(), 42.857, 0.001)
	assert.True(t, containsWithText(got.TaskList(), "Task 1", true))
	assert.True(t, containsWithText(got.TaskList(), "Task 2", true))
	assert.True(t, containsWithText(got.TaskList(), "Task 3", true))
	assert.True(t, containsWithText(got.TaskList(), "Task 4", false))
	assert.True(t, containsWithText(got.TaskList(), "Task 5", false))
	assert.True(t, containsWithText(got.TaskList(), "Task 6", false))
	assert.True(t, containsWithText(got.TaskList(), "Task 7", false))
}

func TestMDParser_Parse_MissingTaskSection(t *testing.T) {
	reader := strings.NewReader(withoutTaskSectionMD)
	p := parser.New(reader)
	_, err := p.Parse()
	require.Error(t, err)
	require.EqualError(t, errors.New("task section is missing"), err.Error())
}

func TestMDParser_Parse_MissingNoteTitle(t *testing.T) {
	reader := strings.NewReader(missingNoteTitleMD)
	p := parser.New(reader)
	_, err := p.Parse()
	require.Error(t, err)
	require.EqualError(t, errors.New("note title is missing"), err.Error())
}

func TestMDParser_Parse_TaskHeaderInWrongPlace(t *testing.T) {
	reader := strings.NewReader(taskHeaderInWrongPlace)
	p := parser.New(reader)
	_, err := p.Parse()
	require.Error(t, err)
	require.EqualError(t, errors.New("task section is missing"), err.Error())
}

func containsWithText(items []parser.Task, text string, completed bool) bool {
	for _, item := range items {
		if item.IsCompleted() == completed && item.Text() == text {
			return true
		}
	}
	return false
}
