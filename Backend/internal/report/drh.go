// Package report membuat dokumen cetak (PDF) dari data personel.
package report

import (
	"Backend/internal/helper"
	"Backend/internal/model"
	"bytes"
	_ "embed"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
)

//go:embed logo.png
var logoSSDM []byte

const (
	kiri     = 15.0
	lebarIsi = 180.0 // A4 210 mm - margin kiri/kanan 15 mm
	usiaBUP  = 58    // batas usia pensiun

	batasBawah = 18.0               // margin bawah (footer di dalamnya)
	akhirIsi   = 297.0 - batasBawah // posisi Y terakhir untuk isi
)

var bulan = []string{"Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}

// drh menyimpan state penulisan supaya fungsi-fungsi kecil di bawah tetap ringkas.
type drh struct {
	pdf *fpdf.Fpdf
	tr  func(string) string // UTF-8 → cp1252 untuk font bawaan
}

// DaftarRiwayatHidup menghasilkan PDF DRH satu personel.
// p harus sudah memuat Pangkat, Satker, dan RiwayatJabatan (+Satker, Fungsi).
func DaftarRiwayatHidup(p *model.Personel, dicetakOleh string, sekarang time.Time) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(kiri, 15, kiri)
	pdf.SetAutoPageBreak(true, batasBawah)
	pdf.SetTitle("Daftar Riwayat Hidup - "+p.Nama, true)
	pdf.SetAuthor("Merit System Personel Polri", true)
	pdf.AliasNbPages("{nb}")
	d := &drh{pdf: pdf, tr: pdf.UnicodeTranslatorFromDescriptor("")}

	pdf.RegisterImageOptionsReader("logo", fpdf.ImageOptions{ImageType: "PNG"}, bytes.NewReader(logoSSDM))
	pdf.SetFooterFunc(func() {
		pdf.SetY(-14)
		pdf.SetDrawColor(203, 213, 225)
		pdf.Line(kiri, pdf.GetY(), kiri+lebarIsi, pdf.GetY())
		pdf.Ln(1.5)
		pdf.SetFont("Helvetica", "", 7)
		pdf.SetTextColor(100, 116, 139)
		pdf.CellFormat(lebarIsi/2, 4, d.tr(fmt.Sprintf("Dicetak oleh %s, %s", dicetakOleh, sekarang.Format("02-01-2006 15:04"))), "", 0, "L", false, 0, "")
		pdf.CellFormat(lebarIsi/2, 4, d.tr(fmt.Sprintf("Halaman %d dari {nb}", pdf.PageNo())), "", 0, "R", false, 0, "")
	})
	pdf.AddPage()

	d.kop()
	d.judul()
	d.identitas(p, sekarang)

	// Riwayat jabatan, terbaru di atas
	riwayat := append([]model.RiwayatJabatan(nil), p.RiwayatJabatan...)
	sort.SliceStable(riwayat, func(i, j int) bool {
		ai, aj := riwayat[i].TmtSelesai == nil, riwayat[j].TmtSelesai == nil
		if ai != aj {
			return ai // jabatan saat ini paling atas
		}
		return riwayat[i].TmtMulai.After(riwayat[j].TmtMulai.Time)
	})
	d.bagian("I. Ringkasan Karier")
	d.ringkasan(riwayat, len(p.Sertifikasi), sekarang)
	pdf.Ln(3)

	d.bagian("II. Riwayat Jabatan")
	d.daftarJabatan(riwayat)
	pdf.Ln(3)

	d.bagian("III. Riwayat Sertifikasi")
	d.daftarSertifikasi(p.Sertifikasi)

	d.tandaTangan(p, sekarang)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Kop: logo + nama instansi, garis bawah
func (d *drh) kop() {
	pdf := d.pdf
	pdf.ImageOptions("logo", kiri, 12, 0, 15, false, fpdf.ImageOptions{}, 0, "")
	pdf.SetXY(kiri+16, 15.5)
	pdf.SetTextColor(19, 41, 75)
	pdf.SetFont("Helvetica", "B", 10)
	pdf.CellFormat(100, 5, "KEPOLISIAN NEGARA REPUBLIK INDONESIA", "", 2, "L", false, 0, "")
	pdf.SetFont("Helvetica", "B", 9)
	pdf.CellFormat(100, 4.5, "STAF SUMBER DAYA MANUSIA", "", 2, "L", false, 0, "")

	// garis bawah kop: hitam, agak tebal, sedikit transparan
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetAlpha(0.75, "Normal")
	pdf.SetLineWidth(0.4)
	pdf.Line(kiri, 29.5, kiri+lebarIsi, 29.5)
	pdf.SetAlpha(1, "Normal")
	pdf.SetLineWidth(0.2)
	pdf.SetY(35)
}

func (d *drh) judul() {
	pdf := d.pdf
	pdf.SetTextColor(15, 23, 42)
	pdf.SetFont("Helvetica", "BU", 13)
	pdf.CellFormat(lebarIsi, 6, "DAFTAR RIWAYAT HIDUP", "", 1, "C", false, 0, "")
	pdf.Ln(6)
}

// Identitas: kotak pas foto di kiri, data di kanan
func (d *drh) identitas(p *model.Personel, sekarang time.Time) {
	pdf := d.pdf
	y0 := pdf.GetY()

	// Kotak pas foto 3x4
	pdf.SetDrawColor(203, 213, 225)
	pdf.SetFillColor(248, 250, 252)
	pdf.Rect(kiri, y0, 30, 40, "FD")
	pdf.SetTextColor(148, 163, 184)
	pdf.SetFont("Helvetica", "", 7.5)
	pdf.SetXY(kiri, y0+16)
	pdf.MultiCell(30, 4, "Pas foto\n3 x 4", "", "C", false)

	pangkat, labelNomor := "-", "NRP"
	if p.Pangkat != nil {
		pangkat = p.Pangkat.Kode + " (" + p.Pangkat.Nama + ")"
	}
	if p.Jenis == "PNS" {
		labelNomor = "NIP"
	}
	satker := "-"
	if p.Satker != nil {
		satker = helper.Rapikan(p.Satker.Nama)
	}
	jabatan, tmt, lama := "Belum ada jabatan aktif", "-", "-"
	if j := p.JabatanSaatIni; j != nil {
		jabatan = helper.Rapikan(j.NamaJabatan)
		if j.StatusJabatan != "DEFINITIF" {
			jabatan += " (" + j.StatusJabatan + ")"
		}
		tmt = tanggalPanjang(j.TmtMulai.Time)
		lama = lamaLengkap(j.TmtMulai.Time, sekarang)
	}
	lahir := tanggalPanjang(p.TanggalLahir.Time)
	if p.TempatLahir != nil && *p.TempatLahir != "" {
		lahir = helper.Rapikan(*p.TempatLahir) + ", " + lahir
	}
	pensiun := p.TanggalLahir.AddDate(usiaBUP, 0, 0)

	baris := [][2]string{
		{"Nama lengkap", strings.ToUpper(p.Nama)},
		{"Pangkat / " + labelNomor, pangkat + " / " + p.NrpNip},
		{"Jenis personel", p.Jenis},
		{"Jabatan", jabatan},
		{"TMT jabatan", tmt},
		{"Lama jabatan", lama},
		{"Satuan kerja", satker},
		{"Tempat, tanggal lahir", lahir},
		{"Usia", fmt.Sprintf("%d Tahun", usia(p.TanggalLahir.Time, sekarang))},
		{"Batas usia pensiun", fmt.Sprintf("%s (usia %d Tahun)", tanggalPanjang(pensiun), usiaBUP)},
		{"Status personel", statusTeks(p.Status)},
	}

	x := kiri + 36
	pdf.SetXY(x, y0)
	for _, b := range baris {
		pdf.SetX(x)
		pdf.SetFont("Helvetica", "", 8.5)
		pdf.SetTextColor(100, 116, 139)
		pdf.CellFormat(36, 4.7, d.tr(b[0]), "", 0, "L", false, 0, "")
		pdf.CellFormat(3, 4.7, ":", "", 0, "L", false, 0, "")
		pdf.SetTextColor(15, 23, 42)
		if b[0] == "Nama lengkap" {
			pdf.SetFont("Helvetica", "B", 9)
		}
		pdf.MultiCell(lebarIsi-36-39, 4.7, d.tr(b[1]), "", "L", false)
	}
	if pdf.GetY() < y0+42 {
		pdf.SetY(y0 + 42)
	}
	pdf.Ln(4)
}

// Judul bagian: pita biru dongker
func (d *drh) bagian(judul string) {
	pdf := d.pdf
	if pdf.GetY() > 260 {
		pdf.AddPage()
	}
	pdf.SetFillColor(19, 41, 75)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Helvetica", "B", 9)
	pdf.CellFormat(lebarIsi, 6.5, "  "+d.tr(judul), "", 1, "L", true, 0, "")
}

// Daftar jabatan sederhana: "1. Kasatresnarkoba Polres Metro Jakarta Timur (2011)".
// Jabatan saat ini di baris pertama dan dicetak tebal.
func (d *drh) daftarJabatan(riwayat []model.RiwayatJabatan) {
	var baris []string
	tebal := -1
	for i, rj := range riwayat {
		teks := helper.Rapikan(rj.NamaJabatan)
		if rj.StatusJabatan != "DEFINITIF" {
			teks += " (" + rj.StatusJabatan + ")"
		}
		baris = append(baris, fmt.Sprintf("%s (%d)", teks, rj.TmtMulai.Year()))
		if rj.TmtSelesai == nil && tebal < 0 {
			tebal = i
		}
	}
	d.daftar(baris, tebal, "Belum ada riwayat jabatan.")
}

// daftarSertifikasi: "1. Sertifikasi Penyidik, Lemdiklat Polri (2008)", terbaru di atas.
func (d *drh) daftarSertifikasi(list []model.Kualifikasi) {
	var baris []string
	for _, k := range list {
		teks := "-"
		if k.Pendidikan != nil {
			teks = k.Pendidikan.Nama
		}
		if k.Institusi != nil && *k.Institusi != "" {
			teks += ", " + helper.Rapikan(*k.Institusi)
		}
		baris = append(baris, fmt.Sprintf("%s (%d)", teks, k.TahunLulus))
	}
	d.daftar(baris, -1, "Belum ada data sertifikasi.")
}

// daftar: list bernomor sederhana. Baris ke-tebal (indeks) dicetak tebal, -1 = tidak ada.
func (d *drh) daftar(baris []string, tebal int, kosong string) {
	pdf := d.pdf
	pdf.Ln(1.5)
	if len(baris) == 0 {
		pdf.SetFont("Helvetica", "I", 8.5)
		pdf.SetTextColor(100, 116, 139)
		pdf.CellFormat(lebarIsi, 5, d.tr(kosong), "", 1, "L", false, 0, "")
		return
	}
	for i, teks := range baris {
		gaya := ""
		if i == tebal {
			gaya = "B"
		}
		pdf.SetTextColor(100, 116, 139)
		pdf.SetFont("Helvetica", "", 8.5)
		pdf.CellFormat(8, 5, fmt.Sprintf("%d.", i+1), "", 0, "R", false, 0, "")
		pdf.CellFormat(2, 5, "", "", 0, "L", false, 0, "")
		pdf.SetTextColor(15, 23, 42)
		pdf.SetFont("Helvetica", gaya, 8.5)
		pdf.MultiCell(lebarIsi-10, 5, d.tr(teks), "", "L", false)
	}
}

func (d *drh) ringkasan(riwayat []model.RiwayatJabatan, jumlahSertifikasi int, sekarang time.Time) {
	pdf := d.pdf
	var baris [][2]string
	if n := len(riwayat); n > 0 {
		pertama := riwayat[n-1] // riwayat sudah urut terbaru → terlama
		baris = append(baris,
			[2]string{"Jabatan pertama", fmt.Sprintf("%s (TMT %s)", helper.Rapikan(pertama.NamaJabatan), pertama.TmtMulai.Format("02-01-2006"))},
			[2]string{"Masa dinas", lamaLengkap(pertama.TmtMulai.Time, sekarang)}, // dihitung dari jabatan pertama
		)
	}
	sertifikasi := "Belum ada"
	if jumlahSertifikasi > 0 {
		sertifikasi = fmt.Sprintf("%d Sertifikasi", jumlahSertifikasi)
	}
	baris = append(baris, [2]string{"Sertifikasi", sertifikasi})
	pdf.Ln(1.5)
	for _, b := range baris {
		pdf.SetFont("Helvetica", "", 8.5)
		pdf.SetTextColor(100, 116, 139)
		pdf.CellFormat(36, 5, d.tr(b[0]), "", 0, "L", false, 0, "")
		pdf.CellFormat(3, 5, ":", "", 0, "L", false, 0, "")
		pdf.SetTextColor(15, 23, 42)
		pdf.MultiCell(lebarIsi-39, 5, d.tr(b[1]), "", "L", false)
	}
}

// Tanda tangan yang bersangkutan, kanan bawah
func (d *drh) tandaTangan(p *model.Personel, sekarang time.Time) {
	pdf := d.pdf
	const tinggiTTD = 37.0
	if pdf.GetY()+tinggiTTD > akhirIsi {
		pdf.AddPage()
	}
	pdf.Ln(4)
	x := kiri + lebarIsi - 75
	baris := func(teks, gaya string, h float64) {
		pdf.SetX(x)
		pdf.SetFont("Helvetica", gaya, 8.5)
		pdf.CellFormat(75, h, d.tr(teks), "", 1, "C", false, 0, "")
	}
	pdf.SetTextColor(15, 23, 42)
	baris(tanggalPanjang(sekarang), "", 5)
	baris("Yang bersangkutan,", "", 5)
	pdf.Ln(14)
	baris(strings.ToUpper(p.Nama), "BU", 5)
	nomor := "NRP " + p.NrpNip
	pangkat := ""
	if p.Jenis == "PNS" {
		nomor = "NIP " + p.NrpNip
	}
	if p.Pangkat != nil {
		pangkat = strings.ToUpper(p.Pangkat.Nama) + " "
	}
	baris(pangkat+nomor, "", 4.5)
}

// ---- format tanggal & durasi ----

func tanggalPanjang(t time.Time) string {
	return fmt.Sprintf("%02d %s %d", t.Day(), bulan[t.Month()-1], t.Year())
}

func usia(lahir, sekarang time.Time) int {
	th, _, _ := selisih(lahir, sekarang)
	return th
}

// selisih dalam tahun, bulan, hari
func selisih(a, b time.Time) (th, bl, hr int) {
	th, bl, hr = b.Year()-a.Year(), int(b.Month()-a.Month()), b.Day()-a.Day()
	if hr < 0 {
		bl--
		hr += time.Date(b.Year(), b.Month(), 0, 0, 0, 0, 0, time.UTC).Day()
	}
	if bl < 0 {
		th--
		bl += 12
	}
	return
}

func lamaLengkap(a, b time.Time) string {
	th, bl, hr := selisih(a, b)
	var s []string
	if th > 0 {
		s = append(s, fmt.Sprintf("%d Tahun", th))
	}
	if bl > 0 {
		s = append(s, fmt.Sprintf("%d Bulan", bl))
	}
	if hr > 0 || len(s) == 0 {
		s = append(s, fmt.Sprintf("%d Hari", hr))
	}
	return strings.Join(s, " ")
}

func statusTeks(s string) string {
	m := map[string]string{"AKTIF": "Aktif", "PENSIUN": "Pensiun", "MENINGGAL": "Meninggal", "DIBERHENTIKAN": "Diberhentikan", "MUTASI_KELUAR": "Mutasi keluar"}
	if v, ok := m[s]; ok {
		return v
	}
	return s
}
