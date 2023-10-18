package huh

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/accessibility"
	"github.com/charmbracelet/lipgloss"
)

// MultiSelect is a form multi-select field.
type MultiSelect[T any] struct {
	title            string
	required         bool
	filterable       bool
	limit            int
	cursor           int
	cursorPrefix     string
	selectedPrefix   string
	unselectedPrefix string
	selected         []bool
	options          []Option[T]
	value            *[]T
	style            *MultiSelectStyle
	blurredStyle     MultiSelectStyle
	focusedStyle     MultiSelectStyle
}

// NewMultiSelect returns a new multi-select field.
func NewMultiSelect[T any](options ...T) *MultiSelect[T] {
	f, b := DefaultMultiSelectStyles()

	var opts []Option[T]
	for _, o := range options {
		opts = append(opts, Option[T]{Key: fmt.Sprint(o), Value: o})
	}

	return &MultiSelect[T]{
		options:          opts,
		cursorPrefix:     "> ",
		selectedPrefix:   "[•] ",
		unselectedPrefix: "[ ] ",
		focusedStyle:     f,
		blurredStyle:     b,
		selected:         make([]bool, len(opts)),
	}
}

// Value sets the value of the multi-select field.
func (m *MultiSelect[T]) Value(value *[]T) *MultiSelect[T] {
	m.value = value
	return m
}

// Title sets the title of the multi-select field.
func (m *MultiSelect[T]) Title(title string) *MultiSelect[T] {
	m.title = title
	return m
}

// Required sets the multi-select field as required.
func (m *MultiSelect[T]) Required(required bool) *MultiSelect[T] {
	m.required = required
	return m
}

// Options sets the options of the multi-select field.
func (m *MultiSelect[T]) Options(options ...Option[T]) *MultiSelect[T] {
	m.options = options
	return m
}

// Filterable sets the multi-select field as filterable.
func (m *MultiSelect[T]) Filterable(filterable bool) *MultiSelect[T] {
	m.filterable = filterable
	return m
}

// Cursor sets the cursor of the multi-select field.
func (m *MultiSelect[T]) Cursor(cursor string) *MultiSelect[T] {
	m.cursorPrefix = cursor
	return m
}

// Limit sets the limit of the multi-select field.
func (m *MultiSelect[T]) Limit(limit int) *MultiSelect[T] {
	m.limit = limit
	return m
}

// Focus focuses the multi-select field.
func (m *MultiSelect[T]) Focus() tea.Cmd {
	m.style = &m.focusedStyle
	return nil
}

// Blur blurs the multi-select field.
func (m *MultiSelect[T]) Blur() tea.Cmd {
	m.style = &m.blurredStyle
	return nil
}

// Init initializes the multi-select field.
func (m *MultiSelect[T]) Init() tea.Cmd {
	m.style = &m.blurredStyle
	return nil
}

// Update updates the multi-select field.
func (m *MultiSelect[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.cursor = max(m.cursor-1, 0)
		case "down", "j":
			m.cursor = min(m.cursor+1, len(m.options)-1)
		case " ", "x":
			m.selected[m.cursor] = !m.selected[m.cursor]
		case "shift+tab":
			m.finalize()
			return m, prevField
		case "tab", "enter":
			m.finalize()
			return m, nextField
		}
	}

	return m, nil
}

func (m *MultiSelect[T]) finalize() {
	*m.value = make([]T, 0)
	for i, option := range m.options {
		if m.selected[i] {
			*m.value = append(*m.value, option.Value)
		}
	}
}

// View renders the multi-select field.
func (m *MultiSelect[T]) View() string {
	var sb strings.Builder
	sb.WriteString(m.style.Title.Render(m.title) + "\n")
	c := m.style.Cursor.Render(m.cursorPrefix)
	for i, option := range m.options {
		if m.cursor == i {
			sb.WriteString(c)
		} else {
			sb.WriteString(strings.Repeat(" ", lipgloss.Width(c)))
		}

		if m.selected[i] {
			sb.WriteString(m.style.SelectedPrefix.Render(m.selectedPrefix))
			sb.WriteString(m.style.Selected.Render(option.Key))
		} else {
			sb.WriteString(m.style.UnselectedPrefix.Render(m.unselectedPrefix))
			sb.WriteString(m.style.Unselected.Render(option.Key))
		}
		if i < len(m.options)-1 {
			sb.WriteString("\n")
		}
	}
	return m.style.Base.Render(sb.String())
}

// Run runs the multi-select field in accessible mode.
func (m *MultiSelect[T]) Run() {
	fmt.Println(m.style.Title.Render(m.title))

	for i, option := range m.options {
		fmt.Printf("%d. %s\n", i+1, option.Key)
	}

	fmt.Println("\nType 0 to finish.\n")

	var choice int
	for {
		choice = accessibility.PromptInt(0, len(m.options))
		if choice == 0 {
			break
		}
		if m.selected[choice-1] {
			fmt.Println("Selected:", m.options[choice-1].Key)
		} else {
			fmt.Println("Unselected:", m.options[choice-1].Key)
		}
	}

	for i, option := range m.options {
		if m.selected[i] {
			*m.value = append(*m.value, option.Value)
		}
	}
}