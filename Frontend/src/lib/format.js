export function formatTanggal(iso) {
  if (!iso) return "-";
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return "-";
  return d.toLocaleDateString("id-ID", { day: "numeric", month: "short", year: "numeric" });
}

export const STATUS_PERSONEL = {
  AKTIF: "Aktif",
  PENSIUN: "Pensiun",
  MENINGGAL: "Meninggal",
  DIBERHENTIKAN: "Diberhentikan",
  MUTASI_KELUAR: "Mutasi keluar",
};

export const NAMA_ROLE = { ADMIN_SSDM: "Admin SSDM", OPERATOR: "Operator" };

export const SCOPES = [
  "personel:read", "personel:write", "personel:delete",
  "jabatan:read", "jabatan:write", "jabatan:delete",
  "pangkat:read", "pangkat:write",
  "kualifikasi:read", "kualifikasi:write", "kualifikasi:delete",
];

// "1999-07-24" untuk input type=date (backend sudah kirim format ini)
export function tanggalInput(v) {
  return v ? String(v).slice(0, 10) : "";
}

// Data satker/jabatan dari sumber campur huruf besar-kecil ("SUBBIDTEKINFO BIDTIK",
// "KASATRESNARKOBA POLRESTABES MEDAN"). Ditampilkan dengan gaya penulisan resmi:
// "Subbidtekinfo Bidtik", "Kasatresnarkoba Polrestabes Medan". Data asli tidak diubah.
const SINGKATAN = new Set(["SPKT", "SPN", "SSDM", "SDM", "PNS", "TNI"]);
const ROMAWI = /^(I|II|III|IV|V|VI|VII|VIII|IX|X|XI|XII)$/;

export function rapikan(teks) {
  if (!teks) return teks ?? "";
  return teks
    .split(/\s+/)
    .map((kata) => {
      if (!/^[A-Z0-9-]{2,}$/.test(kata) || !/[A-Z]/.test(kata)) return kata; // sudah campuran / angka
      if (ROMAWI.test(kata) || SINGKATAN.has(kata)) return kata;
      if (kata === "DAN") return "dan";
      return kata
        .split("-")
        .map((k) => k.charAt(0) + k.slice(1).toLowerCase())
        .join("-");
    })
    .join(" ");
}

export function labelSatker(s) {
  return s ? rapikan(s.nama) : "-";
}

// Keterangan letak satker yang bisa dibaca, mis. Polsek Cakung →
// "Polres Metro Jakarta Timur · Polda Metro Jaya". Pakai daftar /referensi/satker.
const PUNCAK = new Set(["POLDA", "SATKER_MABES", "MABES"]);
const petaCache = new WeakMap();

export function letakSatker(satker, daftar) {
  if (!satker || !daftar?.length) return "";
  let peta = petaCache.get(daftar);
  if (!peta) {
    peta = new Map(daftar.map((x) => [x.id, x]));
    petaCache.set(daftar, peta);
  }
  const diri = peta.get(satker.id) ?? satker;
  const sudah = [diri.nama.toLowerCase()];
  const hasil = [];
  let induk = peta.get(diri.parent_id);
  let pertama = true;
  while (induk && induk.tingkat !== "MABES") {
    const nama = induk.nama.toLowerCase();
    const penting = pertama || induk.tingkat === "POLRES" || PUNCAK.has(induk.tingkat);
    if (penting && !sudah.some((x) => x.includes(nama))) {
      hasil.push(rapikan(induk.nama));
      sudah.push(nama);
      pertama = false;
    }
    if (PUNCAK.has(induk.tingkat)) break;
    induk = peta.get(induk.parent_id);
  }
  return hasil.join(" · ");
}

// Lama menjabat, mis. "2 thn 3 bln"
export function lamaJabatan(mulai, selesai) {
  if (!mulai) return "";
  const a = new Date(mulai);
  const b = selesai ? new Date(selesai) : new Date();
  let bulan = (b.getFullYear() - a.getFullYear()) * 12 + (b.getMonth() - a.getMonth());
  if (b.getDate() < a.getDate()) bulan -= 1;
  if (bulan < 1) return "< 1 bln";
  const thn = Math.floor(bulan / 12);
  const bln = bulan % 12;
  return [thn && `${thn} thn`, bln && `${bln} bln`].filter(Boolean).join(" ");
}

// H-1 dari tanggal "YYYY-MM-DD"
export function sehariSebelum(iso) {
  if (!iso) return "";
  const d = new Date(iso + "T00:00:00Z");
  d.setUTCDate(d.getUTCDate() - 1);
  return d.toISOString().slice(0, 10);
}

// Tanggal mencapai usia 58 tahun (batas usia pensiun). 29 Feb → 28 Feb.
export function tanggalPensiun(lahir) {
  if (!lahir) return null;
  const [y, m, d] = lahir.slice(0, 10).split("-").map(Number);
  let t = new Date(Date.UTC(y + 58, m - 1, d));
  if (t.getUTCMonth() !== m - 1) t = new Date(Date.UTC(y + 58, m, 0));
  return t.toISOString().slice(0, 10);
}

// "12 hari lagi", "3 bulan lagi", "1 tahun 2 bulan lagi"
export function sisaWaktu(iso) {
  const hari = Math.round((new Date(iso) - new Date(new Date().toDateString())) / 864e5);
  if (hari <= 0) return "hari ini";
  if (hari < 31) return `${hari} hari lagi`;
  const bulan = Math.floor(hari / 30.44);
  if (bulan < 12) return `${bulan} bulan lagi`;
  const th = Math.floor(bulan / 12);
  return bulan % 12 ? `${th} tahun ${bulan % 12} bulan lagi` : `${th} tahun lagi`;
}
