package model

const (
	RoleAdminSSDM = "ADMIN_SSDM"
	RoleOperator  = "OPERATOR"
)

const (
	PermPersonelRead   = "personel:read"
	PermPersonelWrite  = "personel:write"
	PermPersonelDelete = "personel:delete"
	PermJabatanRead    = "jabatan:read"
	PermJabatanWrite   = "jabatan:write"
	PermJabatanDelete  = "jabatan:delete"
	PermUserManage     = "user:manage" // khusus Admin SSDM
)

// RolePermissions: APA yang boleh dilakukan tiap role.
// Batas DATA MANA yang boleh disentuh diatur terpisah lewat Principal.SatkerID.
var RolePermissions = map[string][]string{
	RoleAdminSSDM: {
		PermPersonelRead, PermPersonelWrite, PermPersonelDelete,
		PermJabatanRead, PermJabatanWrite, PermJabatanDelete,
		PermUserManage,
	},
	RoleOperator: {
		PermPersonelRead, PermPersonelWrite, PermPersonelDelete,
		PermJabatanRead, PermJabatanWrite, PermJabatanDelete,
	},
}

// ValidScopes: scope yang boleh diberikan ke API client.
// Harus sama dengan CHECK constraint chk_scopes_valid di tabel api_clients.
var ValidScopes = []string{
	"personel:read", "personel:write", "personel:delete",
	"jabatan:read", "jabatan:write", "jabatan:delete",
	"pangkat:read", "pangkat:write",
	"kualifikasi:read", "kualifikasi:write", "kualifikasi:delete",
}
