package model

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

const DateLayout = "2006-01-02"

type Date struct {
	time.Time
}

// JSON → Date (input dari client)
func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		return fmt.Errorf("tanggal wajib diisi dengan format YYYY-MM-DD")
	}
	t, err := time.Parse(DateLayout, s)
	if err != nil {
		return fmt.Errorf("format tanggal '%s' tidak valid, gunakan YYYY-MM-DD", s)
	}
	d.Time = t
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.Format(DateLayout) + `"`), nil
}

// DB → Date (dipakai GORM saat SELECT)
func (d *Date) Scan(value any) error {
	switch v := value.(type) {
	case time.Time:
		d.Time = v
		return nil
	case string:
		t, err := time.Parse(DateLayout, v)
		d.Time = t
		return err
	case []byte:
		t, err := time.Parse(DateLayout, string(v))
		d.Time = t
		return err
	case nil:
		d.Time = time.Time{}
		return nil
	}
	return fmt.Errorf("tidak bisa scan %T ke Date", value)
}

// Date → DB (dipakai GORM saat INSERT/UPDATE)
func (d Date) Value() (driver.Value, error) {
	return d.Format(DateLayout), nil
}
