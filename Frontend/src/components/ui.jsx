import { IconX } from "./icons";
import { STATUS_PERSONEL } from "../lib/format";

export function Button({ variant = "primary", size = "md", className = "", ...props }) {
  const gaya = {
    primary: "bg-dongker text-white hover:bg-dongker/90 shadow-sm",
    outline: "border border-slate-300 bg-white text-slate-700 hover:bg-slate-50",
    ghost: "text-slate-600 hover:bg-slate-100",
    danger: "border border-red-200 bg-white text-red-700 hover:bg-red-50",
    dangerSolid: "bg-red-600 text-white hover:bg-red-700",
  }[variant];
  const ukuran = size === "sm" ? "px-2.5 py-1.5 text-xs" : "px-3.5 py-2 text-sm";
  return (
    <button
      className={`inline-flex items-center justify-center gap-1.5 rounded-lg font-medium transition disabled:opacity-50 ${ukuran} ${gaya} ${className}`}
      {...props}
    />
  );
}

export function IconButton({ label, className = "", children, ...props }) {
  return (
    <button title={label} aria-label={label} className={`rounded-md p-1.5 text-slate-500 hover:bg-slate-100 hover:text-slate-800 ${className}`} {...props}>
      {children}
    </button>
  );
}

// grup = isinya beberapa pilihan (mis. checkbox), dibungkus fieldset bukan label.
export function Field({ label, error, hint, grup = false, children }) {
  const judul = "mb-1 block text-sm font-medium text-slate-700";
  return (
    <div>
      {grup ? (
        <fieldset>
          <legend className={judul}>{label}</legend>
          {children}
        </fieldset>
      ) : (
        // input di dalam <label> → label otomatis terhubung ke input (klik & screen reader)
        <label className="block">
          <span className={judul}>{label}</span>
          {children}
        </label>
      )}
      {hint && !error && <p className="mt-1 text-xs text-slate-500">{hint}</p>}
      {error && <p className="mt-1 text-xs text-red-600">{error}</p>}
    </div>
  );
}

const inputClass =
  "w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm focus:border-dongker focus:outline-none focus:ring-2 focus:ring-dongker/15";

export function Input(props) {
  return <input className={inputClass} {...props} />;
}

export function Select({ children, ...props }) {
  return (
    <select className={inputClass} {...props}>
      {children}
    </select>
  );
}

export function Alert({ children }) {
  if (!children) return null;
  return <div className="rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">{children}</div>;
}

export function Card({ title, subtitle, actions, children, className = "" }) {
  return (
    <section className={`rounded-xl border border-slate-200 bg-white shadow-sm ${className}`}>
      {(title || actions) && (
        <div className="flex items-center justify-between gap-3 border-b border-slate-100 px-5 py-3.5">
          <div>
            <h2 className="font-semibold">{title}</h2>
            {subtitle && <p className="text-xs text-slate-500">{subtitle}</p>}
          </div>
          <div className="flex gap-2">{actions}</div>
        </div>
      )}
      {children}
    </section>
  );
}

export function PageHeader({ title, subtitle, actions }) {
  return (
    <div className="mb-5 flex flex-wrap items-end justify-between gap-3">
      <div>
        <h1 className="text-xl font-semibold tracking-tight">{title}</h1>
        {subtitle && <p className="mt-0.5 text-sm text-slate-500">{subtitle}</p>}
      </div>
      {actions && <div className="flex gap-2">{actions}</div>}
    </div>
  );
}

export function Table({ head, children, empty, loading }) {
  return (
    <div className="overflow-x-auto">
      <table className="w-full text-left text-sm">
        <thead className="border-b border-slate-100 text-xs uppercase tracking-wide text-slate-500">
          <tr>
            {head.map((h, i) => (
              <th key={i} className="whitespace-nowrap px-5 py-3 font-medium">{h}</th>
            ))}
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-100">{children}</tbody>
      </table>
      {loading && <p className="px-5 py-10 text-center text-sm text-slate-500">Memuat...</p>}
      {!loading && empty && <p className="px-5 py-10 text-center text-sm text-slate-500">{empty}</p>}
    </div>
  );
}

export function Pagination({ page, limit, total, onChange }) {
  const totalHal = Math.max(1, Math.ceil(total / limit));
  const dari = total ? (page - 1) * limit + 1 : 0;
  const sampai = Math.min(page * limit, total);
  return (
    <div className="flex items-center justify-between border-t border-slate-100 px-5 py-3 text-sm text-slate-600">
      <span>
        {dari.toLocaleString("id-ID")}–{sampai.toLocaleString("id-ID")} dari {total.toLocaleString("id-ID")}
      </span>
      <div className="flex items-center gap-2">
        <Button variant="outline" size="sm" disabled={page <= 1} onClick={() => onChange(page - 1)}>Sebelumnya</Button>
        <span className="tabular-nums">{page} / {totalHal}</span>
        <Button variant="outline" size="sm" disabled={page >= totalHal} onClick={() => onChange(page + 1)}>Berikutnya</Button>
      </div>
    </div>
  );
}

export function Modal({ title, onClose, children, wide }) {
  return (
    <div className="fixed inset-0 z-50 flex items-start justify-center overflow-y-auto bg-slate-900/40 p-4">
      <div className={`mt-12 w-full rounded-xl bg-white shadow-xl ${wide ? "max-w-2xl" : "max-w-lg"}`}>
        <div className="flex items-center justify-between border-b border-slate-100 px-5 py-3.5">
          <h2 className="font-semibold">{title}</h2>
          <IconButton label="Tutup" onClick={onClose}><IconX className="size-4" /></IconButton>
        </div>
        <div className="p-5">{children}</div>
      </div>
    </div>
  );
}

const TONE = {
  green: "bg-emerald-50 text-emerald-700 ring-emerald-600/20",
  gray: "bg-slate-100 text-slate-600 ring-slate-500/20",
  red: "bg-red-50 text-red-700 ring-red-600/20",
  amber: "bg-amber-50 text-amber-700 ring-amber-600/20",
  blue: "bg-blue-50 text-blue-700 ring-blue-600/20",
};

export function Badge({ tone = "gray", children }) {
  return <span className={`inline-flex rounded-full px-2 py-0.5 text-xs font-medium ring-1 ring-inset ${TONE[tone]}`}>{children}</span>;
}

const TONE_STATUS = { AKTIF: "green", PENSIUN: "gray", MENINGGAL: "gray", DIBERHENTIKAN: "red", MUTASI_KELUAR: "amber" };

export function StatusBadge({ status }) {
  return <Badge tone={TONE_STATUS[status] ?? "gray"}>{STATUS_PERSONEL[status] ?? status}</Badge>;
}

export function Avatar({ nama, size = "md" }) {
  const kata = (nama ?? "").trim().split(/\s+/).filter(Boolean);
  const inisial = ((kata[0]?.[0] ?? "?") + (kata.length > 1 ? kata[kata.length - 1][0] : "")).toUpperCase();
  const ukuran = size === "lg" ? "size-16 text-xl" : size === "sm" ? "size-8 text-xs" : "size-10 text-sm";
  return <div className={`grid shrink-0 place-items-center rounded-full bg-dongker/10 font-semibold text-dongker ${ukuran}`}>{inisial}</div>;
}
