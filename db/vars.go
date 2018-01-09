package db

var (
	ColumnBytePrefix = "byte_"
	bSQLSelect       = []byte("select ")
	bSQLDelete       = []byte("delete ")
	bSQLUpdate       = []byte("update ")
	bSQLSet          = []byte(" set ")
	bSQLFrom         = []byte(" from ")
	bSQLWhere        = []byte(" where ")
	bSQLOrder        = []byte(" order by ")
	bSQLLimit        = []byte(" limit ")
	bSQLOffset       = []byte(" offset ")

	// DBTimeFormat = time.RFC3339
	DBTimeFormat   = "2006-01-02 15:04:05"
	B_QuestionMark = []byte("?")
	B_Equal        = []byte("=")

	B_CommaSplit  = []byte(",")
	B_EmptyString = []byte("")
	B_StarKey     = []byte("*")
	B_Space       = []byte(" ")
)
