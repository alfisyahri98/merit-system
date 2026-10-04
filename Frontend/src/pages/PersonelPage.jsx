import { useEffect, useMemo, useState } from "react";
import { Link, useNavigate, useSearchParams } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { useToast } from "../components/feedback";
import { Alert, Avatar, Button, Card, Pagination, StatusBadge, Table } from "../components/ui";
import SatkerPicker from "../components/SatkerPicker";
import KolomTabel, { PilihanFilter } from "../components/KolomTabel";
import { IconDownload, IconPlus, IconSearch, IconX } from "../components/icons";
import { STATUS_PERSONEL, formatTanggal, labelSatker, letakSatker, rapikan, sisaWaktu, tanggalPensiun } from "../lib/format";
import { useReferensi } from "../lib/useReferensi";

const LIMIT = 20;
// Nama parameter di URL yang dihitung sebagai filter (sort & page tidak termasuk).
const FILTER = ["q", "jenis", "kelompok", "satker_id", "tanpa_jabatan", "status", "akan_pensiun"];

export default function PersonelPage() {
  const { request, download } = useAuth();
  const toast = useToast();
  const navigate = useNavigate();
  const satker = useReferensi("/referensi/satker");
  const pangkat = useReferensi("/referensi/pangkat");
  const [params, setParams] = useSearchParams();
  const [mengunduh, setMengunduh] = useState(false);

  const q = params.get("q") ?? "";
  const status = params.get("status") ?? "";
  const jenis = params.get("jenis") ?? "";
  const satkerId = params.get("satker_id") ?? "";
  const tanpaJabatan = params.get("tanpa_jabatan") === "true";
  const kelompok = params.get("kelompok") ?? "";
  const akanPensiun = params.get("akan_pensiun") === "true";
  const sort = params.get("sort") ?? "";
  const page = Number(params.get("page")) || 1;

  const [hasil, setHasil] = useState({ items: [], total: 0 });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  // Ubah satu parameter → kembali ke halaman 1
  function setParam(key, value) {
    const baru = new URLSearchParams(params);
    value ? baru.set(key, value) : baru.delete(key);
    baru.delete("page");
    setParams(baru, { replace: true });
  }

  // Query yang sama dipakai list dan export
  const query = useMemo(() => {
    const qs = new URLSearchParams();
    for (const k of [...FILTER, "sort"]) if (params.get(k)) qs.set(k, params.get(k));
    return qs.toString();
  }, [params]);

  useEffect(() => {
    let batal = false;
    setLoading(true);
    setError("");
    request(`/personel?${query}&page=${page}&limit=${LIMIT}`)
      .then((d) => !batal && setHasil(d))
      .catch((err) => !batal && setError(err.message))
      .finally(() => !batal && setLoading(false));
    return () => {
      batal = true;
    };
  }, [query, page, request]);

  async function exportExcel() {
    setMengunduh(true);
    try {
      await download(`/export/personel?${query}`, akanPensiun ? "personel_akan_pensiun.xlsx" : "personel.xlsx");
      toast(`${hasil.total.toLocaleString("id-ID")} personel diexport ke Excel`);
    } catch (err) {
      toast(err.message, "error");
    } finally {
      setMengunduh(false);
    }
  }

  // Pilihan kelompok pangkat, urut dari terendah sesuai referensi
  const kelompokPangkat = useMemo(() => {
    const per = { POLRI: [], PNS: [] };
    for (const p of pangkat) if (per[p.jenis] && !per[p.jenis].includes(p.kelompok)) per[p.jenis].push(p.kelompok);
    return per;
  }, [pangkat]);

  const satkerTerpilih = satker.find((s) => String(s.id) === satkerId);

  // Filter yang sedang aktif, ditampilkan sebagai chip di atas tabel
  const chip = [
    q && ["q", `Cari: ${q}`],
    jenis && ["jenis", jenis],
    kelompok && ["kelompok", `Pangkat: ${kelompok}`],
    satkerId && ["satker_id", `Satker: ${satkerTerpilih ? labelSatker(satkerTerpilih) : satkerId}`],
    tanpaJabatan && ["tanpa_jabatan", "Belum punya jabatan aktif"],
    status && ["status", `Status: ${STATUS_PERSONEL[status] ?? status}`],
    akanPensiun && ["akan_pensiun", "Pensiun dalam 12 bulan"],
  ].filter(Boolean);

  const kolom = { sort, onSort: (v) => setParam("sort", v) };

  const head = [
    <KolomTabel key="p" label="Personel" sortKey="nama" judulUrut="Urut nama" filterAktif={!!(q || jenis)} {...kolom}>
      <FilterPersonel q={q} jenis={jenis} setParam={setParam} />
    </KolomTabel>,
    <KolomTabel key="pg" label="Pangkat" sortKey="pangkat" arahAwal="desc" judulUrut="Urut pangkat (tertinggi dulu)" filterAktif={!!kelompok} {...kolom}>
      {(tutup) => (
        <div className="max-h-80 space-y-2 overflow-auto">
          <PilihanFilter value={kelompok} onChange={(v) => { setParam("kelompok", v); tutup(); }} options={[["", "Semua pangkat"]]} />
          {Object.entries(kelompokPangkat).map(([j, list]) => list.length > 0 && (
            <div key={j}>
              <p className="px-2 pb-1 text-xs font-medium text-slate-400">{j}</p>
              <PilihanFilter value={kelompok} onChange={(v) => { setParam("kelompok", v); tutup(); }} options={[...list].reverse().map((k) => [k, k])} />
            </div>
          ))}
        </div>
      )}
    </KolomTabel>,
    <KolomTabel key="s" label="Satker" sortKey="satker" judulUrut="Urut nama satker" filterAktif={!!satkerId} {...kolom}>
      {(tutup) => (
        <div className="space-y-2">
          <p className="text-xs text-slate-500">Personel satker ini beserta satker di bawahnya.</p>
          <SatkerPicker options={satker} value={satkerId} onChange={(id) => { setParam("satker_id", id); tutup(); }} emptyLabel="Semua satker" placeholder="Ketik nama satker, mis. jaksel" />
        </div>
      )}
    </KolomTabel>,
    <KolomTabel key="j" label="Jabatan saat ini" filterAktif={tanpaJabatan} {...kolom}>
      <label className="flex items-center gap-2 px-1 py-1">
        <input type="checkbox" checked={tanpaJabatan} onChange={(e) => setParam("tanpa_jabatan", e.target.checked ? "true" : "")} className="accent-dongker" />
        Hanya yang belum punya jabatan aktif
      </label>
    </KolomTabel>,
    <KolomTabel
      key="st"
      label={akanPensiun ? "Pensiun" : "Status"}
      sortKey={akanPensiun ? "pensiun" : undefined}
      judulUrut="Urut tanggal pensiun"
      filterAktif={!!(status || akanPensiun)}
      {...kolom}
    >
      {(tutup) => (
        <div className="space-y-2">
          <label className="flex items-center gap-2 rounded-md bg-amber-50 px-2 py-2 text-amber-800">
            <input type="checkbox" checked={akanPensiun} onChange={(e) => { setParam("akan_pensiun", e.target.checked ? "true" : ""); tutup(); }} className="accent-amber-600" />
            Pensiun dalam 12 bulan
          </label>
          <PilihanFilter
            value={status}
            onChange={(v) => { setParam("status", v); tutup(); }}
            options={[["", "Semua status"], ...Object.entries(STATUS_PERSONEL)]}
          />
        </div>
      )}
    </KolomTabel>,
  ];

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <p className="text-sm text-slate-500">
          {loading ? "Memuat..." : `${hasil.total.toLocaleString("id-ID")} personel ditemukan`}
        </p>
        <div className="flex gap-2">
          <Button variant="outline" onClick={exportExcel} disabled={mengunduh || !hasil.total}>
            <IconDownload className="size-4" /> {mengunduh ? "Menyiapkan..." : "Export Excel"}
          </Button>
          <Link to="/personel/baru">
            <Button><IconPlus className="size-4" /> Tambah Personel</Button>
          </Link>
        </div>
      </div>

      <Card>
        {chip.length > 0 && (
          <div className="flex flex-wrap items-center gap-2 border-b border-slate-100 px-5 py-3">
            {chip.map(([key, label]) => (
              <span key={key} className="inline-flex items-center gap-1 rounded-full bg-dongker/10 py-1 pl-3 pr-1.5 text-xs font-medium text-dongker">
                {label}
                <button onClick={() => setParam(key, "")} aria-label={`Hapus filter ${label}`} className="rounded-full p-0.5 hover:bg-dongker/15">
                  <IconX className="size-3.5" />
                </button>
              </span>
            ))}
            {chip.length > 1 && (
              <button
                onClick={() => setParams(sort ? { sort } : {})}
                className="text-xs text-slate-500 hover:text-slate-800"
              >
                Hapus semua
              </button>
            )}
          </div>
        )}

        {error && <div className="px-4 pt-4"><Alert>{error}</Alert></div>}

        <Table head={head} loading={loading} empty={!hasil.items.length && "Tidak ada personel yang cocok dengan filter."}>
          {!loading &&
            hasil.items.map((p) => (
              <tr key={p.id} onClick={() => navigate(`/personel/${p.id}`)} className="cursor-pointer hover:bg-slate-50">
                <td className="px-5 py-3">
                  <div className="flex items-center gap-3">
                    <Avatar nama={p.nama} size="sm" />
                    <div>
                      <p className="font-medium">{p.nama}</p>
                      <p className="font-mono text-xs text-slate-500">{p.nrp_nip}</p>
                    </div>
                  </div>
                </td>
                <td className="px-5 py-3">{p.pangkat?.kode}</td>
                <td className="px-5 py-3">
                  <p>{rapikan(p.satker?.nama)}</p>
                  <p className="text-xs text-slate-500">{letakSatker(p.satker, satker)}</p>
                </td>
                <td className="max-w-xs px-5 py-3">
                  {p.jabatan_saat_ini ? <span className="line-clamp-2">{rapikan(p.jabatan_saat_ini.nama_jabatan)}</span> : <span className="text-slate-400">Belum ada</span>}
                </td>
                {akanPensiun ? (
                  <td className="whitespace-nowrap px-5 py-3">
                    <p className="font-medium">{formatTanggal(tanggalPensiun(p.tanggal_lahir))}</p>
                    <p className="text-xs text-amber-700">{sisaWaktu(tanggalPensiun(p.tanggal_lahir))}</p>
                  </td>
                ) : (
                  <td className="px-5 py-3"><StatusBadge status={p.status} /></td>
                )}
              </tr>
            ))}
        </Table>
        {hasil.total > 0 && (
          <Pagination page={page} limit={LIMIT} total={hasil.total} onChange={(p) => { const b = new URLSearchParams(params); b.set("page", p); setParams(b); }} />
        )}
      </Card>
    </div>
  );
}

