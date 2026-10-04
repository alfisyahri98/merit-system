package model

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

// ApiClient: aplikasi lain yang boleh mengakses REST API (tabel api_clients).
type ApiClient struct {
	ID            int         `json:"id"`
	NamaAplikasi  string      `json:"nama_aplikasi"`
	ClientID      string      `json:"client_id"`
	SecretHash    string      `json:"-"`
	Scopes        StringArray `json:"scopes"`
	AksesSatkerID *int        `json:"akses_satker_id"` // nil = semua satker
	IsActive      bool        `json:"is_active"`
	CreatedAt     time.Time   `json:"created_at"`
}

func (ApiClient) TableName() string { return "api_clients" }

// StringArray: mapping kolom PostgreSQL TEXT[] ↔ []string.
type StringArray []string

func (a *StringArray) Scan(value any) error {
	var s string
	switch v := value.(type) {
	case nil:
		*a = nil
		return nil
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fmt.Errorf("tidak bisa scan %T ke StringArray", value)
	}

	s = strings.TrimSpace(s)
	if len(s) < 2 || s[0] != '{' || s[len(s)-1] != '}' {
		return fmt.Errorf("format array postgres tidak valid: %q", s)
	}
	inner := s[1 : len(s)-1]
	if inner == "" {
		*a = StringArray{}
		return nil
	}
	parts := strings.Split(inner, ",")
	out := make(StringArray, 0, len(parts))
	for _, p := range parts {
		out = append(out, strings.Trim(p, `"`))
	}
	*a = out
	return nil
}

func (a StringArray) Value() (driver.Value, error) {
	quoted := make([]string, len(a))
	for i, s := range a {
		quoted[i] = `"` + strings.ReplaceAll(s, `"`, `\"`) + `"`
	}
	return "{" + strings.Join(quoted, ",") + "}", nil
}
