package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/kere/gno/db/drivers"
	"github.com/kere/gno/libs/conf"
	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/libs/myerr"
)

var (
	dbpool *databasePool
	// dbConf conf.Conf
)

func init() {
	dbpool = &databasePool{dblist: make(map[string]*Database)}
}

type databasePool struct {
	dblist  map[string]*Database
	current *Database
}

// Current database
func (dp *databasePool) Current() *Database {
	return dp.current
}

// SetCurrent database
func (dp *databasePool) SetCurrent(d *Database) {
	dp.current = d
}

// Use database
func (dp *databasePool) Use(name string) {
	c := dp.GetDatabase(name)
	if c == nil {
		fmt.Println(name, " database is not found!")
		return
	}
	dp.current = c
}

// SetDatabase by name
func (dp *databasePool) SetDatabase(name string, d *Database) {
	dp.dblist[name] = d
}

// GetDatabase by name
func (dp *databasePool) GetDatabase(name string) *Database {
	if v, ok := dp.dblist[name]; ok {
		return v
	}
	return nil
}

// Init it
func Init(name string, config map[string]string) {
	fmt.Println("Init Database", config)
	dbpool.SetCurrent(New(name, config))
}

func confGet(config map[string]string, key string) string {
	if v, ok := config[key]; ok {
		return v
	}
	return ""
}

// New func
// create a database instance
func New(name string, c map[string]string) *Database {
	if dbpool.GetDatabase(name) != nil {
		panic(name + " this database is already exists!")
	}

	if c == nil {
		return nil
	}

	driverName := confGet(c, "driver")
	logger := NewLogger(c)

	var driver IDriver
	switch driverName {
	case "postgres", "psql":
		driver = &drivers.Postgres{DBName: confGet(c, "dbname"),
			User:     confGet(c, "user"),
			Password: confGet(c, "password"),
			Host:     confGet(c, "host"),
			HostAddr: confGet(c, "hostaddr"),
			Port:     confGet(c, "port"),
		}

	case "mysql":
		driver = &drivers.Mysql{DBName: confGet(c, "dbname"),
			User:       confGet(c, "user"),
			Password:   confGet(c, "password"),
			Protocol:   confGet(c, "protocol"),
			Parameters: confGet(c, "parameters"),
			Addr:       confGet(c, "addr")}

	case "sqlite3":
		driver = &drivers.Sqlite3{File: confGet(c, "file")}

	default:
		logger.Println("you may need regist a custom driver: db.RegistDriver(Mysql{})")
		driver = &drivers.Common{}

	}

	driver.SetConnectString(confGet(c, "connect"))
	// poolSize, err := strconv.Atoi(confGet(c, "pool_size"))
	// if poolSize == 0 || err != nil {
	// 	poolSize = 3
	// }
	// maxCount, err := strconv.Atoi(confGet(c, "max_count"))
	// if maxCount == 0 || err != nil {
	// 	maxCount = 10
	// }

	d := NewDatabase(name, driver, conf.Conf(c), logger)

	// ------- time zone --------
	// if confGet(conf, "timezone") != "" {
	// 	loc, err := time.LoadLocation(confGet(conf, "timezone"))
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	d.Location = loc
	// }

	dbpool.SetDatabase(name, d)
	dbpool.SetCurrent(d)
	return d
}

// CQuery from database on prepare mode
// This function use the current database from database bool
// You can set another database use Use(i) or create an new database use New(name, conf)
func CQuery(conn *sql.DB, sqlstr string, args ...interface{}) (DataSet, error) {
	var rows *sql.Rows
	var err error
	if len(args) == 0 {
		rows, err = conn.Query(sqlstr)
	} else {
		rows, err = conn.Query(sqlstr, args...)
	}

	if err != nil {
		logSQLErr(sqlstr, args)
		return nil, myerr.New(err).Log().Stack()
	}
	dataset, err := ScanRows(rows)
	if err != nil {
		return nil, myerr.New(err).Log().Stack()
	}

	return dataset, nil
	// return Current().Query(NewSqlState([]byte(sqlstr), args))
}

