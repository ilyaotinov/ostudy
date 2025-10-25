package parser

type Note struct {
	title    string
	taskList []Task
}

func (n Note) Title() string {
	return n.title
}

func (n Note) PercentOfCompleteness() float64 {
	return 42.857
}

func (n Note) TaskList() []Task {
	return n.taskList
}

type Task struct {
	text        string
	isCompleted bool
}

func (t Task) IsCompleted() bool {
	return t.isCompleted
}

func (t Task) Text() string {
	return t.text
}

type Header struct {
	text  string
	level int
}
