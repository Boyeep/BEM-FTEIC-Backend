package database

import (
	"fmt"
	"io"

	"gorm.io/gorm"
)

type columnInfo struct {
	TableName  string `gorm:"column:table_name"`
	ColumnName string `gorm:"column:column_name"`
	DataType   string `gorm:"column:data_type"`
	IsNullable string `gorm:"column:is_nullable"`
}

// PrintSchema writes a read-only inventory of application tables. It is
// intended for deployment diagnostics before applying migrations.
func PrintSchema(db *gorm.DB, output io.Writer) error {
	var columns []columnInfo
	err := db.Raw(`
		SELECT table_name, column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_schema = 'public'
		  AND table_name IN ('admins', 'blogs', 'events', 'galeri', 'galleries', 'profiles')
		ORDER BY table_name, ordinal_position
	`).Scan(&columns).Error
	if err != nil {
		return fmt.Errorf("inspect public schema: %w", err)
	}

	if len(columns) == 0 {
		_, _ = fmt.Fprintln(output, "No application tables found in public schema.")
		return nil
	}

	currentTable := ""
	for _, column := range columns {
		if column.TableName != currentTable {
			currentTable = column.TableName
			_, _ = fmt.Fprintf(output, "\n[%s]\n", currentTable)
		}
		_, _ = fmt.Fprintf(
			output,
			"- %s: %s (nullable=%s)\n",
			column.ColumnName,
			column.DataType,
			column.IsNullable,
		)
	}
	return nil
}
