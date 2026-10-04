import { useCallback, useEffect, useState } from "react";
import { useAuth } from "../context/AuthContext";
import { useConfirm, useToast } from "../components/feedback";
import { Alert, Badge, Button, Card, Field, Input, Modal, Pagination, Table } from "../components/ui";
import SatkerPicker from "../components/SatkerPicker";
import { IconPlus } from "../components/icons";
import { SCOPES } from "../lib/format";
import { useReferensi } from "../lib/useReferensi";

const LIMIT = 20;

export default function AplikasiPage() {
  const { request } = useAuth();
  const toast = useToast();
  const confirm = useConfirm();
  const satker = useReferensi("/referensi/satker");
  const [page, setPage] = useState(1);
  const [hasil, setHasil] = useState({ items: [], total: 0 });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [tambah, setTambah] = useState(false);
  const [secret, setSecret] = useState(null); // { client_id, client_secret }

  const muat = useCallback(() => {
    setLoading(true);
    request(`/api-clients?page=${page}&limit=${LIMIT}`)
      .then(setHasil)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, [page, request]);

  useEffect(muat, [muat]);

  const namaSatker = (id) => satker.find((s) => s.id === id)?.nama ?? `Satker #${id}`;

  async function ubahStatus(c) {
    const aktif = !c.is_active;
    const ok = await confirm({
      title: `${aktif ? "Aktifkan" : "Cabut akses"} ${c.client_id}?`,
      message: aktif ? "Aplikasi bisa meminta token lagi." : "Token yang sedang dipakai aplikasi ini langsung ditolak.",
      confirmLabel: aktif ? "Aktifkan" : "Cabut akses",
      danger: !aktif,
    });
    if (!ok) return;
    try {
      await request(`/api-clients/${c.id}/status`, { method: "PATCH", body: { is_active: aktif } });
      toast(aktif ? "Akses diaktifkan" : "Akses dicabut");
      muat();
    } catch (err) {
      toast(err.message, "error");
    }
  }

  async function gantiSecret(c) {
    const ok = await confirm({
      title: `Ganti secret ${c.client_id}?`,
      message: "Secret lama langsung tidak bisa dipakai untuk meminta token.",
      confirmLabel: "Ganti secret",
    });
    if (!ok) return;
    try {
      const d = await request(`/api-clients/${c.id}/rotate-secret`, { method: "POST" });
      setSecret({ client_id: d.client.client_id, client_secret: d.client_secret });
    } catch (err) {
      toast(err.message, "error");
    }
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-end">
        <Button onClick={() => setTambah(true)}><IconPlus className="size-4" /> Tambah API Client</Button>
      </div>
      <Alert>{error}</Alert>

      <Card>
        <Table head={["Aplikasi", "Scope", "Wilayah", "Status", ""]} loading={loading} empty={!hasil.items.length && "Belum ada API client."}>
          {!loading &&
            hasil.items.map((c) => (
              <tr key={c.id} className="hover:bg-slate-50">
                <td className="px-5 py-3">
                  <p className="font-medium">{c.nama_aplikasi}</p>
                  <p className="font-mono text-xs text-slate-500">{c.client_id}</p>
                </td>
                <td className="px-5 py-3">
                  <div className="flex max-w-xs flex-wrap gap-1">
                    {c.scopes?.map((s) => <Badge key={s}>{s}</Badge>)}
                  </div>
                </td>
                <td className="px-5 py-3">{c.akses_satker_id ? namaSatker(c.akses_satker_id) : "Seluruh satker"}</td>
                <td className="px-5 py-3"><Badge tone={c.is_active ? "green" : "red"}>{c.is_active ? "Aktif" : "Dicabut"}</Badge></td>
                <td className="whitespace-nowrap px-5 py-3 text-right">
                  <Button variant="ghost" size="sm" onClick={() => gantiSecret(c)}>Ganti secret</Button>
                  <Button variant="ghost" size="sm" onClick={() => ubahStatus(c)}>{c.is_active ? "Cabut" : "Aktifkan"}</Button>
                </td>
              </tr>
            ))}
        </Table>
        {hasil.total > 0 && <Pagination page={page} limit={LIMIT} total={hasil.total} onChange={setPage} />}
      </Card>

      {tambah && (
        <TambahClient
          satker={satker}
          onClose={() => setTambah(false)}
          onSaved={(d) => {
            setTambah(false);
            setSecret({ client_id: d.client.client_id, client_secret: d.client_secret });
            muat();
          }}
        />
      )}

      {secret && (
        <Modal title="Client secret" onClose={() => setSecret(null)}>
          <p className="text-sm text-slate-600">
            Simpan secret untuk <strong>{secret.client_id}</strong> sekarang. Secret tidak akan ditampilkan lagi.
          </p>
          <div className="mt-3 flex gap-2">
            <code className="flex-1 break-all rounded-lg bg-slate-100 px-3 py-2 font-mono text-sm">{secret.client_secret}</code>
            <Button
              variant="outline"
              onClick={() => navigator.clipboard.writeText(secret.client_secret).then(() => toast("Secret disalin"))}
            >
              Salin
            </Button>
          </div>
          <div className="mt-4 flex justify-end">
            <Button onClick={() => setSecret(null)}>Sudah disimpan</Button>
          </div>
        </Modal>
      )}
    </div>
  );
}

function TambahClient({ satker, onClose, onSaved }) {
  const { request } = useAuth();
  const [v, setV] = useState({ nama_aplikasi: "", client_id: "", scopes: ["personel:read"], akses_satker_id: "" });
  const [errors, setErrors] = useState({});
  const [error, setError] = useState("");
  const [saving, setSaving] = useState(false);
  const set = (k) => (e) => setV({ ...v, [k]: e.target.value });

  const toggle = (s) => setV({ ...v, scopes: v.scopes.includes(s) ? v.scopes.filter((x) => x !== s) : [...v.scopes, s] });

  async function submit(e) {
    e.preventDefault();
    setSaving(true);
    setErrors({});
    setError("");
    try {
      const d = await request("/api-clients", {
        method: "POST",
        body: {
          nama_aplikasi: v.nama_aplikasi.trim(),
          client_id: v.client_id.trim(),
          scopes: v.scopes,
          akses_satker_id: v.akses_satker_id ? Number(v.akses_satker_id) : null,
        },
      });
      onSaved(d);
    } catch (err) {
      setErrors(err.fields);
      setError(err.message);
    } finally {
      setSaving(false);
    }
  }

  return (
    <Modal title="Tambah API Client" onClose={onClose}>
      <form onSubmit={submit} className="space-y-4">
        <Alert>{error}</Alert>
        <Field label="Nama aplikasi" error={errors.nama_aplikasi}>
          <Input value={v.nama_aplikasi} onChange={set("nama_aplikasi")} required />
        </Field>
        <Field label="Client ID" error={errors.client_id} hint="Huruf kecil, angka, tanda minus">
          <Input value={v.client_id} onChange={set("client_id")} placeholder="app-pmj" required />
        </Field>
        <Field label="Wilayah akses" error={errors.akses_satker_id}>
          <SatkerPicker options={satker} value={v.akses_satker_id} onChange={(id) => setV({ ...v, akses_satker_id: id })} emptyLabel="Seluruh satker" />
        </Field>
        <Field label="Scope" error={errors.scopes} grup>
          <div className="grid grid-cols-2 gap-1.5 rounded-lg border border-slate-200 p-3 text-sm">
            {SCOPES.map((s) => (
              <label key={s} className="flex items-center gap-2">
                <input type="checkbox" checked={v.scopes.includes(s)} onChange={() => toggle(s)} className="accent-dongker" />
                {s}
              </label>
            ))}
          </div>
        </Field>
        <div className="flex justify-end gap-2 border-t border-slate-100 pt-4">
          <Button type="button" variant="outline" onClick={onClose}>Batal</Button>
          <Button type="submit" disabled={saving}>{saving ? "Menyimpan..." : "Simpan"}</Button>
        </div>
      </form>
    </Modal>
  );
}
