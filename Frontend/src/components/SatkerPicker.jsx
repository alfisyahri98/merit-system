import { useEffect, useMemo, useRef, useState } from "react";
import { IconChevronDown } from "./icons";
import { letakSatker, rapikan } from "../lib/format";

// Singkatan yang biasa diketik pengguna
const ALIAS = {
  jakpus: "jakarta pusat",
  jakut: "jakarta utara",
  jakbar: "jakarta barat",
  jaksel: "jakarta selatan",
  jaktim: "jakarta timur",
  pmj: "polda metro jaya",
  mabes: "mabes polri",
};

// Semua kata yang diketik harus ada di nama atau kode satker (urutan bebas).
function cocok(satker, query) {
  const teks = `${satker.nama} ${satker.kode}`.toLowerCase();
  return query
    .split(/\s+/)
    .flatMap((k) => (ALIAS[k] ?? k).split(" "))
    .every((k) => teks.includes(k));
}

// Pilih satker dengan mengetik (nama, kode, atau singkatan). Menampilkan maks 50 hasil.
export default function SatkerPicker({ options, value, onChange, placeholder = "Ketik nama satker, mis. jaksel / reskrim", emptyLabel }) {
  const [open, setOpen] = useState(false);
  const [q, setQ] = useState("");
  const [aktif, setAktif] = useState(0);
  const box = useRef(null);

  const selected = options.find((o) => o.id === Number(value));

  const hasil = useMemo(() => {
    const s = q.trim().toLowerCase();
    const list = s ? options.filter((o) => cocok(o, s)) : options;
    const top = list.slice(0, 50);
    return emptyLabel ? [{ id: "", nama: emptyLabel, kode: "" }, ...top] : top;
  }, [q, options, emptyLabel]);

  useEffect(() => {
    const klikLuar = (e) => box.current && !box.current.contains(e.target) && setOpen(false);
    document.addEventListener("mousedown", klikLuar);
    return () => document.removeEventListener("mousedown", klikLuar);
  }, []);

  function pilih(o) {
    onChange(o.id);
    setOpen(false);
    setQ("");
  }

  function onKey(e) {
    if (e.key === "ArrowDown") {
      e.preventDefault();
      setOpen(true);
      setAktif((i) => Math.min(i + 1, hasil.length - 1));
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      setAktif((i) => Math.max(i - 1, 0));
    } else if (e.key === "Enter" && open && hasil[aktif]) {
      e.preventDefault();
      pilih(hasil[aktif]);
    } else if (e.key === "Escape") {
      setOpen(false);
    }
  }

  const tampil = open ? q : selected ? rapikan(selected.nama) : emptyLabel && value === "" ? emptyLabel : "";

  return (
    <div ref={box} className="relative">
      <input
        className="w-full rounded-lg border border-slate-300 bg-white py-2 pl-3 pr-9 text-sm focus:border-dongker focus:outline-none focus:ring-2 focus:ring-dongker/15"
        value={tampil}
        placeholder={placeholder}
        onFocus={() => {
          setOpen(true);
          setQ("");
          setAktif(0);
        }}
        onChange={(e) => {
          setQ(e.target.value);
          setAktif(0);
          setOpen(true);
        }}
        onKeyDown={onKey}
      />
      <IconChevronDown className="pointer-events-none absolute right-3 top-1/2 size-4 -translate-y-1/2 text-slate-400" />
      {selected && !open && <p className="mt-1 text-xs text-slate-500">{letakSatker(selected, options)}</p>}

      {open && (
        <ul className="absolute z-30 mt-1 max-h-72 w-full overflow-auto rounded-lg border border-slate-200 bg-white py-1 text-sm shadow-lg">
          {hasil.length === 0 && <li className="px-3 py-2 text-slate-500">Satker tidak ditemukan.</li>}
          {hasil.map((o, i) => (
            <li
              key={o.id || "kosong"}
              onMouseDown={(e) => {
                e.preventDefault();
                pilih(o);
              }}
              onMouseEnter={() => setAktif(i)}
              className={`cursor-pointer px-3 py-2 ${i === aktif ? "bg-slate-100" : ""}`}
            >
              <p className={o.id === Number(value) ? "font-medium text-dongker" : ""}>{o.id ? rapikan(o.nama) : o.nama}</p>
              {o.id && <p className="text-xs text-slate-500">{letakSatker(o, options)}</p>}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
