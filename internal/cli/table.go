package cli

import (
	"fmt"
	"strings"
)

const (
	cellPadWidth = 2
)

// Table represents a console-printed text table.
type Table struct {
	grid         [][]string
	columnWidths map[int]int
}

// NewTable creates a new table instance with an empty grid.
func NewTable() *Table {
	return &Table{
		columnWidths: make(map[int]int),
	}
}

// Add statefully adds a row to the grid.
func (t *Table) Add(row []string) error {
	if len(t.grid) > 0 && len(t.grid[0]) != len(row) {
		return fmt.Errorf("table: column quantity in row inconsistent with existing grid")
	}

	t.grid = append(t.grid, row)

	for idx, column := range row {
		if width, ok := t.columnWidths[idx]; !ok || width < len(column) {
			t.columnWidths[idx] = len(column)
		}
	}

	return nil
}

// IsEmpty indicates whether the table is not populated with any rows.
func (t *Table) IsEmpty() bool {
	return len(t.grid) == 0
}

// String returns a string representation of the grid, with some padding between columns.
func (t *Table) String() string {
	var table []string

	for _, row := range t.grid {
		var serializedRow []string

		for idx, column := range row {
			if idx == len(row)-1 || t.columnWidths[idx] == 0 {
				serializedRow = append(serializedRow, column)
			} else {
				serializedRow = append(
					serializedRow,
					Pad(column, t.columnWidths[idx]+cellPadWidth),
				)
			}
		}

		table = append(table, strings.Join(serializedRow, ""))
	}

	return strings.Join(table, "\n")
}