func logSQLErr(sqlstr string, args []interface{}) {
	sep := ": "
	var s strings.Builder
	s.WriteString(sqlstr)
	s.WriteString(SLineBreak)
	l := len(args)
	for i := 0; i < l; i++ {
		s.WriteString(fmt.Sprint(i, sep))

		switch args[i].(type) {
		case []byte:
			s.Write(args[i].([]byte))
		default:
			s.WriteString(fmt.Sprint(args[i]))
		}
		s.WriteString(SLineBreak)
	}
	s.WriteString(SLineBreak)
	log.App.Error(s.String())
}

// CQueryPrepare from database on prepare mode
// In prepare mode, the sql command will be cached by database
// This function use the current database from database bool
// You can set another database by Use(i) or New(name, conf) an new database
func CQueryPrepare(conn *sql.DB, sqlstr string, args ...interface{}) (DataSet, error) {
	var rows *sql.Rows
	var err error

	s, err := conn.Prepare(sqlstr)
	if err != nil {
		logSQLErr(sqlstr, args)
		return nil, myerr.New(err).Log().Stack()
	}

	defer s.Close()

	if len(args) == 0 {
		rows, err = s.Query()
	} else {
		rows, err = s.Query(args...)
	}
	if err != nil {
		logSQLErr(sqlstr, args)
		return nil, myerr.New(err).Log().Stack()
	}

	dataset, err := ScanRows(rows)
	if err != nil {
		return nil, myerr.New(err).Log().Stack()
	}

	return dataset, nil
}

// func Find(conn *sql.DB, cls IVO, sqlstr []byte, args ...interface{}) (VODataSet, error) {
// 	dataset, err := Query(conn, sqlstr, args...)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return NewStructConvert(cls).DataSet2Struct(dataset)
// 	// return Current().Find(cls, NewSqlState([]byte(sqlstr), args))
// }
// func FindPrepare(conn *sql.DB, cls IVO, sqlstr []byte, args ...interface{}) (VODataSet, error) {
// 	dataset, err := QueryPrepare(conn, sqlstr, args...)
//
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return NewStructConvert(cls).DataSet2Struct(dataset)
// }

// // ExecFromFile sql from a file
// // This function run sql under not transaction mode and use the current database from database bool
// func ExecFromFile(file string) error {
// 	var filebytes []byte
// 	var err error
//
// 	if filebytes, err = ioutil.ReadFile(file); err != nil {
// 		return myerr.New(err).Log().Stack()
// 	}
// 	conn := Current().Connection.Connect()
// 	b := bytes.Split(filebytes, []byte(";"))
// 	for _, i := range b {
// 		if len(bytes.TrimSpace(i)) == 0 {
// 			continue
// 		}
//
// 		_, err = Exec(conn, i)
// 		if err != nil {
// 			return myerr.New(err).Log().Stack()
// 		}
// 	}
// 	return nil
// }

// Exec sql.
// If your has more than on sql command, it will only excute the first.
// This function use the current database from database bool
func Exec(conn *sql.DB, sqlstr string, args ...interface{}) (result sql.Result, err error) {
	if len(args) == 0 {
		result, err = conn.Exec(sqlstr)
	} else {
		result, err = conn.Exec(sqlstr, args...)
	}
	if err != nil {
		logSQLErr(sqlstr, args)
		return result, myerr.New(err).Log().Stack()
	}
	return result, nil
}

// ExecPrepare sql on prepare mode
// This function use the current database from database bool
func ExecPrepare(conn *sql.DB, sqlstr string, args ...interface{}) (result sql.Result, err error) {
	s, err := conn.Prepare(sqlstr)
	if err != nil {
		return nil, myerr.New(err).Log().Stack()
	}

	defer s.Close()

	if len(args) == 0 {
		result, err = s.Exec()
	} else {
		result, err = s.Exec(args...)
	}
	if err != nil {
		logSQLErr(sqlstr, args)
		return result, myerr.New(err).Log().Stack()
	}
	return result, nil
}

// Get a database instance by name from database pool
func Get(name string) *Database {
	return dbpool.GetDatabase(name)
}

//Current Return the current database from database pool
func Current() *Database {
	if dbpool.Current() == nil {
		panic("db is not initalized")
	}
	return dbpool.Current()
}

// Use current database by name
func Use(name string) {
	dbpool.Use(name)
}

// DatabaseCount Get database count
func DatabaseCount() int {
	return len(dbpool.dblist)
}
