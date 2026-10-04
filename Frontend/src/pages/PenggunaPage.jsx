import { useCallback, useEffect, useState } from "react";
import { useAuth } from "../context/AuthContext";
import { useConfirm, useToast } from "../components/feedback";
import { Alert, Avatar, Badge, Button, Card, Field, Input, Modal, Pagination, Select, Table } from "../components/ui";
import SatkerPicker from "../components/SatkerPicker";
import { IconPlus } from "../components/icons";
import { NAMA_ROLE } from "../lib/format";
import { useReferensi } from "../lib/useReferensi";

const LIMIT = 20;

export default function PenggunaPage() {
  const { request, user: saya } = useAuth();
  const toast = useToast();
  const confirm = useConfirm();
  const satker = useReferensi("/referensi/satker");
  const [page, setPage] = useState(1);
  const [hasil, setHasil] = useState({ items: [], total: 0 });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [tambah, setTambah] = useState(false);
  const [resetUntuk, setResetUntuk] = useState(null);

  const muat = useCallback(() => {
    setLoading(true);
    request(`/users?page=${page}&limit=${LIMIT}`)
      .then(setHasil)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, [page, request]);

  useEffect(muat, [muat]);

  const namaSatker = (id) => satker.find((s) => s.id === id)?.nama ?? `Satker #${id}`;

  async function ubahStatus(u) {
    const aktif = !u.is_active;
    const ok = await confirm({
      title: `${aktif ? "Aktifkan" : "Nonaktifkan"} ${u.username}?`,
      message: aktif ? "Akun bisa login kembali." : "Akun tidak bisa login dan sesi yang sedang berjalan langsung berakhir.",
      confirmLabel: aktif ? "Aktifkan" : "Nonaktifkan",
      danger: !aktif,
    });
    if (!ok) return;
    try {
      await request(`/users/${u.id}/status`, { method: "PATCH", body: { is_active: aktif } });
      toast(`Akun ${u.username} ${aktif ? "diaktifkan" : "dinonaktifkan"}`);
      muat();
    } catch (err) {
      toast(err.message, "error");
    }
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-end">
        <Button onClick={() => setTambah(true)}><IconPlus className="size-4" /> Tambah Pengguna</Button>
      </div>
      <Alert>{error}</Alert>

      <Card>
        <Table head={["Pengguna", "Role", "Wilayah", "Status", ""]} loading={loading} empty={!hasil.items.length && "Belum ada pengguna."}>
          {!loading &&
            hasil.items.map((u) => (
              <tr key={u.id} className="hover:bg-slate-50">
                <td className="px-5 py-3">
                  <div className="flex items-center gap-3">
                    <Avatar nama={u.username.replace(/[._]/g, " ")} size="sm" />
                    <span className="font-medium">{u.username}</span>
                    {u.id === saya?.id && <Badge tone="blue">Anda</Badge>}
                  </div>
                </td>
                <td className="px-5 py-3">{NAMA_ROLE[u.role] ?? u.role}</td>
                <td className="px-5 py-3">{u.satker_id ? namaSatker(u.satker_id) : "Seluruh satker"}</td>
                <td className="px-5 py-3"><Badge tone={u.is_active ? "green" : "gray"}>{u.is_active ? "Aktif" : "Nonaktif"}</Badge></td>
                <td className="whitespace-nowrap px-5 py-3 text-right">
                  <Button variant="ghost" size="sm" onClick={() => setResetUntuk(u)}>Reset password</Button>
                  {u.id !== saya?.id && (
                    <Button variant="ghost" size="sm" onClick={() => ubahStatus(u)}>{u.is_active ? "Nonaktifkan" : "Aktifkan"}</Button>
                  )}
                </td>
              </tr>
            ))}
        </Table>
        {hasil.total > 0 && <Pagination page={page} limit={LIMIT} total={hasil.total} onChange={setPage} />}
      </Card>

      {tambah && (
        <TambahPengguna
          satker={satker}
          onClose={() => setTambah(false)}
          onSaved={(u) => {
            setTambah(false);
            toast(`Akun ${u.username} dibuat`);
            muat();
          }}
        />
      )}
      {resetUntuk && (
        <ResetPassword
          user={resetUntuk}
          onClose={() => setResetUntuk(null)}
          onSaved={() => {
            toast(`Password ${resetUntuk.username} direset`);
            setResetUntuk(null);
          }}
        />
      )}
    </div>
  );
}

function TambahPengguna({ satker, onClose, onSaved }) {
  const { request } = useAuth();
  const [v, setV] = useState({ username: "", password: "", role: "OPERATOR", satker_id: "" });
  const [errors, setErrors] = useState({});
  const [error, setError] = useState("");
  const [saving, setSaving] = useState(false);
  const set = (k) => (e) => setV({ ...v, [k]: e.target.value });

  async function submit(e) {
    e.preventDefault();
    setSaving(true);
    setErrors({});
    setError("");
    try {
      const u = await request("/users", {
        method: "POST",
        body: {
          username: v.username.trim(),
          password: v.password,
          role: v.role,
          satker_id: v.role === "OPERATOR" && v.satker_id ? Number(v.satker_id) : null,
        },
      });
      onSaved(u);
    } catch (err) {
      setErrors(err.fields);
      setError(err.message);
    } finally {
      setSaving(false);
    }
  }

  return (
    <Modal title="Tambah Pengguna" onClose={onClose}>
      <form onSubmit={submit} className="space-y-4">
        <Alert>{error}</Alert>
        <Field label="Username" error={errors.username} hint="Huruf kecil, angka, titik, garis bawah">
          <Input value={v.username} onChange={set("username")} placeholder="opr.jaktim" required />
        </Field>
        <Field label="Password" error={errors.password} hint="Minimal 8 karakter">
          <Input type="password" value={v.password} onChange={set("password")} required />
        </Field>
        <Field label="Role" error={errors.role}>
          <Select value={v.role} onChange={set("role")}>
            <option value="OPERATOR">Operator (terbatas satker)</option>
            <option value="ADMIN_SSDM">Admin SSDM (seluruh satker)</option>
          </Select>
        </Field>
        {v.role === "OPERATOR" && (
          <Field label="Satker" error={errors.satker_id}>
            <SatkerPicker options={satker} value={v.satker_id} onChange={(id) => setV({ ...v, satker_id: id })} />
          </Field>
        )}
        <div className="flex justify-end gap-2 border-t border-slate-100 pt-4">
          <Button type="button" variant="outline" onClick={onClose}>Batal</Button>
          <Button type="submit" disabled={saving}>{saving ? "Menyimpan..." : "Simpan"}</Button>
        </div>
      </form>
    </Modal>
  );
}

function ResetPassword({ user, onClose, onSaved }) {
  const { request } = useAuth();
  const [password, setPassword] = useState("");
  const [errors, setErrors] = useState({});
  const [saving, setSaving] = useState(false);

  async function submit(e) {
    e.preventDefault();
    setSaving(true);
    setErrors({});
    try {
      await request(`/users/${user.id}/reset-password`, { method: "POST", body: { password } });
      onSaved();
    } catch (err) {
      setErrors(Object.keys(err.fields).length ? err.fields : { password: err.message });
    } finally {
      setSaving(false);
    }
  }

  return (
    <Modal title={`Reset password ${user.username}`} onClose={onClose}>
      <form onSubmit={submit} className="space-y-4">
        <Field label="Password baru" error={errors.password} hint="Minimal 8 karakter">
          <Input type="password" value={password} onChange={(e) => setPassword(e.target.value)} autoFocus required />
        </Field>
        <div className="flex justify-end gap-2 border-t border-slate-100 pt-4">
          <Button type="button" variant="outline" onClick={onClose}>Batal</Button>
          <Button type="submit" disabled={saving}>{saving ? "Menyimpan..." : "Simpan"}</Button>
        </div>
      </form>
    </Modal>
  );
}
