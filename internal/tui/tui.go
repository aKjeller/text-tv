package tui

import (
	"os"
	"strings"

	"github.com/aKjeller/text-tv/pkg/svttext"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"golang.org/x/term"
)

type Model struct {
	page svttext.Page

	// SubPage Pagination
	activeIndex int
	activeDot   string
	inactiveDot string

	client *svttext.Client
}

func NewModel() Model {
	return Model{
		client:      svttext.NewClient(svttext.WithCaching()),
		page:        svttext.Page{},
		activeIndex: 0,
		activeDot:   lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render("●"),
		inactiveDot: lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Render("●"),
	}
}

func (m Model) Init() tea.Cmd {
	return m.getPage("100")
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "k":
			return m, m.getPage(m.page.NextPage)
		case "j":
			return m, m.getPage(m.page.PrevPage)
		case "h":
			m.activeIndex = m.prevIndex()
		case "l":
			m.activeIndex = m.nextIndex()
		}
	case newPage:
		m.page = svttext.Page(msg)
		m.activeIndex = 0
		return m, tea.Batch(m.preLoadPage(m.page.PrevPage), m.preLoadPage(m.page.NextPage))
	}
	return m, nil
}

func (m Model) View() string {
	if len(m.page.SubPages) == 0 {
		return ""
	}

	var sb strings.Builder
	for _, row := range m.page.SubPages[m.activeIndex].Grid {
		sb.WriteString(row.ColorString())
		sb.WriteString("\n")
	}

	sb.WriteString(" ")
	for i := range m.page.SubPages {
		if i == m.activeIndex {
			sb.WriteString(m.activeDot + " ")
		} else {
			sb.WriteString(m.inactiveDot + " ")
		}
	}

	// TODO: add debug options
	// sb.WriteString(fmt.Sprintf("\nCache size: %d", m.client.CacheSize()))

	tv := lipgloss.NewStyle().
		Align(lipgloss.Center).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.BrightBlue)

	width, height, _ := term.GetSize(int(os.Stdout.Fd()))
	screen := lipgloss.NewStyle().
		Width(width).
		Height(height).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Padding(2, 0)

	return screen.Render(tv.Render(sb.String()))
}

func (m Model) nextIndex() int {
	if m.activeIndex < len(m.page.SubPages)-1 {
		return m.activeIndex + 1
	}
	return m.activeIndex
}

func (m Model) prevIndex() int {
	if m.activeIndex > 0 {
		return m.activeIndex - 1
	}
	return m.activeIndex
}

type newPage svttext.Page

func (m Model) getPage(pageId string) tea.Cmd {
	return func() tea.Msg {
		page, err := m.client.GetPage(pageId)
		if err != nil {
			// TODO: add error / debug handling
			// log.Printf("failed to get page: %v", err)
			return nil
		}
		return newPage(page)
	}
}

func (m Model) preLoadPage(pageId string) tea.Cmd {
	return func() tea.Msg {
		_, _ = m.client.GetPage(pageId)
		return nil
	}
}
