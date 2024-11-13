package main

import (
	"fmt"

	crf "github.com/GSA-TTS/jemison/cmd/admin-cli/crawl_request_form"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
	dotChar           = " • "
)

// General stuff for styling the view
var (
	keywordStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	ticksStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("79"))
	checkboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	dotStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	mainStyle     = lipgloss.NewStyle().MarginLeft(2)
)

// MODEL DATA

type viewselection int

const (
	MainMenu viewselection = iota
	NextPage
	CrawlRequest
)

func selectView(n int) viewselection {
	return []viewselection{
		// 0          1        2         3
		MainMenu, CrawlRequest, NextPage, NextPage, NextPage,
	}[n]
}

type model struct {
	Choice   int
	Chosen   viewselection
	CRF      tea.Model // huh.Form is just a tea.Model
	Loaded   bool
	Quitting bool
}

func NewModel() model {
	m := model{0, MainMenu, nil, false, false}
	m.CRF = crf.NewModel()
	m.CRF.Init()
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

// Main update function.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if m.Chosen == MainMenu {
		if msg, ok := msg.(tea.KeyMsg); ok {
			k := msg.String()
			switch k {
			case "q":
				fallthrough
			case "esc":
				fallthrough
			case "ctrl+c":
				m.Quitting = true
				return m, tea.Quit
			case "m":
				m.Chosen = MainMenu
				return m, nil
			}
		} else if m.Chosen == CrawlRequest {
			m.CRF.Update(msg)
		}

	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	switch m.Chosen {
	case MainMenu:
		return updateChoices(msg, m)
	default:
		return updateChosen(msg, m)
	}
}

// The main view, which just calls the appropriate sub-view
func (m model) View() string {
	var s string
	if m.Quitting {
		return "\n  See you later!\n\n"
	}

	switch m.Chosen {
	case MainMenu:
		s = choicesView(m)
	case NextPage:
		s = chosenView(m)
	case CrawlRequest:
		s = m.CRF.View()
	}

	return mainStyle.Render("\n" + s + "\n\n")
}

// As we go up-and-down the menu, the model is updated.
// Specifically, the `Choice` field is updated with an integer value
// so that we know what menu to jump to when the enter key is pressed.
// The boolean flag `Chosen` tells the event loop handler that we can
// jump somewhere.
func updateChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.Choice++
			if m.Choice > 3 {
				m.Choice = 3
			}
		case "k", "up":
			m.Choice--
			if m.Choice < 0 {
				m.Choice = 0
			}
		case "enter":
			m.Chosen = selectView(m.Choice + 1)
			// return m, frameMsg{}
		}
	}

	return m, nil
}

// Update loop for the second view after a choice has been made
func updateChosen(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	return m, nil
}

// The first view, where you're choosing a task
func choicesView(m model) string {
	c := m.Choice

	tpl := "jemison pilot\n\n"
	tpl += "%s\n\n"
	tpl += subtleStyle.Render("j/k, up/down: select") + dotStyle +
		subtleStyle.Render("m: main menu") + dotStyle +
		subtleStyle.Render("enter: choose") + dotStyle +
		subtleStyle.Render("q, esc: quit")

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n",
		checkbox("launch a site crawl", c == 0),
		checkbox("monitor jobs", c == 1),
		checkbox("get toes in the grass", c == 2),
	)

	return fmt.Sprintf(tpl, choices)
}

func checkbox(label string, checked bool) string {
	if checked {
		return checkboxStyle.Render("[x] " + label)
	}
	return fmt.Sprintf("[ ] %s", label)
}

// The second view, after a task has been chosen
func chosenView(m model) string {
	var msg string

	switch m.Choice {
	case 0:
		msg = fmt.Sprintf("0 Carrot planting?\n\nCool, we'll need %s and %s...", keywordStyle.Render("libgarden"), keywordStyle.Render("vegeutils"))
	case 1:
		msg = fmt.Sprintf("1 A trip to the market?\n\nOkay, then we should install %s and %s...", keywordStyle.Render("marketkit"), keywordStyle.Render("libshopping"))
	case 2:
		msg = fmt.Sprintf("2 Reading time?\n\nOkay, cool, then we’ll need a library. Yes, an %s.", keywordStyle.Render("actual library"))
	default:
		msg = fmt.Sprintf("It’s always good to see friends.\n\nFetching %s and %s...", keywordStyle.Render("social-skills"), keywordStyle.Render("conversationutils"))
	}

	return msg + "\n\n"
}

func main() {
	initialModel := NewModel()
	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}
