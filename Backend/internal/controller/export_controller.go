package controller

import (
	"Backend/internal/helper"
	"Backend/internal/middleware"
	"Backend/internal/repository"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type ExportController struct {
	dashboard repository.DashboardRepository
	export    repository.ExportRepository
}

func NewExportController(d repository.DashboardRepository, e repository.ExportRepository) *ExportController {
	return &ExportController{dashboard: d, export: e}
}

// GET /api/export/personel → file .xlsx. Menerima filter yang sama dengan GET /personel.
func (ctl *ExportController) Personel(c *gin.Context) {
	p := middleware.CurrentPrincipal(c)
	if p == nil {
		helper.Fail(c, http.StatusUnauthorized, "unauthenticated", nil)
		return
	}
	var q listPersonelQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		helper.Fail(c, http.StatusBadRequest, "parameter tidak valid", err.Error())
		return
	}
	filter := repository.PersonelFilter{
		Q: q.Q, SatkerID: q.SatkerID, Status: q.Status, Jenis: q.Jenis,
		Kelompok: q.Kelompok, TanpaJabatan: q.TanpaJabatan, AkanPensiun: q.AkanPensiun,
		Urut: q.Sort,
	}
	// Tanpa filter → sertakan sheet rekap per satker.
	adaFilter := q.Q != "" || q.SatkerID != nil || q.Status != "" || q.Jenis != "" || q.Kelompok != "" || q.TanpaJabatan || q.AkanPensiun

	daftar, err := ctl.export.DaftarPersonel(filter, p.SatkerID)
	if err != nil {
		respondError(c, err)
		return
	}

	f := excelize.NewFile()
	defer f.Close()
	judul, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"13294B"}},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})

	tulis := func(sheet string, header []string, lebar []float64, rows [][]any) {
		for i, h := range header {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheet, cell, h)
			col, _ := excelize.ColumnNumberToName(i + 1)
			f.SetColWidth(sheet, col, col, lebar[i])
		}
		last, _ := excelize.CoordinatesToCellName(len(header), 1)
		f.SetCellStyle(sheet, "A1", last, judul)
		for r, row := range rows {
			for i, v := range row {
				cell, _ := excelize.CoordinatesToCellName(i+1, r+2)
				f.SetCellValue(sheet, cell, v)
			}
		}
		f.SetPanes(sheet, &excelize.Panes{Freeze: true, YSplit: 1, TopLeftCell: "A2", ActivePane: "bottomLeft"})
		end, _ := excelize.CoordinatesToCellName(len(header), len(rows)+1)
		f.AutoFilter(sheet, "A1:"+end, nil)
	}

	// Sheet daftar personel
	f.SetSheetName("Sheet1", "Daftar Personel")
	tgl := func(t *time.Time) any {
		if t == nil {
			return ""
		}
		return t.Format("2006-01-02")
	}
	str := func(s *string) string {
		if s == nil {
			return ""
		}
		return *s
	}
	var baris [][]any
	for i, d := range daftar {
		baris = append(baris, []any{
			i + 1, d.NrpNip, d.Nama, d.Jenis, d.Pangkat, helper.Rapikan(d.Satker), helper.Rapikan(str(d.Jabatan)),
			tgl(d.TmtJabatan), str(d.TempatLahir), d.TanggalLahir.Format("2006-01-02"),
			pensiun(d.TanggalLahir).Format("2006-01-02"), d.Status,
		})
	}
	tulis("Daftar Personel",
		[]string{"No", "NRP/NIP", "Nama", "Jenis", "Pangkat", "Satker", "Jabatan saat ini", "TMT jabatan", "Tempat lahir", "Tanggal lahir", "Pensiun (58 th)", "Status"},
		[]float64{6, 20, 30, 8, 12, 35, 50, 13, 18, 13, 15, 14}, baris)

	// Sheet rekap per satker (hanya untuk export tanpa filter)
	if !adaFilter {
		ringkasan, err := ctl.dashboard.Ringkasan(p.SatkerID)
		if err != nil {
			respondError(c, err)
			return
		}
		var rekap [][]any
		for _, s := range ringkasan.PerSatker {
			rekap = append(rekap, []any{helper.Rapikan(s.Nama), s.Kode, s.Total, s.Polri, s.Pns, s.AkanPensiun})
		}
		f.NewSheet("Rekap Satker")
		tulis("Rekap Satker",
			[]string{"Satker", "Kode", "Total", "POLRI", "PNS", "Pensiun ≤ 12 bulan"},
			[]float64{40, 30, 10, 10, 10, 18}, rekap)
		idx, _ := f.GetSheetIndex("Rekap Satker")
		f.SetActiveSheet(idx)
		f.MoveSheet("Rekap Satker", "Daftar Personel")
	}

	nama := fmt.Sprintf("personel_%s.xlsx", time.Now().Format("2006-01-02"))
	c.Header("Content-Disposition", `attachment; filename="`+nama+`"`)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	if err := f.Write(c.Writer); err != nil {
		helper.Fail(c, http.StatusInternalServerError, "gagal membuat file", nil)
	}
}

// pensiun: tanggal mencapai usia 58 tahun (29 Feb → 28 Feb, sama seperti Postgres).
func pensiun(lahir time.Time) time.Time {
	t := lahir.AddDate(58, 0, 0)
	if t.Day() != lahir.Day() {
		t = t.AddDate(0, 0, -t.Day())
	}
	return t
}