// Isi filter kolom Personel: cari nama/NRP + jenis
function FilterPersonel({ q, jenis, setParam }) {
  const [ketik, setKetik] = useState(q);

  // Cari jalan sendiri 400 ms setelah berhenti mengetik
  useEffect(() => {
    if (ketik.trim() === q) return;
    const t = setTimeout(() => setParam("q", ketik.trim()), 400);
    return () => clearTimeout(t);
  }, [ketik]); // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <div className="space-y-3">
      <div className="relative">
        <IconSearch className="pointer-events-none absolute left-2.5 top-1/2 size-4 -translate-y-1/2 text-slate-400" />
        <input
          autoFocus
          value={ketik}
          onChange={(e) => setKetik(e.target.value)}
          placeholder="Nama atau NRP/NIP"
          className="w-full rounded-lg border border-slate-300 py-2 pl-8 pr-3 focus:border-dongker focus:outline-none focus:ring-2 focus:ring-dongker/15"
        />
      </div>
      <div>
        <p className="px-2 pb-1 text-xs font-medium text-slate-400">Jenis</p>
        <PilihanFilter value={jenis} onChange={(v) => setParam("jenis", v)} options={[["", "POLRI & PNS"], ["POLRI", "POLRI"], ["PNS", "PNS"]]} />
      </div>
    </div>
  );
}
