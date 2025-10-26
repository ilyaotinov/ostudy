package parser

type Note struct {
	title    string
	taskList []Task
}

func (n Note) Title() string {
	return n.title
}

func (n Note) PercentOfCompleteness() float64 {
	var completedTasks int
	for _, task := range n.taskList {
		if task.isCompleted {
			completedTasks++
		}
	}

	return float64(completedTasks) / float64(len(n.taskList)) * 100
}

func (n Note) TaskList() []Task {
	return n.taskList
}

// Task represent single task
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

type header struct {
	text  string
	level int
}
