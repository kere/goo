package util

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func Test_Func(t *testing.T) {
	num2 := 1.712774821
	val := Round(num2, 2)
	if val != 1.71 {
		t.Fatal("round failed", val)
	}
	val = Round(num2, 3)
	if val != 1.713 {
		t.Fatal("round failed", val)
	}
	val = Round(num2, 0)
	if val != 2.0 {
		t.Fatal("round failed", val)
	}
	val = Round(num2, 4)
	if val != 1.7128 {
		t.Fatal("round failed", val)
	}
	num2 = -1.712774821
	val = Round(num2, 2)
	if val != -1.71 {
		t.Fatal("round failed", val)
	}
	val = Round(num2, 3)
	if val != -1.713 {
		t.Fatal("round failed", val)
	}
	val = Round(num2, 0)
	if val != -2.0 {
		t.Fatal("round failed", val)
	}
	val = Round(num2, 4)
	if val != -1.7128 {
		t.Fatal("round failed", val)
	}

	// str := HumanFloatC(3.1415926)
	str := HumanInt64C(1234567890)
	if str != "12,3456,7890" {
		t.Fatal(str)
	}
	str = HumanFloatC(12345678.12345)
	if str != "1234,5678.12345" {
		t.Fatal(str)
	}
}

func TestSort(t *testing.T) {
	arr := []int64{10, 322, 3, 43, 65, 30, 230, 44, 56, 76, 20, 430, 659}
	arrSort := Int64sOrder(arr)
	arrSort.Sort()

	index := IndexOfInt64s(43, arrSort)
	if index != 4 {
		t.Fatal("index:", index)
	}
	index = IndexOfInt64s(45, arrSort)
	if index != -1 {
		t.Fatal("index:", index)
	}

	arr = []int64{3, 10, 20, 30, 43, 44, 45, 56, 65, 76, 230, 322, 430, 659}
	arrSort = Int64sOrder(arr)
	index = SearchInt64s(50, arrSort)

	if index != 7 {
		fmt.Println(arr)
		t.Fatal(index)
	}

	index = SearchInt64s(57, arrSort)
	if index != 8 {
		fmt.Println(arr)
		t.Fatal(index)
	}

	arr = []int64{3, 10, 20, 30, 43, 44, 45, 49, 56, 65, 76, 230, 322, 430, 659}
	arrSort = Int64sOrder(arr)
	index = SearchInt64s(50, arrSort)
	if index != 8 {
		fmt.Println(arr)
		t.Fatal(index)
	}
	index = SearchInt64s(48, arrSort)
	if index != 7 {
		fmt.Println(arr)
		t.Fatal(index)
	}

	arr = []int64{10, 322, 3, 3, 43, 65, 30, 230, 30, 44, 56, 20, 76, 20, 430, 659}
	tmp := Int64sUnique(arr, false)
	if tmp[1] != 10 || tmp[2] != 20 || tmp[3] != 30 {
		fmt.Println(tmp)
		t.Fatal("Int64sUniqueP")
	}
}

func TestSync(t *testing.T) {
	cpt := NewComputation(20)
	arr := make([]int, 100)

	counter := 0
	cpt.RunA(100, func(i int) (interface{}, error) {
		arr[i] = i + 1
		return i, nil
	}, func(i int, dat interface{}) {
		counter++
	})

	if counter != 100 {
		fmt.Println(counter, "===========")
		fmt.Println(arr, "===========")
		t.Fatal(counter)
	}

	s := []byte("110010")
	v := MaskBytes2Int(s)
	if v != 50 {
		t.Fatal(v)
	}
	if IsMaskTrueAt(v, 0) {
		t.Fatal(string(s))
	}
	if !IsMaskTrueAt(v, 1) {
		t.Fatal(string(s))
	}

	v = SetIntMask(v, 0, true)
	str := strconv.FormatInt(int64(v), 2)
	if str != "110011" {
		t.Fatal(str)
	}
	v = SetIntMask(v, 2, true)
	str = strconv.FormatInt(int64(v), 2)
	if str != "110111" {
		t.Fatal(str)
	}

}

