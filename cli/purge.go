package cli

import (
	"fmt"
	"github.com/dormunis/gitd/adapters"
	"github.com/dormunis/gitd/taskmanagers/taskmanager"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	taskmanager adapters.TaskManagerAdapter
	cursor      int
	actions     *[]adapters.TaskAction
	tasks       *[]adapters.Task
	keys        keymap
	help        help.Model
}

type keymap struct {
	Up         key.Binding
	Down       key.Binding
	Defer      key.Binding
	Ignore     key.Binding
	Complete   key.Binding
	Revalidate key.Binding
	Delete     key.Binding
	Save       key.Binding
	Quit       key.Binding
	Help       key.Binding
}

func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Save, k.Quit}
}

func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Ignore, k.Complete, k.Defer, k.Delete, k.Revalidate},
		{k.Help, k.Save, k.Quit},
	}
}

var keys = keymap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Complete: key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("v", "complete"),
	),
	Ignore: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "ignore"),
	),
	Defer: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "defer"),
	),
	Delete: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "delete"),
	),
	Revalidate: key.NewBinding(
		key.WithKeys("backspace", "delete"),
		key.WithHelp("backspace/delete", "remove selection"),
	),
	Save: key.NewBinding(
		key.WithKeys("esc", "q"),
		key.WithHelp("esc/q", "save and quit"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit without saving"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
}

func Purge(taskManager adapters.TaskManagerAdapter, timespan adapters.TimeSpan) {
	// TODO: make this use a loader
	tasks, err := taskManager.FetchTasks()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	filterRequest := adapters.FilterRequest{
		AfterTimeSpan: &timespan,
	}
	filteredTasks, err := taskmanager.FilterTasks(&tasks, &filterRequest)
	actions := make([]adapters.TaskAction, len(*&filteredTasks))

	programModel := model{
		tasks:   &filteredTasks,
		actions: &actions,
		keys:    keys,
		help:    help.New(),
	}

	p := tea.NewProgram(programModel, tea.WithAltScreen())
	_, err = p.Run()
	if err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}

	SavePurge(taskManager, &actions)
}

func (m model) Init() tea.Cmd {
	for i := range *m.tasks {
		(*m.actions)[i].Task = &(*m.tasks)[i]
		(*m.actions)[i].Action = adapters.ActionRevalidate
	}
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(*m.tasks)-1 {
				m.cursor++
			}
		case key.Matches(msg, m.keys.Revalidate):
			(*m.actions)[m.cursor].Action = adapters.ActionRevalidate

		case key.Matches(msg, m.keys.Ignore):
			(*m.actions)[m.cursor].Action = adapters.ActionIgnore
			if m.cursor < len(*m.tasks)-1 {
				m.cursor++
			}
		case key.Matches(msg, m.keys.Complete):
			(*m.actions)[m.cursor].Action = adapters.ActionComplete
			if m.cursor < len(*m.tasks)-1 {
				m.cursor++
			}
		case key.Matches(msg, m.keys.Defer):
			(*m.actions)[m.cursor].Action = adapters.ActionDefer
			if m.cursor < len(*m.tasks)-1 {
				m.cursor++
			}
		case key.Matches(msg, m.keys.Delete):
			(*m.actions)[m.cursor].Action = adapters.ActionDelete
			if m.cursor < len(*m.tasks)-1 {
				m.cursor++
			}
		case key.Matches(msg, m.keys.Save):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Quit):
			os.Exit(1)
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
	}

	return m, nil
}

func (m model) createTable() table.Model {
	columns := []table.Column{
		{Title: "Action", Width: 8},
		{Title: "Task", Width: 65},
		{Title: "Project", Width: 30},
		{Title: "Creation Date", Width: 18},
		{Title: "Last Modified Date", Width: 18},
	}

	var rows []table.Row

	for _, taskAction := range *m.actions {
		mark := getActionString(&taskAction)
		rows = append(rows, table.Row{
			fmt.Sprintf("[%s]", mark),
			taskAction.Task.Content,
			taskAction.Task.Project,
			taskAction.Task.CreatedDate.Format("2006-01-02"),
			taskAction.Task.UpdatedDate.Format("2006-01-02"),
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)
	t.SetCursor(m.cursor)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	t.SetStyles(s)
	return t
}

func (m model) View() string {
	helpView := m.help.View(m.keys)
	t := m.createTable()
	return helpView + "\n\n" + t.View()
}

func getActionString(action *adapters.TaskAction) string {
	switch action.Action {
	case adapters.ActionDelete:
		return "x"
	case adapters.ActionComplete:
		return "v"
	case adapters.ActionDefer:
		return "d"
	case adapters.ActionIgnore:
		return "i"
	case adapters.ActionRevalidate:
		return " "
	default:
		return " "
	}
}

func SavePurge(taskManager interface{ adapters.TaskManagerAdapter }, actions *[]adapters.TaskAction) {
	if !verifyPurge(actions) {
		return
	}
	taskManager.UpdateTasks(actions)
}

func verifyPurge(actions *[]adapters.TaskAction) bool {
	deleteCount := 0
	deferCount := 0
	revalidateCount := 0
	completeCount := 0

	for _, action := range *actions {
		switch action.Action {
		case adapters.ActionDelete:
			deleteCount++
		case adapters.ActionDefer:
			deferCount++
		case adapters.ActionRevalidate:
			revalidateCount++
		case adapters.ActionComplete:
			completeCount++
		}
	}

	if deleteCount == 0 && deferCount == 0 && revalidateCount == 0 && completeCount == 0 {
		return false
	}

	fmt.Println("You are about to perform the following actions:")
	fmt.Printf("Delete %d items\n", deleteCount)
	fmt.Printf("Archiving %d items\n", deferCount)
	fmt.Printf("Revalidating %d items\n", revalidateCount)
	fmt.Printf("Completing %d items\n", completeCount)
	fmt.Printf("Are you sure you want to continue? (y/n): ")
	var response string
	fmt.Scanln(&response)
	return strings.ToLower(response) == "y"
}
