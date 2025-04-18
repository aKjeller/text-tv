package tui

import (
	"log"
	"strings"

	"github.com/aKjeller/text-tv/pkg/svttext"
	tea "github.com/charmbracelet/bubbletea/v2"
)

type Model struct {
	page  svttext.Page
	index int
}

func (m Model) Init() tea.Cmd {
	return getPage("377")
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "j":
			return m, getPage(m.page.NextPage)
		case "k":
			return m, getPage(m.page.PrevPage)
		case "h":
			m.index = m.prevIndex()
		case "l":
			m.index = m.nextIndex()
		}
	case newPage:
		m.page = svttext.Page(msg)
		m.index = 0
	}
	return m, nil
}

func (m Model) View() string {
	if len(m.page.SubPages) == 0 {
		return ""
	}

	var sb strings.Builder
	for _, row := range m.page.SubPages[m.index].Grid {
		sb.WriteString(row.ColorString())
		sb.WriteString("\n")
	}

	for i := range m.page.SubPages {
		if i == m.index {
			sb.WriteRune('•')
		} else {
			sb.WriteRune('◦')
		}
	}

	return sb.String()
}

func (m Model) nextIndex() int {
	if m.index < len(m.page.SubPages)-1 {
		return m.index + 1
	}
	return m.index
}

func (m Model) prevIndex() int {
	if m.index > 0 {
		return m.index - 1
	}
	return m.index
}

type newPage svttext.Page

func getPage(pageId string) tea.Cmd {
	return func() tea.Msg {
		page, err := svttext.GetPage(pageId)
		if err != nil {
			log.Fatalf("failed to get page: %v", err)
		}
		return newPage(page)
	}
}
