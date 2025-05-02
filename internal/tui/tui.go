package tui

import (
	"slices"
	"strings"
	"time"

	"github.com/aKjeller/text-tv/pkg/svttext"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type Model struct {
	page svttext.Page

	// SubPage Pagination
	activeIndex int
	activeDot   string
	inactiveDot string

	// terminal size
	width  int
	height int

	input []rune

	client *svttext.Client
}

func NewModel() Model {
	return Model{
		client:      svttext.NewClient(svttext.WithCacheTime(1 * time.Minute)),
		page:        svttext.Page{},
		activeIndex: 0,
		activeDot:   lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render("●"),
		inactiveDot: lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Render("●"),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.getPage("100"),
		m.preLoadPage("130"),
		m.preLoadPage("377"),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
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
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			m.input = append(m.input, msg.Key().Code)
			if len(m.input) == 3 {
				return m, m.getPage(string(m.input))
			}
		}
	case pageLoadedMsg:
		m.page = msg.page
		m.activeIndex = 0
		m.input = []rune{}
		return m, tea.Batch(m.preLoadPage(m.page.PrevPage), m.preLoadPage(m.page.NextPage))
	case pageFallbackMsg:
		return m, m.getPage(msg.fallbackPageId)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

func (m Model) View() string {
	if len(m.page.SubPages) == 0 {
		return ""
	}

	var sb strings.Builder
	for i, row := range m.page.SubPages[m.activeIndex].Grid {
		if i == 0 && len(m.input) != 0 {
			sb.WriteString(replacePageNbr(row, m.input).ColorString())
		} else {
			sb.WriteString(row.ColorString())
		}
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

	screen := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Padding(2, 0)

	return screen.Render(tv.Render(sb.String()))
}

func replacePageNbr(row svttext.Row, nbr []rune) svttext.Row {
	res := slices.Clone(row)
	for i := 0; i < 3; i++ {
		pos := 3 + i
		if i < len(nbr) && nbr[i] != 0 {
			res[pos].Char = nbr[i]
		} else {
			res[pos].Char = ' '
		}
	}
	return res
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

type pageLoadedMsg struct {
	page svttext.Page
}

type pageFallbackMsg struct {
	fallbackPageId string
}

func (m Model) getPage(pageId string) tea.Cmd {
	return func() tea.Msg {
		page, err := m.client.GetPage(pageId)
		if err != nil {
			// TODO: add error / debug handling
			// log.Printf("failed to get page: %v", err)
			return pageFallbackMsg{"100"}
		}

		if len(page.SubPages) == 0 {
			if page.PrevPage != "" {
				return pageFallbackMsg{page.PrevPage}
			} else if page.NextPage != "" {
				return pageFallbackMsg{page.NextPage}
			} else {
				return nil
			}
		}

		return pageLoadedMsg{page}
	}
}

func (m Model) preLoadPage(pageId string) tea.Cmd {
	return func() tea.Msg {
		_, _ = m.client.GetPage(pageId)
		return nil
	}
}
