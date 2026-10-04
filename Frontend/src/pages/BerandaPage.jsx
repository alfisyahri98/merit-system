import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { useToast } from "../components/feedback";
import { Alert, Avatar, Button, Card } from "../components/ui";
import Donut from "../components/Donut";
import { IconAlert, IconBriefcase, IconDownload, IconPlus, IconSearch, IconUsers } from "../components/icons";
import { formatTanggal, labelSatker, rapikan } from "../lib/format";
import { useReferensi } from "../lib/useReferensi";

const angka = (n) => (typeof n === "number" ? n.toLocaleString("id-ID") : "–");
const TAMPIL_AWAL = 9;

function Stat({ icon: Icon, tone, label, value, to, hint }) {
  const Wrap = to ? Link : "div";
  return (
    <Wrap {...(to ? { to } : {})} className="flex items-center gap-4 rounded-xl border border-slate-200 bg-white p-5 shadow-sm transition hover:border-dongker/30 hover:shadow-md">
      <span className={`grid size-11 shrink-0 place-items-center rounded-lg ${tone}`}>
        <Icon className="size-5" />
      </span>
      <div className="min-w-0">
        <p className="text-sm text-slate-500">{label}</p>
        <p className="text-2xl font-semibold tracking-tight tabular-nums">{value}</p>
        {hint && <p className="text-xs text-slate-500">{hint}</p>}
      </div>
    </Wrap>
  );
}

// Seluruh kartu bisa diklik (ke daftar personel satker),
// badge pensiun punya link sendiri (ke daftar yang akan pensiun).
function KartuSatker({ s }) {
  return (
    <div className="group relative flex flex-col rounded-lg border border-slate-200 p-4 transition hover:border-dongker/40 hover:bg-slate-50">
      <Link
        to={`/personel?satker_id=${s.satker_id}`}
        className="line-clamp-2 min-h-[2.5rem] text-sm font-medium after:absolute after:inset-0 after:rounded-lg group-hover:text-dongker"
      >
        {rapikan(s.nama)}
      </Link>
      <p className="mt-2 text-2xl font-semibold tabular-nums">{angka(s.total)}</p>
      <p className="text-xs text-slate-500">
        {angka(s.polri)} POLRI · {angka(s.pns)} PNS
      </p>
      {s.akan_pensiun > 0 && (
        <Link
          to={`/personel?satker_id=${s.satker_id}&akan_pensiun=true`}
          title="Lihat personel yang akan pensiun"
          className="relative z-10 mt-2 inline-flex w-fit items-center gap-1 rounded-full bg-amber-50 px-2 py-0.5 text-xs font-medium text-amber-700 ring-amber-300 transition hover:bg-amber-100 hover:ring-1"
        >
          <IconAlert className="size-3.5" /> {s.akan_pensiun} pensiun ≤ 12 bln
        </Link>
      )}
    </div>
  );
}

