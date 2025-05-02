package tui

import (
	"slices"
	"strings"
	"time"

	"github.com/aKjeller/text-tv/pkg/svttext"
	"github.com/charmbracelet/bubbles/v2/help"
	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type keyMap struct {
	Next     key.Binding
	Prev     key.Binding
	Right    key.Binding
	Left     key.Binding
	Continue key.Binding
	Number   key.Binding
	Help     key.Binding
	Quit     key.Binding
}

var defaultKeyMap = keyMap{
	Next: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k/↑", "next page")),
	Prev: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j/↓", "previous page")),
	Right: key.NewBinding(
		key.WithKeys("l", "right"),
		key.WithHelp("l/→", "scroll right")),
	Left: key.NewBinding(
		key.WithKeys("h", "left"),
		key.WithHelp("h/←", "scroll left")),
	Continue: key.NewBinding(
		key.WithKeys("space"),
		key.WithHelp("space", "continue")),
	Number: key.NewBinding(
		key.WithKeys("0", "1", "2", "3", "4", "5", "6", "7", "8", "9"),
		key.WithHelp("0-9", "search")),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help")),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("q", "quit")),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Left, k.Prev, k.Next, k.Right},
		{k.Number, k.Continue, k.Help, k.Quit},
	}
}

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

	// keybinds
	keys keyMap
	help help.Model

	client *svttext.Client
}

func NewModel() Model {
	return Model{
		client:      svttext.NewClient(svttext.WithCacheTime(1 * time.Minute)),
		page:        svttext.Page{},
		activeDot:   lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render("●"),
		inactiveDot: lipgloss.NewStyle().Foreground(lipgloss.Color("238")).Render("●"),
		keys:        defaultKeyMap,
		help:        help.New(),
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
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Next):
			return m, m.getPage(m.page.NextPage)
		case key.Matches(msg, m.keys.Prev):
			return m, m.getPage(m.page.PrevPage)
		case key.Matches(msg, m.keys.Left):
			m.activeIndex = m.prevIndex()
		case key.Matches(msg, m.keys.Right):
			m.activeIndex = m.nextIndex()
		case key.Matches(msg, m.keys.Continue):
			newIndex := m.nextIndex()
			if newIndex != m.activeIndex {
				m.activeIndex = newIndex
			} else {
				return m, m.getPage(m.page.NextPage)
			}
		case key.Matches(msg, m.keys.Number):
			m.input = append(m.input, msg.Key().Code)
			if len(m.input) == 3 {
				return m, m.getPage(string(m.input))
			}
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
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
			sb.WriteString(replacePageNbr(row, m.input).ColorString() + "\n")
		} else if i < svttext.PageHeight-1 {
			sb.WriteString(row.ColorString() + "\n")
		}
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

	sb.WriteString("\n" + m.help.View(m.keys))

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
