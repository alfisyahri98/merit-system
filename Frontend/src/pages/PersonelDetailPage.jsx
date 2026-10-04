import { useCallback, useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { useConfirm, useToast } from "../components/feedback";
import { Alert, Avatar, Badge, Button, Card, Field, IconButton, Input, Modal, Select, StatusBadge } from "../components/ui";
import SatkerPicker from "../components/SatkerPicker";
import { IconArrowLeft, IconDownload, IconPencil, IconPlus, IconTrash } from "../components/icons";
import { formatTanggal, labelSatker, lamaJabatan, letakSatker, rapikan, sehariSebelum, tanggalInput } from "../lib/format";
import { useReferensi } from "../lib/useReferensi";

export default function PersonelDetailPage() {
  const daftarSatker = useReferensi("/referensi/satker?semua=1");
  const { id } = useParams();
  const { request, download } = useAuth();
  const toast = useToast();
  const confirm = useConfirm();
  const navigate = useNavigate();
  const [p, setP] = useState(null);
  const [error, setError] = useState("");
  const [form, setForm] = useState(null); // null = tertutup, {} = tambah, objek riwayat = ubah
  const [mengunduh, setMengunduh] = useState(false);

  // PDF Daftar Riwayat Hidup dibuat di backend
  async function unduhDRH() {
    setMengunduh(true);
    try {
      await download(`/personel/${id}/drh`, `DRH_${p?.nrp_nip ?? id}.pdf`);
    } catch (err) {
      toast(err.message, "error");
    } finally {
      setMengunduh(false);
    }
  }

  const muat = useCallback(() => {
    request(`/personel/${id}`).then(setP).catch((err) => setError(err.message));
  }, [id, request]);

  useEffect(muat, [muat]);

  async function hapusPersonel() {
    const ok = await confirm({
      title: `Hapus ${p.nama}?`,
      message: "Data personel beserta seluruh riwayat jabatannya akan dihapus.",
      confirmLabel: "Hapus",
      danger: true,
    });
    if (!ok) return;
    try {
      await request(`/personel/${id}`, { method: "DELETE" });
      toast("Personel dihapus");
      navigate("/personel", { replace: true });
    } catch (err) {
      toast(err.message, "error");
    }
  }

  async function hapusRiwayat(rj) {
    const ok = await confirm({ title: "Hapus riwayat jabatan?", message: rj.nama_jabatan, confirmLabel: "Hapus", danger: true });
    if (!ok) return;
    try {
      await request(`/riwayat-jabatan/${rj.id}`, { method: "DELETE" });
      toast("Riwayat jabatan dihapus");
      muat();
    } catch (err) {
      toast(err.message, "error");
    }
  }

  if (error && !p) {
    return (
      <div className="space-y-4">
        <Link to="/personel" className="inline-flex items-center gap-1 text-sm text-slate-500 hover:text-slate-800"><IconArrowLeft className="size-4" /> Kembali</Link>
        <Alert>{error}</Alert>
      </div>
    );
  }
  if (!p) return <p className="text-sm text-slate-500">Memuat...</p>;

  const jabatan = p.jabatan_saat_ini;
  const riwayat = p.riwayat_jabatan ?? [];
  const sertifikasi = p.sertifikasi ?? [];
  const identitas = [
    ["NRP / NIP", <span className="font-mono">{p.nrp_nip}</span>],
    ["Jenis", p.jenis],
    ["Pangkat", p.pangkat ? `${p.pangkat.kode} (${p.pangkat.nama})` : "-"],
    ["Satker", <>{labelSatker(p.satker)}<span className="block text-xs font-normal text-slate-500">{letakSatker(p.satker, daftarSatker)}</span></>],
    ["Tempat lahir", p.tempat_lahir ?? "-"],
    ["Tanggal lahir", formatTanggal(p.tanggal_lahir)],
  ];

  return (
    <div className="space-y-5">
      <Link to="/personel" className="inline-flex items-center gap-1 text-sm text-slate-500 hover:text-slate-800"><IconArrowLeft className="size-4" /> Data personel</Link>

      {/* Kepala profil */}
      <Card>
        <div className="flex flex-wrap items-center gap-5 p-5">
          <Avatar nama={p.nama} size="lg" />
          <div className="min-w-0 flex-1">
            <div className="flex flex-wrap items-center gap-2">
              <h2 className="text-xl font-semibold tracking-tight">{p.nama}</h2>
              <StatusBadge status={p.status} />
            </div>
            <p className="mt-1 text-sm text-slate-500">
              {p.pangkat?.kode} · <span className="font-mono">{p.nrp_nip}</span> · {rapikan(p.satker?.nama)}
            </p>
            <p className="mt-2 text-sm">
              {jabatan ? (
                <>
                  <span className="font-medium">{rapikan(jabatan.nama_jabatan)}</span>
                  <span className="text-slate-500"> · sejak {formatTanggal(jabatan.tmt_mulai)} ({lamaJabatan(jabatan.tmt_mulai)})</span>
                </>
              ) : (
                <span className="text-slate-500">Belum ada jabatan aktif</span>
              )}
            </p>
          </div>
          <div className="flex flex-wrap gap-2">
            <Button variant="outline" onClick={unduhDRH} disabled={mengunduh}>
              <IconDownload className="size-4" /> {mengunduh ? "Menyiapkan..." : "Unduh data diri"}
            </Button>
            <Link to={`/personel/${id}/ubah`}><Button variant="outline"><IconPencil className="size-4" /> Ubah</Button></Link>
            <Button variant="danger" onClick={hapusPersonel}><IconTrash className="size-4" /> Hapus</Button>
          </div>
        </div>
      </Card>

      <div className="grid gap-5 lg:grid-cols-3">
        <div className="h-fit space-y-5">
          <Card title="Identitas">
            <dl className="divide-y divide-slate-100 text-sm">
              {identitas.map(([k, v]) => (
                <div key={k} className="flex justify-between gap-4 px-5 py-3">
                  <dt className="text-slate-500">{k}</dt>
                  <dd className="text-right font-medium">{v}</dd>
                </div>
              ))}
            </dl>
          </Card>

          <Card title="Riwayat Sertifikasi" subtitle={`${sertifikasi.length} sertifikasi / pelatihan`}>
            {sertifikasi.length === 0 ? (
              <p className="px-5 py-6 text-center text-sm text-slate-500">Belum ada data sertifikasi.</p>
            ) : (
              <ul className="divide-y divide-slate-100 text-sm">
                {sertifikasi.map((k) => (
                  <li key={k.id} className="flex items-start justify-between gap-3 px-5 py-3">
                    <div className="min-w-0">
                      <p className="font-medium">{k.pendidikan?.nama}</p>
                      <p className="text-xs text-slate-500">{rapikan(k.institusi) || "-"}</p>
                    </div>
                    <span className="shrink-0 text-xs text-slate-500 tabular-nums">{k.tahun_lulus}</span>
                  </li>
                ))}
              </ul>
            )}
          </Card>
        </div>

        {/* Timeline karier, urut dari jabatan pertama sampai saat ini */}
        <Card
          title="Riwayat Jabatan"
          subtitle={`${riwayat.length} jabatan, urut dari yang pertama`}
          actions={<Button size="sm" onClick={() => setForm({})}><IconPlus className="size-4" /> Tambah jabatan</Button>}
          className="min-w-0 lg:col-span-2"
        >
          {riwayat.length === 0 ? (
            <p className="px-5 py-10 text-center text-sm text-slate-500">Belum ada riwayat jabatan.</p>
          ) : (
            <ol className="px-5 py-4">
              {riwayat.map((rj, i) => {
                const aktif = !rj.tmt_selesai;
                const terakhir = i === riwayat.length - 1;
                return (
                  <li key={rj.id} className="relative flex gap-4 pb-5 last:pb-0">
                    {!terakhir && <span className="absolute left-[7px] top-5 h-full w-px bg-slate-200" />}
                    <span className={`relative mt-1.5 size-[15px] shrink-0 rounded-full border-2 ${aktif ? "border-dongker bg-dongker" : "border-slate-300 bg-white"}`} />
                    <div className="group min-w-0 flex-1 rounded-lg px-3 py-2 hover:bg-slate-50">
                      <div className="flex items-start justify-between gap-3">
                        <div className="min-w-0">
                          <p className="font-medium">{rapikan(rj.nama_jabatan)}</p>
                          <p className="text-sm text-slate-500">
                            {rapikan(rj.satker?.nama)}
                            {rj.fungsi && ` · ${rj.fungsi.nama}`}
                            {rj.nivelering && ` · ${rj.nivelering.kode}`}
                          </p>
                          <p className="mt-1 flex flex-wrap items-center gap-2 text-xs text-slate-500">
                            <span>
                              {formatTanggal(rj.tmt_mulai)} – {aktif ? "sekarang" : formatTanggal(rj.tmt_selesai)} · {lamaJabatan(rj.tmt_mulai, rj.tmt_selesai)}
                            </span>
                            {aktif && <Badge tone="blue">Saat ini</Badge>}
                            {rj.status_jabatan !== "DEFINITIF" && <Badge tone="amber">{rj.status_jabatan}</Badge>}
                          </p>
                          {rj.nomor_skep && <p className="mt-1 text-xs text-slate-400">{rj.nomor_skep}</p>}
                        </div>
                        <div className="flex shrink-0 gap-1 opacity-60 group-hover:opacity-100">
                          <IconButton label="Ubah" onClick={() => setForm(rj)}><IconPencil className="size-4" /></IconButton>
                          <IconButton label="Hapus" onClick={() => hapusRiwayat(rj)} className="hover:text-red-600"><IconTrash className="size-4" /></IconButton>
                        </div>
                      </div>
                    </div>
                  </li>
                );
              })}
            </ol>
          )}
        </Card>
      </div>

      {form && (
        <RiwayatForm
          personel={p}
          data={form}
          onClose={() => setForm(null)}
          onSaved={(pesan) => {
            setForm(null);
            toast(pesan);
            muat();
          }}
        />
      )}
    </div>
  );
}

function RiwayatForm({ personel, data, onClose, onSaved }) {
  const { request } = useAuth();
  const satker = useReferensi("/referensi/satker?semua=1");
  const fungsi = useReferensi("/referensi/fungsi");
  const nivelering = useReferensi("/referensi/nivelering");
  const ubah = Boolean(data.id);

  const [v, setV] = useState({
    nama_jabatan: data.nama_jabatan ?? "",
    satker_id: data.satker_id ?? personel.satker_id ?? "",
    fungsi_id: data.fungsi_id ?? "",
    nivelering_id: data.nivelering_id ?? "",
    tmt_mulai: tanggalInput(data.tmt_mulai),
    tmt_selesai: tanggalInput(data.tmt_selesai),
    status_jabatan: data.status_jabatan ?? "DEFINITIF",
    nomor_skep: data.nomor_skep ?? "",
    keterangan: data.keterangan ?? "",
  });
  const [errors, setErrors] = useState({});
  const [error, setError] = useState("");
  const [saving, setSaving] = useState(false);
  const set = (k) => (e) => setV({ ...v, [k]: e.target.value });

  // Info mutasi: jabatan definitif aktif yang akan ditutup otomatis oleh backend
  const aktif = personel.jabatan_saat_ini;
  const akanMenutup =
    !ubah && aktif && aktif.status_jabatan === "DEFINITIF" && v.status_jabatan === "DEFINITIF" && !v.tmt_selesai && v.tmt_mulai;

  async function submit(e) {
    e.preventDefault();
    setSaving(true);
    setErrors({});
    setError("");
    const angka = (x) => (x === "" || x === null ? null : Number(x));
    const teks = (x) => (x.trim() === "" ? null : x.trim());
    try {
      await request(ubah ? `/riwayat-jabatan/${data.id}` : `/personel/${personel.id}/riwayat-jabatan`, {
        method: ubah ? "PUT" : "POST",
        body: {
          nama_jabatan: v.nama_jabatan,
          satker_id: angka(v.satker_id),
          fungsi_id: angka(v.fungsi_id),
          nivelering_id: angka(v.nivelering_id),
          tmt_mulai: v.tmt_mulai,
          tmt_selesai: v.tmt_selesai || null,
          status_jabatan: v.status_jabatan,
          nomor_skep: teks(v.nomor_skep),
          keterangan: teks(v.keterangan),
        },
      });
      onSaved(ubah ? "Riwayat jabatan diperbarui" : "Jabatan ditambahkan");
    } catch (err) {
      setErrors(err.fields);
      setError(err.message);
    } finally {
      setSaving(false);
    }
  }

  return (
    <Modal title={ubah ? "Ubah Riwayat Jabatan" : "Tambah Jabatan"} onClose={onClose} wide>
      <form onSubmit={submit} className="space-y-4">
        <Alert>{error}</Alert>
        <Field label="Nama jabatan" error={errors.nama_jabatan}>
          <Input value={v.nama_jabatan} onChange={set("nama_jabatan")} placeholder="contoh: Kanit 2 Satreskrim Polres Metro Jakarta Barat" required />
        </Field>
        <Field label="Satker" error={errors.satker_id}>
          <SatkerPicker options={satker} value={v.satker_id} onChange={(id) => setV({ ...v, satker_id: id })} />
        </Field>
        <div className="grid gap-4 sm:grid-cols-3">
          <Field label="Status jabatan" error={errors.status_jabatan}>
            <Select value={v.status_jabatan} onChange={set("status_jabatan")}>
              <option value="DEFINITIF">Definitif</option>
              <option value="PLT">Plt</option>
              <option value="PLH">Plh</option>
            </Select>
          </Field>
          <Field label="Fungsi" error={errors.fungsi_id}>
            <Select value={v.fungsi_id} onChange={set("fungsi_id")}>
              <option value="">–</option>
              {fungsi.map((f) => <option key={f.id} value={f.id}>{f.nama}</option>)}
            </Select>
          </Field>
          <Field label="Nivelering" error={errors.nivelering_id}>
            <Select value={v.nivelering_id} onChange={set("nivelering_id")}>
              <option value="">–</option>
              {nivelering.map((n) => <option key={n.id} value={n.id}>{n.kode}</option>)}
            </Select>
          </Field>
          <Field label="TMT mulai" error={errors.tmt_mulai}>
            <Input type="date" value={v.tmt_mulai} onChange={set("tmt_mulai")} required />
          </Field>
          <Field label="TMT selesai" error={errors.tmt_selesai} hint="Kosongkan bila masih menjabat">
            <Input type="date" value={v.tmt_selesai} onChange={set("tmt_selesai")} />
          </Field>
          <Field label="Nomor Skep" error={errors.nomor_skep}>
            <Input value={v.nomor_skep} onChange={set("nomor_skep")} />
          </Field>
        </div>
        <Field label="Keterangan" error={errors.keterangan}>
          <Input value={v.keterangan} onChange={set("keterangan")} />
        </Field>

        {akanMenutup && (
          <div className="rounded-lg border border-blue-200 bg-blue-50 px-3 py-2 text-sm text-blue-800">
            Jabatan saat ini <strong>{aktif.nama_jabatan}</strong> akan otomatis ditutup per{" "}
            <strong>{formatTanggal(sehariSebelum(v.tmt_mulai))}</strong>.
          </div>
        )}

        <div className="flex justify-end gap-2 border-t border-slate-100 pt-4">
          <Button type="button" variant="outline" onClick={onClose}>Batal</Button>
          <Button type="submit" disabled={saving}>{saving ? "Menyimpan..." : "Simpan"}</Button>
        </div>
      </form>
    </Modal>
  );
}
