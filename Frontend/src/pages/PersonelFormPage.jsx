import { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { useToast } from "../components/feedback";
import { Alert, Button, Card, Field, Input, Select } from "../components/ui";
import SatkerPicker from "../components/SatkerPicker";
import { IconArrowLeft } from "../components/icons";
import { STATUS_PERSONEL, tanggalInput } from "../lib/format";
import { useReferensi } from "../lib/useReferensi";

const KOSONG = {
  jenis: "POLRI",
  nrp_nip: "",
  nama: "",
  pangkat_id: "",
  satker_id: "",
  tempat_lahir: "",
  tanggal_lahir: "",
  status: "AKTIF",
};

export default function PersonelFormPage() {
  const { id } = useParams();
  const ubah = Boolean(id);
  const { request } = useAuth();
  const toast = useToast();
  const navigate = useNavigate();
  const pangkat = useReferensi("/referensi/pangkat");
  const satker = useReferensi("/referensi/satker");

  const [v, setV] = useState(KOSONG);
  const [errors, setErrors] = useState({});
  const [error, setError] = useState("");
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    if (!ubah) return;
    request(`/personel/${id}`)
      .then((p) =>
        setV({
          jenis: p.jenis,
          nrp_nip: p.nrp_nip,
          nama: p.nama,
          pangkat_id: p.pangkat_id,
          satker_id: p.satker_id,
          tempat_lahir: p.tempat_lahir ?? "",
          tanggal_lahir: tanggalInput(p.tanggal_lahir),
          status: p.status,
        })
      )
      .catch((err) => setError(err.message));
  }, [id, ubah, request]);

  const set = (k) => (e) => setV({ ...v, [k]: e.target.value });

  async function submit(e) {
    e.preventDefault();
    setSaving(true);
    setErrors({});
    setError("");
    try {
      const p = await request(ubah ? `/personel/${id}` : "/personel", {
        method: ubah ? "PUT" : "POST",
        body: {
          ...v,
          nrp_nip: v.nrp_nip.trim(),
          pangkat_id: Number(v.pangkat_id),
          satker_id: Number(v.satker_id),
          tempat_lahir: v.tempat_lahir.trim() || null,
        },
      });
      toast(ubah ? "Data personel diperbarui" : "Personel ditambahkan");
      navigate(`/personel/${p.id}`, { replace: true });
    } catch (err) {
      setErrors(err.fields);
      setError(err.message);
    } finally {
      setSaving(false);
    }
  }

  const kembali = ubah ? `/personel/${id}` : "/personel";

  return (
    <div className="mx-auto max-w-3xl space-y-4">
      <Link to={kembali} className="inline-flex items-center gap-1 text-sm text-slate-500 hover:text-slate-800">
        <IconArrowLeft className="size-4" /> Kembali
      </Link>

      <form onSubmit={submit} className="space-y-4">
        <Alert>{error}</Alert>

        <Card title="Identitas" subtitle="Data dasar personel">
          <div className="grid gap-4 p-5 sm:grid-cols-2">
            <Field label="Jenis" error={errors.jenis}>
              <Select value={v.jenis} onChange={(e) => setV({ ...v, jenis: e.target.value, pangkat_id: "" })}>
                <option value="POLRI">POLRI</option>
                <option value="PNS">PNS</option>
              </Select>
            </Field>
            <Field label={v.jenis === "POLRI" ? "NRP" : "NIP"} error={errors.nrp_nip} hint={v.jenis === "POLRI" ? "8 digit angka" : "18 digit angka"}>
              <Input value={v.nrp_nip} onChange={set("nrp_nip")} inputMode="numeric" required />
            </Field>
            <div className="sm:col-span-2">
              <Field label="Nama lengkap" error={errors.nama}>
                <Input value={v.nama} onChange={set("nama")} required />
              </Field>
            </div>
            <Field label="Tempat lahir" error={errors.tempat_lahir}>
              <Input value={v.tempat_lahir} onChange={set("tempat_lahir")} />
            </Field>
            <Field label="Tanggal lahir" error={errors.tanggal_lahir}>
              <Input type="date" value={v.tanggal_lahir} onChange={set("tanggal_lahir")} required />
            </Field>
          </div>
        </Card>

        <Card title="Kepegawaian" subtitle="Pangkat, satuan kerja, dan status">
          <div className="grid gap-4 p-5 sm:grid-cols-2">
            <Field label={v.jenis === "POLRI" ? "Pangkat" : "Golongan"} error={errors.pangkat_id}>
              <Select value={v.pangkat_id} onChange={set("pangkat_id")} required>
                <option value="">Pilih</option>
                {pangkat.filter((pg) => pg.jenis === v.jenis).map((pg) => (
                  <option key={pg.id} value={pg.id}>{pg.kode} ({pg.nama})</option>
                ))}
              </Select>
            </Field>
            <Field label="Status" error={errors.status}>
              <Select value={v.status} onChange={set("status")}>
                {Object.entries(STATUS_PERSONEL).map(([k, label]) => <option key={k} value={k}>{label}</option>)}
              </Select>
            </Field>
            <div className="sm:col-span-2">
              <Field label="Satuan kerja" error={errors.satker_id}>
                <SatkerPicker options={satker} value={v.satker_id} onChange={(sid) => setV({ ...v, satker_id: sid })} />
              </Field>
            </div>
          </div>
        </Card>

        <div className="flex justify-end gap-2">
          <Link to={kembali}><Button type="button" variant="outline">Batal</Button></Link>
          <Button type="submit" disabled={saving}>{saving ? "Menyimpan..." : ubah ? "Simpan perubahan" : "Simpan"}</Button>
        </div>
      </form>
    </div>
  );
}