export default function BerandaPage() {
  const { request, download, user } = useAuth();
  const toast = useToast();
  const satker = useReferensi("/referensi/satker");
  const [r, setR] = useState(null);
  const [error, setError] = useState("");
  const [cari, setCari] = useState("");
  const [semua, setSemua] = useState(false);
  const [mengunduh, setMengunduh] = useState(false);

  useEffect(() => {
    request("/dashboard").then(setR).catch((err) => setError(err.message));
  }, [request]);

  const wilayah = user?.satker_id ? labelSatker(satker.find((s) => s.id === user.satker_id)) : "Seluruh satker";

  const daftarSatker = useMemo(() => {
    const s = cari.trim().toLowerCase();
    const list = r?.per_satker ?? [];
    return s ? list.filter((x) => x.nama.toLowerCase().includes(s) || x.kode.toLowerCase().includes(s)) : list;
  }, [r, cari]);
  const tampil = semua || cari ? daftarSatker : daftarSatker.slice(0, TAMPIL_AWAL);

  const pangkat = r?.per_kelompok_pangkat ?? [];
  const polri = pangkat.filter((p) => p.jenis === "POLRI");
  const pns = pangkat.filter((p) => p.jenis === "PNS");

  async function exportExcel() {
    setMengunduh(true);
    try {
      await download("/export/personel", "personel.xlsx");
      toast("File Excel diunduh");
    } catch (err) {
      toast(err.message, "error");
    } finally {
      setMengunduh(false);
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h2 className="text-xl font-semibold tracking-tight">Ringkasan personel</h2>
          <p className="text-sm text-slate-500">Wilayah: {wilayah}</p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" onClick={exportExcel} disabled={mengunduh}>
            <IconDownload className="size-4" /> {mengunduh ? "Menyiapkan..." : "Export Excel"}
          </Button>
          <Link to="/personel/baru">
            <Button><IconPlus className="size-4" /> Tambah Personel</Button>
          </Link>
        </div>
      </div>

      <Alert>{error}</Alert>

      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <Stat icon={IconUsers} tone="bg-blue-50 text-blue-700" label="Total personel" value={angka(r?.total)} to="/personel" />
        <Stat icon={IconUsers} tone="bg-slate-100 text-slate-700" label="Anggota POLRI" value={angka(r?.polri)} to="/personel?jenis=POLRI" />
        <Stat icon={IconBriefcase} tone="bg-slate-100 text-slate-700" label="PNS" value={angka(r?.pns)} to="/personel?jenis=PNS" />
        <Stat icon={IconAlert} tone="bg-amber-50 text-amber-700" label="Pensiun dalam 12 bulan" value={angka(r?.akan_pensiun)} hint="Batas usia 58 tahun" to="/personel?akan_pensiun=true" />
      </div>

      <div className="grid gap-4 xl:grid-cols-3">
        <Card
          title="Personel per satker"
          subtitle={`${daftarSatker.length} satker · klik kartu untuk melihat personelnya`}
          actions={
            <div className="relative">
              <IconSearch className="pointer-events-none absolute left-2.5 top-1/2 size-4 -translate-y-1/2 text-slate-400" />
              <input
                value={cari}
                onChange={(e) => setCari(e.target.value)}
                placeholder="Cari satker"
                className="w-44 rounded-lg border border-slate-200 py-1.5 pl-8 pr-2 text-sm focus:border-dongker focus:outline-none"
              />
            </div>
          }
          className="min-w-0 xl:col-span-2"
        >
          <div className="grid gap-3 p-5 sm:grid-cols-2 lg:grid-cols-3">
            {tampil.map((s) => <KartuSatker key={s.satker_id} s={s} />)}
            {r && !tampil.length && <p className="col-span-full py-6 text-center text-sm text-slate-500">Satker tidak ditemukan.</p>}
            {!r && <p className="col-span-full py-6 text-center text-sm text-slate-500">Memuat...</p>}
          </div>
          {!cari && daftarSatker.length > TAMPIL_AWAL && (
            <div className="border-t border-slate-100 px-5 py-3 text-center">
              <button onClick={() => setSemua(!semua)} className="text-sm font-medium text-dongker hover:underline">
                {semua ? "Tampilkan lebih sedikit" : `Tampilkan semua (${daftarSatker.length})`}
              </button>
            </div>
          )}
        </Card>

        <Card title="Komposisi pangkat" subtitle="Klik untuk melihat personelnya" className="h-fit min-w-0">
          <div className="space-y-6 p-5">
            {polri.length > 0 && <Donut judul="POLRI" data={polri} linkTo={(k) => `/personel?kelompok=${encodeURIComponent(k)}`} />}
            {pns.length > 0 && <Donut judul="PNS" data={pns} linkTo={(k) => `/personel?kelompok=${encodeURIComponent(k)}`} />}
            {r && !pangkat.length && <p className="text-center text-sm text-slate-500">Belum ada data.</p>}
          </div>
        </Card>
      </div>

      <div className="grid gap-4 lg:grid-cols-2">
        <Card
          title="Memasuki pensiun"
          subtitle="Personel aktif yang mencapai usia 58 tahun dalam 12 bulan"
          actions={
            r?.akan_pensiun > 0 && (
              <Link to="/personel?akan_pensiun=true" className="whitespace-nowrap text-sm font-medium text-dongker hover:underline">
                Lihat semua ({angka(r.akan_pensiun)})
              </Link>
            )
          }
          className="min-w-0"
        >
          {r?.akan_pensiun_list?.length ? (
            <ul className="divide-y divide-slate-100">
              {r.akan_pensiun_list.map((p) => (
                <li key={p.id}>
                  <Link to={`/personel/${p.id}`} className="flex items-start gap-3 px-5 py-3 hover:bg-slate-50">
                    <Avatar nama={p.nama} size="sm" />
                    <div className="min-w-0 flex-1">
                      <p className="truncate text-sm font-medium">{p.pangkat} {p.nama}</p>
                      <p className="truncate text-xs text-slate-600">{p.jabatan ? rapikan(p.jabatan) : "Belum ada jabatan aktif"}</p>
                      <p className="truncate text-xs text-slate-400">{rapikan(p.satker)}</p>
                    </div>
                    <span className="shrink-0 text-right text-xs text-slate-500">
                      Pensiun
                      <br />
                      <span className="font-medium text-slate-700">{formatTanggal(p.tanggal_pensiun)}</span>
                    </span>
                  </Link>
                </li>
              ))}
            </ul>
          ) : (
            <p className="px-5 py-8 text-center text-sm text-slate-500">{r ? "Tidak ada personel yang pensiun dalam 12 bulan." : "Memuat..."}</p>
          )}
        </Card>

        <Card title="Mutasi terbaru" subtitle="Jabatan yang paling baru dimulai" className="min-w-0">
          {r?.mutasi_terbaru?.length ? (
            <ul className="divide-y divide-slate-100">
              {r.mutasi_terbaru.map((m, i) => (
                <li key={i}>
                  <Link to={`/personel/${m.personel_id}`} className="flex items-start gap-3 px-5 py-3 hover:bg-slate-50">
                    <Avatar nama={m.nama} size="sm" />
                    <div className="min-w-0 flex-1">
                      <p className="truncate text-sm font-medium">{m.pangkat} {m.nama}</p>
                      <p className="truncate text-xs text-slate-600">{rapikan(m.nama_jabatan)}</p>
                      <p className="truncate text-xs text-slate-400">{rapikan(m.satker)}</p>
                    </div>
                    <span className="shrink-0 text-xs text-slate-400">{formatTanggal(m.tmt_mulai)}</span>
                  </Link>
                </li>
              ))}
            </ul>
          ) : (
            <p className="px-5 py-8 text-center text-sm text-slate-500">{r ? "Belum ada data." : "Memuat..."}</p>
          )}
        </Card>
      </div>
    </div>
  );
}
