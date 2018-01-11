package drivers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

var (
	b_r_HSTORE = []byte("\"=>\"")
	b_r_JSON   = []byte("\":\"")
)

type Postgres struct {
	Common
	DBName   string
	User     string
	Password string
	Host     string
	HostAddr string
	Port     string
}

func (this *Postgres) DriverName() string {
	return "postgres"
}

func (this *Postgres) AdaptSql(bSQL []byte) []byte {
	arr := bytes.Split(bSQL, B_QuestionMark)
	l := len(arr)
	s := bytes.Buffer{}
	s.Write(arr[0])

	for i := 1; i < l; i++ {
		if bytes.HasPrefix(arr[i], b_Dollar) {
			s.Write(B_QuestionMark)
			s.Write(arr[i][1:])
		} else {
			s.Write(b_Dollar)
			s.WriteString(fmt.Sprint(i))
			s.Write(arr[i])
		}
	}
	return s.Bytes()
}

func (this *Postgres) ConnectString() string {
	if this.Host == "" {
		this.Host = "127.0.0.1"
	}
	if this.Port == "" {
		this.Port = "5432"
	}

	if this.HostAddr != "" {
		return fmt.Sprintf("dbname=%s user=%s password=%s hostaddr=%s sslmode=disable",
			this.DBName,
			this.User,
			this.Password,
			this.HostAddr)

	} else {
		return fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s sslmode=disable",
			this.DBName,
			this.User,
			this.Password,
			this.Host,
			this.Port)
	}
}

func (this *Postgres) QuoteField(s string) string {
	return fmt.Sprint("\"", s, "\"")
}

func (this *Postgres) LastInsertId(table, pkey string) string {
	// return "select currval(pg_get_serial_sequence('" + table + "','" + pkey + "'))"
	return fmt.Sprint("select currval(pg_get_serial_sequence('", table, "','", pkey, "')) as count")
}

func (this *Postgres) sliceToStore(typ reflect.Type, v interface{}) string {
	switch typ.Kind() {
	case reflect.Slice, reflect.Array:
		value := reflect.ValueOf(v)
		arr := make([]string, value.Len())
		l := value.Len()
		if l == 0 {
			return "{}"
		}

		var tmpV reflect.Value
		for i := 0; i < l; i++ {
			tmpV = value.Index(i)
			arr[i] = this.sliceToStore(tmpV.Type(), tmpV.Interface())
		}
		return fmt.Sprint("{", strings.Join(arr, ","), "}")

	case reflect.String:
		return fmt.Sprint("'", v, "'")

	default:
		return fmt.Sprint(v)

	}

}

func (this *Postgres) FlatData(v interface{}) interface{} {
	if v == nil {
		return nil
	}

	typ := reflect.TypeOf(v)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	switch typ.Kind() {
	// case reflect.String, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
	default:
		return v

	case reflect.Bool:
		if v.(bool) {
			return "t"
		} else {
			return "f"
		}

	case reflect.Array:
		return this.sliceToStore(typ, v)

	case reflect.Slice:
		switch v.(type) {
		case []byte:
			return v
		default:
			return this.sliceToStore(typ, v)
		}

	case reflect.Map:
		switch v.(type) {
		case map[string]string:
			var valStr string
			hdata := v.(map[string]string)
			arr := make([]string, len(hdata))
			i := 0
			for kk, vv := range hdata {
				valStr = strings.Replace(fmt.Sprint(vv), "\"", "\\\"", -1)
				arr[i] = fmt.Sprint("\"", kk, "\"", "=>", "\"", valStr, "\"")
				i++
			}
			return fmt.Sprint(strings.Join(arr, ","))

		default:
			b, err := json.Marshal(v)
			if err != nil {
				return []byte("")
			}
			return b

		}

	case reflect.Struct:
		switch v.(type) {
		case time.Time:
			return v

		default:
			b, err := json.Marshal(v)
			if err != nil {
				return []byte("")
			}
			return b
		}
	}

}

func (this *Postgres) StringSlice(src []byte) ([]string, error) {
	if len(src) == 0 {
		return []string{}, nil
	}

	src = bytes.TrimPrefix(src, b_BRACE_LEFT)
	src = bytes.TrimSuffix(src, b_BRACE_RIGHT)
	if len(src) == 0 {
		return []string{}, nil
	}

	l := bytes.Split(src, b_COMMA)
	v := make([]string, len(l))
	for i, _ := range l {
		v[i] = string(bytes.Trim(l[i], "'"))
	}

	return v, nil
}

func (this *Postgres) Int64Slice(src []byte) ([]int64, error) {
	if len(src) == 0 {
		return []int64{}, nil
	}
	var arr = make([]int64, 0)
	if err := this.ParseNumberSlice(src, &arr); err != nil {
		return nil, err
	}

	return arr, nil
}

func (this *Postgres) ParseStringSlice(src []byte, ptr interface{}) error {
	src = bytes.Replace(src, b_BRACE_LEFT, b_BRACKET_LEFT, -1)
	src = bytes.Replace(src, b_BRACE_RIGHT, b_BRACKET_RIGHT, -1)
	src = bytes.Replace(src, b_Quote, b_DoubleQuote, -1)

	if err := json.Unmarshal(src, ptr); err != nil {
		return fmt.Errorf("json parse error: %s \nsrc=%s", err.Error(), src)
	}

	return nil
}

func (this *Postgres) HStore(src []byte) (map[string]string, error) {
	src = bytes.Replace(src, b_r_HSTORE, b_r_JSON, -1)
	src = append(b_BRACE_LEFT, src...)
	v := make(map[string]string)

	if err := json.Unmarshal(append(src, b_BRACE_RIGHT...), &v); err != nil {
		return nil, fmt.Errorf("json parse error: %s \nsrc=%s", err.Error(), src)
	}
	return v, nil
}

func (this *Postgres) ParseNumberSlice(src []byte, ptr interface{}) error {
	if len(src) == 0 {
		return nil
	}

	fmt.Println("==a==")

	src = bytes.Replace(src, b_BRACE_LEFT, b_BRACKET_LEFT, -1)
	src = bytes.Replace(src, b_BRACE_RIGHT, b_BRACKET_RIGHT, -1)
	src = bytes.Replace(src, bNaN, bZero, -1)

	if err := json.Unmarshal(src, ptr); err != nil {
		return err
	}

	return nil
}