package helper

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

// jadi status HTTP + pesan yang bisa dibaca user.
func MapPgError(err error) (int, string, bool) {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return 0, "", false
	}
	switch pgErr.Code {
	case "23505": // unique_violation
		switch pgErr.ConstraintName {
		case "personel_nrp_nip_key", "uq_personel_nrp_aktif":
			return http.StatusConflict, "NRP/NIP sudah terdaftar", true
		case "users_username_key":
			return http.StatusConflict, "username sudah dipakai", true
		case "api_clients_client_id_key":
			return http.StatusConflict, "client_id sudah dipakai aplikasi lain", true
		case "uq_jabatan_aktif_per_personel":
			return http.StatusConflict, "personel sudah memiliki jabatan definitif aktif (tmt_selesai kosong)", true
		}
		return http.StatusConflict, "data sudah terdaftar (duplikat): " + pgErr.ConstraintName, true
	case "23503": // foreign_key_violation
		return http.StatusUnprocessableEntity, "referensi tidak valid (pangkat/satker tidak ada atau tidak sesuai jenis personel)", true
	case "23514": // check_violation
		return http.StatusUnprocessableEntity, "data melanggar ketentuan sistem: " + pgErr.ConstraintName, true
	case "23P01": // exclusion_violation
		return http.StatusConflict, "periode jabatan tumpang tindih dengan jabatan lain", true
	case "P0001": // RAISE EXCEPTION dari trigger
		return http.StatusUnprocessableEntity, pgErr.Message, true
	}
	return 0, "", false
}
