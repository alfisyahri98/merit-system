package model

import "slices"

const (
	KindUser   = "user"   // manusia: Admin SSDM / Operator (tabel users)
	KindClient = "client" // aplikasi lain (tabel api_clients)
)

// Principal: identitas pemanggil yang sudah terverifikasi.
// User dan aplikasi lain sama-sama jadi Principal, sehingga
// RequirePermission & scope satker berlaku sama untuk keduanya.
type Principal struct {
	Kind        string
	ID          int
	Name        string
	Role        string
	SatkerID    *int // nil = semua satker
	Permissions []string
}

func (p *Principal) Can(perm string) bool {
	return slices.Contains(p.Permissions, perm)
}

func (p *Principal) AllSatker() bool {
	return p.SatkerID == nil
}
