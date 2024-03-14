package tuttableevertras

// https://github.com/Evertras/bubble-table

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	keyID          = "id"
	keyName        = "name"
	keyDescription = "description"
	keyCount       = "count"
	keyMeta        = "meta"
)

var (
	customBorder = table.Border{
		Top:    "─",
		Right:  "│",
		Bottom: "─",
		Left:   "│",

		TopRight:    "┐",
		BottomRight: "┘",
		BottomLeft:  "└",
		TopLeft:     "┌",

		TopJunction:    "┬",
		RightJunction:  "┤",
		BottomJunction: "┴",
		LeftJunction:   "├",
		InnerJunction:  "┼",
		InnerDivider:   "│",
	}
)

func Run() {
	os.Truncate("debug.log", 0)
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	//
	p := tea.NewProgram(NewModel())
	_, err = p.Run()
	if err != nil {
		log.Fatal(err)
	}
}

type Model struct {
	Table           table.Model
	Rows            []*Row
	SortKey         string
	CurrSortKey     string //To be able to start with ASC for any new selected column.
	SortDirection   string
	FilterTextInput textinput.Model
}

type Row struct {
	id          string
	name        string
	description string
	count       float64
}

func MakeTableRow(id, name, description string, count float64) table.Row {
	countVal := fmt.Sprintf("% 6.2f", count)
	countStyled := table.NewStyledCell(countVal, lipgloss.NewStyle().Foreground(lipgloss.Color("#8f8")).Align(lipgloss.Right))
	if count < 0 {
		countStyled = table.NewStyledCell(countVal, lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Align(lipgloss.Right))
	}

	return table.NewRow(table.RowData{
		keyID:          id,
		keyName:        name,
		keyDescription: description,
		keyCount:       countStyled,
		keyMeta: Row{
			id:          id,
			name:        name,
			description: description,
			count:       count,
		},
	})
}

func NewModel() Model {
	columns := []table.Column{
		table.NewColumn(keyID, "(I)D", 5).WithFiltered(true),
		table.NewColumn(keyName, "(N)ame", 10).WithFiltered(true),
		table.NewColumn(keyDescription, "(D)escription", 30).WithFiltered(true),
		table.NewColumn(keyCount, "(M)oney", 10).WithStyle(lipgloss.NewStyle().Align(lipgloss.Right)).WithFiltered(true),
	}

	// rows := []table.Row{}
	// for _, r := range items {
	// 	rows = append(rows, MakeTableRow(r.id, r.name, r.description, r.count))
	// }

	// Start with the default key map and change it slightly, just for demoing
	keys := table.DefaultKeyMap()
	keys.RowDown.SetKeys("j", "down", "s")
	keys.RowUp.SetKeys("k", "up", "w")

	model := Model{
		Table: table.
			New(columns).
			Filtered(true).
			Focused(true).
			WithPageSize(60).
			HeaderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#0000aa")).Bold(true)).
			SelectableRows(true).
			Border(customBorder).
			WithKeyMap(keys).
			WithStaticFooter("Footer!").
			WithSelectedText(" ", "✓").
			WithBaseStyle(
				lipgloss.NewStyle().
					BorderForeground(lipgloss.Color("#a38")).
					Foreground(lipgloss.Color("#a7a")).
					Align(lipgloss.Left),
			).
			SortByAsc(keyID),
		CurrSortKey:     keyID,
		SortKey:         keyID,
		SortDirection:   "asc",
		FilterTextInput: textinput.New(),
	}

	model.updateFooter()

	return model
}

