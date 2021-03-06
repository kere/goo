package util

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

const (
	// PrintLine for print
	PrintLine = "-----------------------"
	// PrintDotted for print
	PrintDotted = "......................."
	// PrintEquals for print
	PrintEquals = "======================="
)

func PrintObj(obj interface{}) {
	val := reflect.ValueOf(obj)
	if val.Type().Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	n := typ.NumField()
	for i := 0; i < n; i++ {
		ftyp := typ.Field(i)
		name := ftyp.Tag.Get("name")
		if name == "" {
			// name = ftyp.Name
			continue
		}
		hasPer := strings.HasSuffix(name, "%") || strings.HasSuffix(name, "率")

		str := ""
		fval := val.Field(i)
		switch ftyp.Type.Kind() {
		case reflect.Float64:
			deci := 2
			decimal := ftyp.Tag.Get("decimal")
			if decimal != "" {
				v, _ := strconv.ParseInt(decimal, 10, 64)
				if v > 0 {
					deci = int(v)
				}
			}

			v := fval.Float()
			if math.IsInf(v, 0) {
				str = "~"
			} else if math.IsNaN(v) {
				str = "-"
			} else if v == 0 {
				str = "0"
			} else {
				if hasPer {
					str = HumanFloatC(v*100, deci) + "%"
				} else {
					str = HumanFloatC(v, deci)
				}
			}
		case reflect.Int, reflect.Int64:
			if hasPer {
				str = fmt.Sprint(fval.Int()*100) + "%"
			} else {
				str = fmt.Sprint(fval.Int())
			}
		case reflect.String:
			str = fval.String()
		case reflect.Slice:
			v := fval.Interface()
			switch fval.Interface().(type) {
			case []byte:
				str = Bytes2Str(v.([]byte))
			default:
				str = fmt.Sprint(v)
			}

		default:
			str = fmt.Sprintf("%v", fval.Interface())
		}

		end := ""
		endArr := strings.Split(ftyp.Tag.Get("end"), " ")
		for _, e := range endArr {
			switch e {
			case "tab":
				end += "\t"
			case "tab2":
				end += "\t\t"
			case "break":
				end += "\n"
			case "line":
				end += PrintLine
			case "dotted":
				end += PrintDotted
			case "equal":
				end += PrintEquals
			}
		}
		inline := ftyp.Tag.Get("inline")
		switch inline {
		case "tab":
			inline = " \t "
		case "tab2":
			inline = "\t\t"
		case "space":
			inline = " "
		case "space2":
			inline = "  "
		case "space4":
			inline = "    "
		case "comma":
			inline = " , "
		case "colon":
			inline = " : "
		case "minus":
			inline = " - "
		default:
			inline = ""
		}

		suffix := ftyp.Tag.Get("suffix")
		if inline == "" {
			fmt.Print(name, ":", str+suffix, end, "\n")
		} else {
			fmt.Print(name, ":", str+suffix, inline, end)
		}
	}
	fmt.Println()
}
