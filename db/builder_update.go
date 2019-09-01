package db

import (
	"database/sql"

	"github.com/valyala/bytebufferpool"
)

// UpdateBuilder class
type UpdateBuilder struct {
	table string
	where string
	args  []interface{}
	builder
}

// NewUpdate func
func NewUpdate(t string) *UpdateBuilder {
	return &UpdateBuilder{table: t}
}

// Table string
func (u *UpdateBuilder) Table(t string) *UpdateBuilder {
	u.table = t
	return u
}

// Parse sql
func (u *UpdateBuilder) Parse(row MapRow) (string, []interface{}) {
	// keys, values, _ := keyValueList(ActionUpdate, data)
	keys, values := sqlUpdateParamsByMapRow(row)

	// buf := bytePool.Get()
	buf := bytebufferpool.Get()

	// s := bytes.Buffer{}
	driver := u.GetDatabase().Driver
	buf.Write(bSQLUpdate)
	buf.Write(driver.QuoteIdentifierB(u.table))
	buf.Write(bSQLSet)
	buf.Write(keys)
	if u.where != "" {
		buf.Write(bSQLWhere)
		buf.Write(driver.Adapt(u.where, len(values)))
		values = append(values, u.args...)
	}
	str := buf.String()
	// bytePool.Put(buf)
	bytebufferpool.Put(buf)

	return str, values
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
func (u *UpdateBuilder) Update(row MapRow) (sql.Result, error) {
	sql, vals := u.Parse(row)
	return u.GetDatabase().ExecPrepare(sql, vals...)
}

// // UpdateByString by string
// func (u *UpdateBuilder) UpdateByString(str string) (sql.Result, error) {
// 	var values []interface{}
// 	s := bytes.Buffer{}
// 	driver := u.GetDatabase().Driver
// 	s.Write(bSQLUpdate)
// 	s.WriteString(driver.QuoteField(u.table))
// 	s.Write(bSQLSet)
// 	s.WriteString(str)
// 	if u.where != "" {
// 		s.Write(bSQLWhere)
// 		s.WriteString(u.where)
// 		values = u.args
// 	}
//
// 	return u.GetDatabase().ExecPrepare(s.String(), values...)
// }

// // TxUpdateByString trunsaction
// func (u *UpdateBuilder) TxUpdateByString(tx *Tx, str string) (sql.Result, error) {
// 	var values []interface{}
// 	s := bytes.Buffer{}
// 	driver := u.GetDatabase().Driver
// 	s.Write(bSQLUpdate)
// 	s.WriteString(driver.QuoteField(u.table))
// 	s.Write(bSQLSet)
// 	s.WriteString(str)
// 	if u.where != "" {
// 		s.Write(bSQLWhere)
// 		s.WriteString(u.where)
// 		values = u.args
// 	}
//
// 	return tx.ExecPrepare(s.String(), values...)
// }

// TxUpdate trunsaction
func (u *UpdateBuilder) TxUpdate(tx *Tx, row MapRow) (sql.Result, error) {
	sql, vals := u.Parse(row)
	return tx.ExecPrepare(sql, vals...)
}
