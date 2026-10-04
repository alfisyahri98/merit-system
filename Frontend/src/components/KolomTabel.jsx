import { useCallback, useEffect, useRef, useState } from "react";
import { createPortal } from "react-dom";
import { IconArrowDown, IconArrowUp, IconFilter, IconSort } from "./icons";

// Judul kolom tabel dengan tombol urut (klik judul) dan tombol filter (ikon corong).
// sortKey: kunci urut kolom ini. arahAwal: arah panah saat pertama diklik ("asc"/"desc").
// Klik judul berputar: urut → urut terbalik → kembali ke default.
export default function KolomTabel({ label, sortKey, sort, onSort, arahAwal = "asc", judulUrut, filterAktif, children }) {
  const [buka, setBuka] = useState(false);
  const tombol = useRef(null);
  const tutup = useCallback(() => setBuka(false), []);

  const aktif = sortKey && (sort === sortKey || sort === `-${sortKey}`);
  const terbalik = sort === `-${sortKey}`;
  const arah = aktif ? (terbalik ? (arahAwal === "asc" ? "desc" : "asc") : arahAwal) : null;

  function klikUrut() {
    if (!sortKey) return;
    onSort(sort === sortKey ? `-${sortKey}` : sort === `-${sortKey}` ? "" : sortKey);
  }

  const Panah = arah === "asc" ? IconArrowUp : arah === "desc" ? IconArrowDown : IconSort;

  return (
    <div className="flex items-center gap-1">
      {sortKey ? (
        <button
          type="button"
          onClick={klikUrut}
          title={judulUrut}
          className={`group inline-flex items-center gap-1 uppercase tracking-wide hover:text-slate-800 ${aktif ? "text-dongker" : ""}`}
        >
          {label}
          <Panah className={`size-3.5 ${aktif ? "" : "opacity-0 group-hover:opacity-60"}`} />
        </button>
      ) : (
        <span>{label}</span>
      )}
      {children && (
        <button
          ref={tombol}
          type="button"
          onClick={() => setBuka((v) => !v)}
          aria-label={`Filter ${label}`}
          className={`grid size-6 place-items-center rounded transition hover:bg-slate-200/70 ${
            filterAktif ? "bg-dongker/10 text-dongker" : "text-slate-400"
          }`}
        >
          <IconFilter className="size-3.5" />
        </button>
      )}
      {buka && (
        <Popover anchor={tombol.current} onClose={tutup}>
          {typeof children === "function" ? children(tutup) : children}
        </Popover>
      )}
    </div>
  );
}

// Panel melayang di bawah tombol. Dirender ke <body> supaya tidak terpotong tabel.
function Popover({ anchor, onClose, children }) {
  const panel = useRef(null);
  const [pos, setPos] = useState(null);

  useEffect(() => {
    if (!anchor) return;
    const r = anchor.getBoundingClientRect();
    const lebar = 288;
    setPos({ top: r.bottom + 6, left: Math.max(8, Math.min(r.left - 8, window.innerWidth - lebar - 8)) });

    const luar = (e) => panel.current && !panel.current.contains(e.target) && !anchor.contains(e.target) && onClose();
    const esc = (e) => e.key === "Escape" && onClose();
    const gulir = (e) => panel.current && !panel.current.contains(e.target) && onClose();
    document.addEventListener("mousedown", luar);
    document.addEventListener("keydown", esc);
    window.addEventListener("scroll", gulir, true);
    window.addEventListener("resize", onClose);
    return () => {
      document.removeEventListener("mousedown", luar);
      document.removeEventListener("keydown", esc);
      window.removeEventListener("scroll", gulir, true);
      window.removeEventListener("resize", onClose);
    };
  }, [anchor, onClose]);

  if (!pos) return null;
  return createPortal(
    <div
      ref={panel}
      style={{ top: pos.top, left: pos.left }}
      className="fixed z-50 w-72 rounded-lg border border-slate-200 bg-white p-3 text-sm font-normal normal-case tracking-normal text-slate-700 shadow-lg"
    >
      {children}
    </div>,
    document.body
  );
}

// Daftar pilihan tunggal (radio) untuk isi filter.
export function PilihanFilter({ value, onChange, options }) {
  return (
    <div className="space-y-0.5">
      {options.map(([v, label]) => (
        <button
          key={v || "semua"}
          type="button"
          onClick={() => onChange(v)}
          className={`flex w-full items-center justify-between rounded-md px-2 py-1.5 text-left hover:bg-slate-100 ${
            value === v ? "font-medium text-dongker" : ""
          }`}
        >
          {label}
          {value === v && <span className="size-1.5 rounded-full bg-dongker" />}
        </button>
      ))}
    </div>
  );
}
