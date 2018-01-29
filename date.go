package date

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// A Date specifies the year, month and day.
type Date struct {
	Year  int
	Month time.Month
	Day   int
}

// After reports whether d is after t.
func (d Date) After(t Date) bool {
	if d.Year != t.Year {
		return d.Year > t.Year
	}
	if d.Month != t.Month {
		return d.Month > t.Month
	}
	return d.Day > t.Day
}

// Before reports whether d is before t.
func (d Date) Before(t Date) bool {
	return t.After(d)
}

// Equal reports whether d and t are the same date.
func (d Date) Equal(t Date) bool {
	return !d.After(t) && !d.Before(t)
}

// MarshalText implements the encoding.TextMarshaler interface.
// The output is in "YYYY-MM-DD" format.
func (d Date) MarshalText() ([]byte, error) {
	s := []byte(d.String())
	b := make([]byte, 0, len(s))
	return append(b, s...), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The formats it supports are "2006-01-02", "2006/01/02" and "02 Jan 2006".
func (d *Date) UnmarshalText(b []byte) error {
	var t time.Time
	var err error

	for _, f := range []string{"2006-01-02", "2006/01/02", "02 Jan 2006"} {
		if t, err = time.Parse(f, string(b)); err == nil {
			break
		}
	}
	if err != nil {
		return fmt.Errorf(`Date.UnmarshalText: Unsupported format %s. Only "2006-01-02", "2006/01/02" and "02 Jan 2006" are supported`, b)
	}

	*d = Of(t)
	return nil
}

// Scan implements the sql.Scanner interface.
func (d *Date) Scan(v interface{}) error {
	switch s := v.(type) {
	case time.Time:
		*d = Date{Year: s.Year(), Month: s.Month(), Day: s.Day()}
	case []byte:
		t, err := time.Parse("2006-01-02", string(s))
		if err != nil {
			return err
		}
		*d = Date{Year: t.Year(), Month: t.Month(), Day: t.Day()}
	default:
		return fmt.Errorf("date: Unsupport scanning type %T", v)
	}
	return nil
}

// String returns a string of date in "YYYY-MM-DD" format.
func (d Date) String() string {
	return fmt.Sprintf("%4d-%02d-%02d", d.Year, d.Month, d.Day)
}

// Value implements the driver.Valuer interface.
func (d Date) Value() (driver.Value, error) {
	return d.String(), nil
}

// Now returns the current local date.
func Now() Date {
	n := time.Now()
	return Of(n)
}

// Of returns the Date of t in t's location.
func Of(t time.Time) Date {
	return Date{t.Year(), t.Month(), t.Day()}
}

// Parse parses the d with layout and returns the value of Date.
// The layout follows the format of time.Parse.
func Parse(layout, d string) (Date, error) {
	t, err := time.Parse(layout, d)
	if err != nil {
		return Date{}, err
	}

	return Of(t), nil
}
