package date

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var (

	// databases
	db map[string]*sql.DB

	// test cases
	d1 = Date{2001, time.March, 5}
	d2 = Date{2009, time.November, 15}
	f1 = []string{
		"2001-03-05",
		"2001/03/05",
		"Mar 05 2001",
	}
	f2 = []string{
		"2009-11-15",
		"2009/11/15",
		"15 Nov 2009",
	}
)

// setup databases
func init() {
	cfg := map[string]string{
		"mysql":    "travis@tcp(localhost)/date_test",
		"postgres": "dbname=date_test user=postgres sslmode=disable",
		"sqlite3":  ":memory:",
	}

	var err error
	for k, v := range cfg {
		db[k], err = sql.Open(k, v)
		if err != nil {
			fmt.Print(err)
			delete(db, k)
		}
	}
}

func TestAfter(t *testing.T) {
	if d1.After(d2) {
		t.Error("Date.After failed: d1 should not after d2")
	}
	if d1.After(d1) {
		t.Error("Date.After failed: d1 should not after d1")
	}
}

func TestBefore(t *testing.T) {
	if d2.Before(d1) {
		t.Error("Date.Before failed: d2 should not before d1")
	}
	if d2.After(d2) {
		t.Error("Date.Before failed: d2 should not after d2")
	}
}

func TestEqual(t *testing.T) {
	if d1.Equal(d2) {
		t.Error("Date.Equal failed: d1 != d2")
	}
	if !d2.Equal(d2) {
		t.Error("Date.Equal failed: d2 == d2")
	}
}

func TestMarshalJSON(t *testing.T) {
	b, err := json.Marshal(d1)
	if err != nil {
		t.Errorf("Date.MarshalJSON exits with error: %v", err)
	}

	w := `"` + f1[0] + `"`
	if string(b) != w {
		t.Errorf("Date.MarshalJSON failed. want: %s, got: %s", w, string(b))
	}
}

func TestUnmarshalJSON(t *testing.T) {
	var d Date
	for _, b := range f2 {
		err := json.Unmarshal([]byte(`"`+b+`"`), &d)
		if err != nil {
			t.Errorf("Date.UnmarshalJSON exits with error: %v", err)
		}
		if !d.Equal(d2) {
			t.Errorf("Date.UnmarshalJSON failed. want: %v, got: %v", d2, d)
		}
	}
}

func TestMarshalText(t *testing.T) {
	b, err := d1.MarshalText()
	if err != nil {
		t.Errorf("Date.MarshalText exits with error: %v", err)
	}

	w := f1[0]
	if string(b) != w {
		t.Errorf("Date.MarshalText failed. want: %s, got: %s", w, string(b))
	}
}

func TestUnmarshalText(t *testing.T) {
	var d Date
	for _, b := range f2 {
		err := d.UnmarshalText([]byte(b))
		if err != nil {
			t.Errorf("Date.UnmarshalText exits with error: %v", err)
		}
		if !d.Equal(d2) {
			t.Errorf("Date.UnmarshalText failed. want: %v, got: %v", d2, d)
		}
	}
}

func TestScan(t *testing.T) {
	q := map[string]string{
		"mysql":    "SELECT CURDATE();",
		"postgres": "SELECT now()::date;",
		"sqlite3":  "SELECT date('now');",
	}

	d := Date{}
	n := time.Now()

	for k := range db {
		err := db[k].QueryRow(q[k]).Scan(&d)
		if err != nil {
			t.Errorf("Scanning from %s exits with error: %v", k, err)
		}

		if d.Year != n.Year() || d.Month != n.Month() || d.Day != n.Day() {
			t.Errorf("Date.Scan failed: want: %s, got: %v", n.Format("2006-01-02"), d)
		}
	}
}

func TestString(t *testing.T) {
	w := f2[0]
	if d2.String() != w {
		t.Errorf("Date.String failed. want: %s, got: %v", w, d2)
	}
}

func TestValue(t *testing.T) {
	q := map[string]string{
		"mysql":    "SELECT DATE(?);",
		"postgres": "SELECT $1::date;",
		"sqlite3":  "SELECT date(?);",
	}

	d := Date{}

	for k := range db {
		err := db[k].QueryRow(q[k], d1).Scan(&d)
		if err != nil {
			t.Errorf("Scanning from %s exits with error: %v", k, err)
		}

		if !d.Equal(d1) {
			t.Errorf("Date.Value failed: want: %v, got: %v", d1, d)
		}
	}
}

func TestNow(t *testing.T) {
	d := Now()
	n := time.Now()
	if d.Year != n.Year() || d.Month != n.Month() || d.Day != n.Day() {
		t.Errorf("Now failed. Want: %s, got: %v", n.Format("2006/01/02"), d)
	}
}

func TestParse(t *testing.T) {
	layouts := []string{
		"2006-01-02",
		"2006/01/02",
		"Jan 2 2006",
	}
	for i, l := range layouts {
		d, err := Parse(l, f1[i])
		if err != nil {
			t.Errorf("Parse exit with error: %v", err)
		}
		if d.Year != d1.Year ||
			d.Month != d1.Month ||
			d.Day != d1.Day {
			t.Errorf("Parse failed: parse %s returns %v", f1[i], d)
		}
	}

	// Invalid date should return error
	_, err := Parse(layouts[0], "2017/13/01")
	if err == nil {
		t.Error("Parse failed: Invalid input doesn't return error")
	}
}
