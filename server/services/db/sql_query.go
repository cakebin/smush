package db

import (
  "bytes"
  "fmt"
)


// MakeMultiInsertStatement generates a PostgreSQL insert statement for a given
// table, set of columns, and the corresponding values; use to easily make a
// statement for inserting multiple rows into a the given table
func MakeMultiInsertStatement(table string, columns []string, numInserts int, returningCol string) string {
  buf := bytes.Buffer{}

  buf.WriteString(fmt.Sprintf("INSERT INTO %s ", table))
  buf.WriteString("(")
  for i, col := range columns {
    buf.WriteString(fmt.Sprintf("%s", col))
    if i < (len(columns) - 1) {
      buf.WriteString(", ")
    }
  }
  buf.WriteString(")")
  buf.WriteString(" VALUES ")

  curValNum := 1

  for j := 1; j <= numInserts; j++ {
    buf.WriteString("(")
    for k := range columns {
      buf.WriteString(fmt.Sprintf("$%d", curValNum))
      curValNum++
      if k < (len(columns) - 1) {
        buf.WriteString(", ")
      }
    }
    buf.WriteString(")")
    if j < numInserts {
      buf.WriteString(", ")
    }
  }

  buf.WriteString(fmt.Sprintf(" RETURNING %s", returningCol))

  return buf.String()
}