func TestSync2(t *testing.T) {
	count := 1000000
	row := []int{0}
	now := time.Now()
	for i := 0; i < count; i++ {
		row[0] += i
	}
	fmt.Println("for:", time.Now().Sub(now).String())
	fmt.Println(row[0])

	row[0] = 0
	now = time.Now()
	cpt := NewComputation(1000)
	lock := sync.Mutex{}
	cpt.Run(count, func(i int) {
		lock.Lock()
		row[0] += i
		lock.Unlock()
	})
	fmt.Println("pool:", time.Now().Sub(now).String())
	fmt.Println(row[0])

	row[0] = 0
	now = time.Now()
	wg := sync.WaitGroup{}
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(i int) {
			lock.Lock()
			row[0] += i
			lock.Unlock()
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println("WaitGroup:", time.Now().Sub(now).String())
	fmt.Println(row[0])
}

func TestPool(t *testing.T) {
	r := make([]float64, 5)
	for i := 0; i < 5; i++ {
		r[i] = float64(i + 1)
	}

	PutFloats(r)
	r = GetFloats(10)
	if len(r) != 10 {
		t.Fatal(r)
	}

	PutFloats(r)
	r = GetFloats(10)
	if len(r) != 10 {
		t.Fatal(r)
	}

	PutFloats(r)
	r = GetFloats(5)
	if len(r) != 5 {
		t.Fatal(r)
	}
}

func TestCamelCase(t *testing.T) {
	s := "created_at_some"
	str := CamelCase(s)
	if str != "CreatedAtSome" {
		t.Fatal(str)
	}
}
func TestNum(t *testing.T) {
	v := StrNumType("0")
	if v != 'i' {
		t.Fatal(rune(v))
	}
	v = StrNumType("+0")
	if v != 'i' {
		t.Fatal(rune(v))
	}
	v = StrNumType("-0")
	if v != 'i' {
		t.Fatal(rune(v))
	}
	v = StrNumType("+1230")
	if v != 'i' {
		t.Fatal(rune(v))
	}
	v = StrNumType("e")
	if v != 's' {
		t.Fatal(rune(v))
	}
	v = StrNumType(".xyce")
	if v != 's' {
		t.Fatal(rune(v))
	}
	v = StrNumType("+xyce")
	if v != 's' {
		t.Fatal(rune(v))
	}
	v = StrNumType("+e")
	if v != 's' {
		t.Fatal(rune(v))
	}
	v = StrNumType("-e")
	if v != 's' {
		t.Fatal(rune(v))
	}
	v = StrNumType("e+")
	if v != 's' {
		t.Fatal(rune(v))
	}
	v = StrNumType("e-0")
	if v != 's' {
		t.Fatal(rune(v))
	}
	v = StrNumType("3e+")
	if v != 's' {
		t.Fatal(rune(v))
	}
	v = StrNumType("3.14e+")
	if v != 's' {
		t.Fatal(rune(v))
	}
	v = StrNumType(".314.")
	if v != 's' {
		t.Fatal(rune(v))
	}
	v = StrNumType(".314e+")
	if v != 's' {
		t.Fatal(rune(v))
	}
	v = StrNumType("3.14e+03")
	if v != 'f' {
		t.Fatal(rune(v))
	}
	v = StrNumType("+3.14e+03")
	if v != 'f' {
		t.Fatal(rune(v))
	}
	v = StrNumType("3.14e-03")
	if v != 'f' {
		t.Fatal(rune(v))
	}
	v = StrNumType("-3.14e-03")
	if v != 'f' {
		t.Fatal(rune(v))
	}
	v = StrNumType("314e+03")
	if v != 'f' {
		t.Fatal(rune(v))
	}
	v = StrNumType("+314e+03")
	if v != 'f' {
		t.Fatal(rune(v))
	}
	v = StrNumType("314e-03")
	if v != 'f' {
		t.Fatal(rune(v))
	}
	v = StrNumType(".314")
	if v != 'f' {
		t.Fatal(rune(v))
	}
	v = StrNumType("314.")
	if v != 'f' {
		t.Fatal(rune(v))
	}
	v = StrNumType("+.314")
	if v != 'f' {
		t.Fatal(rune(v))
	}
	v = StrNumType(".314e-03")
	if v != 'f' {
		t.Fatal(rune(v))
	}
	v = StrNumType("+.314e-03")
	if v != 'f' {
		t.Fatal(rune(v))
	}
}

func TestSplit(t *testing.T) {
	src := "a,b,c"
	arr := SplitStrNotSafe(src, SComma)
	if len(arr) != 3 || arr[0] != "a" || arr[2] != "c" {
		t.Fatal(strings.Join(arr, "-"))
	}
	src = "a,b,c,"
	arr = SplitStrNotSafe(src, SComma)
	if len(arr) != 4 || arr[0] != "a" || arr[2] != "c" {
		t.Fatal(strings.Join(arr, "-"))
	}
	src = ",a,b,c"
	arr = SplitStrNotSafe(src, SComma)
	if len(arr) != 4 || arr[1] != "a" || arr[3] != "c" {
		t.Fatal(strings.Join(arr, "-"))
	}
	src = ",a,b,c,"
	arr = SplitStrNotSafe(src, SComma)
	if len(arr) != 5 || arr[1] != "a" || arr[3] != "c" {
		t.Fatal(strings.Join(arr, "-"))
	}
	src = "a,b,c,,"
	arr = SplitStrNotSafe(src, SComma)
	if len(arr) != 5 || arr[0] != "a" || arr[2] != "c" {
		t.Fatal(strings.Join(arr, "-"))
	}
	src = ",,a,b,c"
	arr = SplitStrNotSafe(src, SComma)
	if len(arr) != 5 || arr[2] != "a" || arr[4] != "c" {
		t.Fatal(strings.Join(arr, "-"))
	}
	src = ",,a,b,c,,,"
	arr = SplitStrNotSafe(src, SComma)
	if len(arr) != 8 || arr[2] != "a" || arr[4] != "c" {
		t.Fatal(strings.Join(arr, "-"))
	}
	fmt.Println(strings.Join(arr, "-"))
}

// func TestEach(t *testing.T) {
// 	arr := []int{0}
// 	EachPartByN(len(arr), 3, func(a int, b int) bool {
// 		fmt.Println(arr[a:b])
// 		return true
// 	})
// }
