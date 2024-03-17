package tuttableevertras

// https://github.com/Evertras/bubble-table

import (
	"context"
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
	stream := make(chan []*Row)
	defer close(stream)
	ctx := context.Background()
	p := tea.NewProgram(NewModel(ctx, stream))
	_, err = p.Run()
	if err != nil {
		log.Fatal(err)
	}
}

type Model struct {
	Table           table.Model
	SortKey         string
	CurrSortKey     string //To be able to start with ASC for any new selected column
	SortDirection   string
	FilterTextInput textinput.Model
	// Streams
	Rows        *Rows
	Ctx         context.Context
	RowsStream  dataStream
	SelectedRow Row
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

func NewModel(ctx context.Context, stream chan []*Row) Model {
	columns := []table.Column{
		table.NewColumn(keyID, "(I)D", 5).WithFiltered(true),
		table.NewColumn(keyName, "(N)ame", 10).WithFiltered(true),
		table.NewColumn(keyDescription, "(D)escription", 30).WithFiltered(true),
		table.NewColumn(keyCount, "(M)oney", 10).WithStyle(lipgloss.NewStyle().Align(lipgloss.Right)).WithFiltered(true),
	}

	// Start with the default key map and change it slightly
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
		Ctx:             ctx,
		Rows:            &Rows{},
		RowsStream:      stream,
		SelectedRow:     Row{},
	}

	model.updateFooter()

	return model
}

type dataStream chan []*Row
type dataStreamMsg dataStream

// Connect to stream
func (m Model) makeConnectCmd() tea.Cmd {
	return func() tea.Msg {
		go m.Rows.generateData(m.Ctx, m.RowsStream)
		return dataStreamMsg(m.RowsStream)
	}
}

// Listen to stream
func listenCmd(m Model) tea.Cmd {
	return func() tea.Msg {
		return dataStreamMsg(m.RowsStream)
	}
}

func (m Model) Init() tea.Cmd {
	return m.makeConnectCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	doSort := false

	switch msg := msg.(type) {

	case dataStreamMsg:
		rows := []table.Row{}
		select {
		case <-m.Ctx.Done():
			return m, tea.Quit
		case rs := <-m.RowsStream:
			for _, r := range rs {
				rows = append(rows, MakeTableRow(r.id, r.name, r.description, r.count))
			}
			m.updateFooter()
			m.Table = m.Table.WithRows(rows)
		}
		return m, listenCmd(m) // listen for next event

	case tea.KeyMsg:
		if m.FilterTextInput.Focused() {
			if msg.String() == "esc" {
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
		case "enter":
			selRow := m.Table.HighlightedRow().Data[keyMeta].(Row)
			log.Printf("SelectedRow: %+v\n", selRow)
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

	// for _, row := range m.Table.SelectedRows() {
	// for _, row := range m.SelectedRows {
	// selectedIDs = append(selectedIDs, row.Data[keyID].(string))
	// }

	body.WriteString(m.Table.View() + "\n")
	body.WriteString(m.FilterTextInput.View() + "\n")
	body.WriteString(fmt.Sprintf("\nSelectedIDs: %s\n\n", strings.Join(selectedIDs, ", ")))

	return body.String()
}

//--------------------------------------------------------------------------------
// Data

type Row struct {
	id          string
	name        string
	description string
	count       float64
}

type Rows []*Row

func (r *Rows) generateData(ctx context.Context, stream dataStream) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if len(*r) == 0 {
				*r = Rows{
					{"100", "Carl", "The head of the organization", 9000},
					{"123", "Abe", "Good with all weapons", 170},
					{"398", "Dave", "Drives and flies any vehicle or plane", -170},
					{"093", "Eve", "Cooks the most effective explosives", -200},
					{"007", "Fiona", "Can open any lock", 445},
				}
			} else {
				sign := []float64{1, -1}
				randSignInt, _ := rand.Int(rand.Reader, big.NewInt(2))
				randSign := sign[randSignInt.Uint64()]

				amount := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
				randAmountInt, _ := rand.Int(rand.Reader, big.NewInt(int64(len(amount))))
				randAmount := amount[randAmountInt.Uint64()]

				nofRows := int64(len(*r))
				randRow, _ := rand.Int(rand.Reader, big.NewInt(nofRows))

				(*r)[randRow.Uint64()].count += randAmount * randSign
			}

			stream <- *r
			time.Sleep(25 * time.Millisecond)
		}
	}
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
		totalCount(*m.Rows),
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
