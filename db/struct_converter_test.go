package db

import (
	"testing"
	"time"
)

// VO value object
type convVO struct {
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	FinishedAt time.Time `json:"finished_at"`
}

func TestConvVO(t *testing.T) {
	now := time.Now()
	row := MapRow{"code": "code1", "name": "tom01", "finished_at": now}
	vo := convVO{}
	row.CopyToWithJSON(&vo)
	if vo.Code != row.String("code") && vo.FinishedAt.String() != now.String() {
		t.Fatal()
	}
}

func TestIsEmpty(t *testing.T) {
	vo := convVO{}
	cv := NewStructConvert(vo)
	row := cv.Struct2DataRow(ActionInsert)

	if row.IsEmpty() {
		t.Fatal("is empty failed", row)
	}

	vo.Name = "tom"
	vo.FinishedAt = time.Now()
	cv = NewStructConvert(vo)
	row = cv.Struct2DataRow(ActionUpdate)

	if len(row) != 3 {
		t.Fatal("is empty failed", row)
	}

}
