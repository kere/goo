package dba

import (
	"database/sql"

	"github.com/valyala/bytebufferpool"
)

// DeleteBuilder class
type DeleteBuilder struct {
	Builder
	IsPrepare bool
}

// NewDelete func
func newDelete(t string) DeleteBuilder {
	d := DeleteBuilder{IsPrepare: false}
	d.table = t
	return d
}

// Delete delete
func (d *DeleteBuilder) Delete(where string, args ...interface{}) (sql.Result, error) {
	if d.IsPrepare {
		return d.ExecPrepare(parseDelete(d, where), args)
	}
	return d.Exec(parseDelete(d, where), args)
}

func parseDelete(d *DeleteBuilder, where string) string {
	buf := bytebufferpool.Get()
	buf.Write(bSQLDelete)
	buf.Write(bSQLFrom)
	driver := d.database.Driver
	driver.WriteQuoteIdentifier(buf, d.table)
	if where != "" {
		buf.Write(bSQLWhere)
		buf.WriteString(where)
	}
	str := buf.String()
	bytebufferpool.Put(buf)
	return str
}