func (m Model) Init() tea.Cmd {
	return initDataCmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	doSort := false

	switch msg := msg.(type) {
	case []*Row: //New data has presented itself
		m.Rows = msg
		rows := []table.Row{}
		for _, r := range m.Rows {
			rows = append(rows, MakeTableRow(r.id, r.name, r.description, r.count))
		}
		m.updateFooter()
		m.Table = m.Table.WithRows(rows)
		cmds = append(cmds, m.updateDataCmd)
	case tea.KeyMsg:
		if m.FilterTextInput.Focused() {
			if msg.String() == "enter" || msg.String() == "esc" {
				m.FilterTextInput.Blur()
			} else {
				m.FilterTextInput, _ = m.FilterTextInput.Update(msg)
			}
			m.Table = m.Table.WithFilterInput(m.FilterTextInput)
		}
		//
		switch msg.String() {
		case "ctrl+c", "q":
			cmds = append(cmds, tea.Quit)
		case "h":
			m.Table = m.Table.WithHeaderVisibility(!m.Table.GetHeaderVisibility())
		case "/":
			m.FilterTextInput.Focus()
		case "i":
			m.SortKey = keyID
			doSort = true
		case "n":
			m.SortKey = keyName
			doSort = true
		case "d":
			m.SortKey = keyDescription
			doSort = true
		case "m":
			m.SortKey = keyCount
			doSort = true
		default:
			m.Table, cmd = m.Table.Update(msg)
			cmds = append(cmds, cmd)
			m.updateFooter()
		}
		if doSort {
			SortTable(&m)
			doSort = false
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	body := strings.Builder{}

	body.WriteString("Press left/right or page up/down to move pages\n")
	body.WriteString("Press 'i' to toggle the header visibility\n")
	body.WriteString("Press space/enter to select a row, q or ctrl+c to quit\n")
	body.WriteString("A filtered simple default table\n" +
		"Currently filter by Title and Author, press / + letters to start filtering, and escape to clear filter.\n" +
		"Press q or ctrl+c to quit\n\n")

	selectedIDs := []string{}

	for _, row := range m.Table.SelectedRows() {
		selectedIDs = append(selectedIDs, row.Data[keyID].(string))
	}

	body.WriteString(m.Table.View() + "\n")
	body.WriteString(m.FilterTextInput.View() + "\n")
	body.WriteString(fmt.Sprintf("\nSelectedIDs: %s\n\n", strings.Join(selectedIDs, ", ")))

	return body.String()
}

//--------------------------------------------------------------------------------
// Data

func initDataCmd() tea.Msg {
	return []*Row{
		{"100", "Carl", "The head of the organization", 9000},
		{"123", "Abe", "Good with all weapons", 170},
		{"398", "Dave", "Drives and flies any vehicle or plane", -170},
		{"093", "Eve", "Cooks the most effective explosives", -200},
		{"007", "Fiona", "Can open any lock", 445},
	}
}

func (m *Model) updateDataCmd() tea.Msg {
	sign := []float64{1, -1}
	randSignInt, _ := rand.Int(rand.Reader, big.NewInt(2))
	randSign := sign[randSignInt.Uint64()]

	amount := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	randAmountInt, _ := rand.Int(rand.Reader, big.NewInt(int64(len(amount))))
	randAmount := amount[randAmountInt.Uint64()]

	nofRows := int64(len(m.Rows))
	randRow, _ := rand.Int(rand.Reader, big.NewInt(nofRows))

	m.Rows[randRow.Uint64()].count += randAmount * randSign

	time.Sleep(100 * time.Millisecond)
	return m.Rows
}

//--------------------------------------------------------------------------------
// Helper functions

func SortTable(m *Model) {
	if m.SortKey == m.CurrSortKey {
		if m.SortDirection == "asc" {
			m.Table = m.Table.SortByAsc(m.SortKey)
			m.SortDirection = "desc"
		} else {
			m.Table = m.Table.SortByDesc(m.SortKey)
			m.SortDirection = "asc"
		}
	} else {
		m.Table = m.Table.SortByAsc(m.SortKey)
		m.SortDirection = "desc"
	}
	m.CurrSortKey = m.SortKey
}

func (m *Model) updateFooter() {
	footerText := fmt.Sprintf(
		"%d/%d    total: %.2f",
		m.Table.CurrentPage(),
		m.Table.MaxPages(),
		totalCount(m.Rows),
	)

	m.Table = m.Table.WithStaticFooter(footerText)
}

func totalCount(rows []*Row) float64 {
	total := 0.0
	for _, r := range rows {
		total += r.count
	}
	return total
}
