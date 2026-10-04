import { useState } from "react";
import { Link } from "react-router-dom";

// Donut satu warna bertingkat (gelap = pangkat tertinggi). Segmen & legenda bisa diklik.
const WARNA = ["#13294b", "#2d4a7a", "#4d6ea6", "#7f9bc9", "#b5c7e3"];

const SINGKAT = {
  "Perwira Tinggi": "Pati",
  "Perwira Menengah": "Pamen",
  "Perwira Pertama": "Pama",
};
const singkat = (k) => SINGKAT[k] ?? k.replace("Golongan", "Gol.");

export default function Donut({ judul, data, linkTo }) {
  const [aktif, setAktif] = useState(null);
  const total = data.reduce((s, d) => s + d.total, 0);
  const r = 42;
  const C = 2 * Math.PI * r;
  const celah = data.length > 1 ? 1.5 : 0;

  let offset = 0;
  const segmen = data.map((d, i) => {
    const panjang = total ? (d.total / total) * C : 0;
    const seg = { ...d, warna: WARNA[i % WARNA.length], dash: Math.max(0, panjang - celah), offset };
    offset += panjang;
    return seg;
  });
  const sorot = aktif !== null ? segmen[aktif] : null;

  return (
    <div className="flex items-center gap-5">
      <div className="relative size-32 shrink-0">
        <svg viewBox="0 0 100 100" className="size-full -rotate-90">
          <circle cx="50" cy="50" r={r} fill="none" stroke="#f1f5f9" strokeWidth="12" />
          {segmen.map((s, i) => (
            <circle
              key={s.kelompok}
              cx="50"
              cy="50"
              r={r}
              fill="none"
              stroke={s.warna}
              strokeWidth={aktif === i ? 15 : 12}
              strokeDasharray={`${s.dash} ${C - s.dash}`}
              strokeDashoffset={-s.offset}
              className="cursor-pointer transition-all"
              onMouseEnter={() => setAktif(i)}
              onMouseLeave={() => setAktif(null)}
            >
              <title>{`${s.kelompok}: ${s.total.toLocaleString("id-ID")}`}</title>
            </circle>
          ))}
        </svg>
        <div className="pointer-events-none absolute inset-0 grid place-items-center text-center">
          <div>
            <p className="text-lg font-semibold tabular-nums leading-none">{(sorot?.total ?? total).toLocaleString("id-ID")}</p>
            <p className="mt-1 text-[11px] text-slate-500">{sorot ? singkat(sorot.kelompok) : judul}</p>
          </div>
        </div>
      </div>

      <ul className="min-w-0 flex-1 space-y-1 text-sm">
        {segmen.map((s, i) => (
          <li key={s.kelompok}>
            <Link
              to={linkTo(s.kelompok)}
              onMouseEnter={() => setAktif(i)}
              onMouseLeave={() => setAktif(null)}
              className={`flex items-center gap-2 rounded-md px-2 py-1 ${aktif === i ? "bg-slate-100" : "hover:bg-slate-50"}`}
            >
              <span className="size-2.5 shrink-0 rounded-sm" style={{ background: s.warna }} />
              <span className="flex-1 truncate" title={s.kelompok}>{singkat(s.kelompok)}</span>
              <span className="tabular-nums text-slate-600">{s.total.toLocaleString("id-ID")}</span>
              <span className="w-9 text-right text-xs tabular-nums text-slate-400">{total ? Math.round((s.total / total) * 100) : 0}%</span>
            </Link>
          </li>
        ))}
      </ul>
    </div>
  );
}
