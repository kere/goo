package db

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/kere/gno/db/drivers"
	"github.com/kere/gno/libs/conf"
)

var (
	dbpool *databasePool
	dbConf conf.Conf
)

type databasePool struct {
	dblist  map[string]*Database
	current *Database
}

func (this *databasePool) Current() *Database {
	return this.current
}

func (this *databasePool) SetCurrent(d *Database) {
	this.current = d
}

func (this *databasePool) Use(name string) {
	c := this.GetDatabase(name)
	if c == nil {
		fmt.Println(name, " database is not found!")
		return
	} else {
		this.current = c
	}
}

func (this *databasePool) SetDatabase(name string, d *Database) {
	this.dblist[name] = d
}

func (this *databasePool) GetDatabase(name string) *Database {
	if v, ok := this.dblist[name]; ok {
		return v
	}
	return nil
}

func init() {
	dbpool = &databasePool{dblist: make(map[string]*Database)}
}

// Init it
func Init(name string, config map[string]string) {
	dbConf = conf.Conf(config)

	dbpool.SetCurrent(New(name, config))
}

func confGet(conf map[string]string, key string) string {
	if v, ok := conf[key]; ok {
		return v
	}
	return ""
}

// New func
// create a database instance
func New(name string, conf map[string]string) *Database {
	if dbpool.GetDatabase(name) != nil {
		panic(name + " this database is already exists!")
	}

	if conf == nil {
		return nil
	}

	driverName := confGet(conf, "driver")
	logger := NewLogger(conf)

	var driver IDriver
	switch driverName {
	case "postgres", "psql":
		driver = &drivers.Postgres{DBName: confGet(conf, "dbname"),
			User:     confGet(conf, "user"),
			Password: confGet(conf, "password"),
			Host:     confGet(conf, "host"),
			HostAddr: confGet(conf, "hostaddr"),
			Port:     confGet(conf, "port"),
		}
		DBTimeFormat = time.RFC3339

	case "mysql":
		driver = &drivers.Mysql{DBName: confGet(conf, "dbname"),
			User:       confGet(conf, "user"),
			Password:   confGet(conf, "password"),
			Protocol:   confGet(conf, "protocol"),
			Parameters: confGet(conf, "parameters"),
			Addr:       confGet(conf, "addr")}

	case "sqlite3":
		driver = &drivers.Sqlite3{File: confGet(conf, "file")}

	default:
		logger.Println("you may need regist a custom driver: db.RegistDriver(Mysql{})")
		driver = &drivers.Common{}

	}

	driver.SetConnectString(confGet(conf, "connect"))
	poolSize, err := strconv.Atoi(confGet(conf, "pool_size"))
	if poolSize == 0 || err != nil {
		poolSize = 3
	}
	maxCount, err := strconv.Atoi(confGet(conf, "max_count"))
	if maxCount == 0 || err != nil {
		maxCount = 10
	}

	d := NewDatabase(name, driver, logger)

	if confGet(conf, "timezone") != "" {
		loc, err := time.LoadLocation(confGet(conf, "timezone"))
		if err != nil {
			panic(err)
		}
		d.Location = loc
	}

	dbpool.SetDatabase(name, d)
	dbpool.SetCurrent(d)
	return d
}

// Query from database on prepare mode
// This function use the current database from database bool
// You can set another database use Use(i) or create an new database use New(name, conf)
func Query(conn *sql.DB, sqlstr []byte, args ...interface{}) (DataSet, error) {
	sqlstring := string(sqlstr)
	Current().Log.Sql(sqlstr, args)

	var rows *sql.Rows
	var err error
	if len(args) == 0 {
		rows, err = conn.Query(sqlstring)
	} else {
		rows, err = conn.Query(sqlstring, args...)
	}

	if err != nil {
		return nil, err
	}
	dataset, err := ScanRows(rows)
	if err != nil {
		return nil, err
	}

	return dataset, nil
	// return Current().Query(NewSqlState([]byte(sqlstr), args))
}

// Query from database on prepare mode
// In prepare mode, the sql command will be cached by database
// This function use the current database from database bool
// You can set another database by Use(i) or New(name, conf) an new database
func QueryPrepare(conn *sql.DB, sqlstr []byte, args ...interface{}) (DataSet, error) {
	sqlstring := string(sqlstr)
	Current().Log.Sql(sqlstr, args)

	var rows *sql.Rows
	var err error

	s, err := conn.Prepare(sqlstring)
	if err != nil {
		return nil, err
	}

	defer s.Close()

	if len(args) == 0 {
		rows, err = s.Query()
	} else {
		rows, err = s.Query(args...)
	}
	if err != nil {
		return nil, err
	}

	dataset, err := ScanRows(rows)
	if err != nil {
		return nil, err
	}

	return dataset, nil
}

func Find(conn *sql.DB, cls IVO, sqlstr []byte, args ...interface{}) (VODataSet, error) {
	dataset, err := Query(conn, sqlstr, args...)
	if err != nil {
		return nil, err
	}

	return NewStructConvert(cls).DataSet2Struct(dataset)
	// return Current().Find(cls, NewSqlState([]byte(sqlstr), args))
}
func FindPrepare(conn *sql.DB, cls IVO, sqlstr []byte, args ...interface{}) (VODataSet, error) {
	dataset, err := QueryPrepare(conn, sqlstr, args...)

	if err != nil {
		return nil, err
	}

	return NewStructConvert(cls).DataSet2Struct(dataset)
}

// Excute sql from a file
// This function run sql under not transaction mode and use the current database from database bool
func ExecFromFile(file string) error {
	var filebytes []byte
	var err error

	if filebytes, err = ioutil.ReadFile(file); err != nil {
		return err
	}
	conn := Current().Connection.Connect()
	b := bytes.Split(filebytes, []byte(";"))
	for _, i := range b {
		if len(bytes.TrimSpace(i)) == 0 {
			continue
		}

		_, err = Exec(conn, i)
		if err != nil {
			return err
		}
	}
	return nil
}

// Excute sql.
// If your has more than on sql command, it will only excute the first.
// This function use the current database from database bool
func Exec(conn *sql.DB, sqlstr []byte, args ...interface{}) (sql.Result, error) {
	sqlstring := string(sqlstr)
	Current().Log.Sql(sqlstr, args)

	if len(args) == 0 {
		return conn.Exec(sqlstring)
	} else {
		return conn.Exec(sqlstring, args...)
	}

}

// Excute sql on prepare mode
// This function use the current database from database bool
func ExecPrepare(conn *sql.DB, sqlstr []byte, args ...interface{}) (sql.Result, error) {
	sqlstring := string(sqlstr)
	Current().Log.Sql(sqlstr, args)

	s, err := conn.Prepare(sqlstring)
	if err != nil {
		return nil, err
	}

	defer s.Close()

	if len(args) == 0 {
		return s.Exec()
	} else {
		return s.Exec(args...)
	}
}

// Get a database instance by name from database pool
func Get(name string) *Database {
	return dbpool.GetDatabase(name)
}

// Return the current database from database pool
func Current() *Database {
	if dbpool.Current() == nil {
		panic("db is not initalized")
	}
	return dbpool.Current()
}

// Set current database by index
func Use(name string) {
	dbpool.Use(name)
}

// Get database count
func DatabaseCount() int {
	return len(dbpool.dblist)
}