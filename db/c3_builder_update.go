package db

import (
	"database/sql"
	"fmt"
	"io"

	"github.com/kere/gno/libs/util"
	"github.com/valyala/bytebufferpool"
)

// UpdateBuilder class
type UpdateBuilder struct {
	where string
	args  []interface{}
	Builder
}

// NewUpdate func
func NewUpdate(t string) UpdateBuilder {
	u := UpdateBuilder{}
	u.table = t
	return u
}

// Table string
func (u *UpdateBuilder) Table(t string) *UpdateBuilder {
	u.table = t
	return u
}

// SetPrepare prepare sql
func (u *UpdateBuilder) Prepare(v bool) *UpdateBuilder {
	u.isPrepare = v
	return u
}

// Where sql
func (u *UpdateBuilder) Where(cond string, args ...interface{}) *UpdateBuilder {
	if cond == "" {
		return u
	}
	u.where = cond
	u.args = args
	return u
}

// Update db
func (u *UpdateBuilder) Update(fields []string, row []interface{}) (sql.Result, error) {
	sqlstr, vals := u.ParseP(fields, row)
	if u.isPrepare {
		return u.ExecPrepare(sqlstr, vals)
	}
	return u.Exec(sqlstr, vals)
}

// ParseP sql
func (u *UpdateBuilder) ParseP(fields []string, row []interface{}) (string, []interface{}) {
	buf := bytebufferpool.Get()
	count := len(u.args)
	values := GetRow()
	if u.where != "" && count != 0 {
		values = append(values, u.args...)
	}
	values = append(values, row...)

	driver := u.GetDatabase().Driver
	buf.Write(bSQLUpdate)
	driver.WriteQuoteIdentifier(buf, u.table)
	buf.Write(bSQLSet)
	writeUpdate(buf, driver, fields, row, len(u.args))

	if u.where != "" {
		buf.Write(bSQLWhere)
		buf.WriteString(u.where)
	}
	str := buf.String()
	bytebufferpool.Put(buf)

	return str, values
}

func writeUpdateItem(w io.Writer, driver IDriver, field string, value interface{}, seq int) int {
	driver.WriteQuoteIdentifier(w, field)
	w.Write(util.BEqual)
	if value == nil {
		w.Write(util.BNull)
		return seq
	}
	w.Write(util.BDoller)
	w.Write(util.Str2Bytes(fmt.Sprint(seq)))
	seq++
	return seq
}

func writeUpdate(w io.Writer, driver IDriver, fields []string, values []interface{}, count int) int {
	n := len(fields)
	seq := 1 + count
	seq = writeUpdateItem(w, driver, fields[0], values[0], seq)
	for i := 1; i < n; i++ {
		w.Write(util.BComma)
		seq = writeUpdateItem(w, driver, fields[i], values[i], seq)
	}

	return seq
}
