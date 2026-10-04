import { useState } from "react";
import { NavLink, Outlet, useLocation, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { NAMA_ROLE } from "../lib/format";
import { Avatar } from "./ui";
import { IconBook, IconHome, IconKey, IconLogout, IconMenu, IconSearch, IconUserCog, IconUsers, IconX } from "./icons";
import logo from "../assets/logo-ssdm.png";

function judul(path) {
  if (path === "/") return "Beranda";
  if (path === "/personel") return "Data Personel";
  if (path === "/personel/baru") return "Tambah Personel";
  if (/^\/personel\/\d+\/ubah$/.test(path)) return "Ubah Personel";
  if (path.startsWith("/personel/")) return "Profil Personel";
  if (path === "/pengguna") return "Pengguna";
  if (path === "/aplikasi") return "API Client";
  return "";
}

function Sidebar({ onNavigate }) {
  const { user, isAdmin, logout } = useAuth();

  const grup = [
    {
      label: "Menu",
      item: [
        { to: "/", label: "Beranda", icon: IconHome, end: true },
        { to: "/personel", label: "Data Personel", icon: IconUsers },
      ],
    },
  ];
  if (isAdmin) {
    grup.push({
      label: "Administrasi",
      item: [
        { to: "/pengguna", label: "Pengguna", icon: IconUserCog },
        { to: "/aplikasi", label: "API Client", icon: IconKey },
      ],
    });
  }

  return (
    <div className="flex h-full flex-col bg-dongker text-white">
      <div className="flex items-center gap-3 px-5 py-5">
        <img src={logo} alt="" className="h-10 w-auto" />
        <div className="leading-tight">
          <p className="font-semibold">Merit System</p>
          <p className="text-xs text-white/60">Personel Polri</p>
        </div>
      </div>
      <div className="mx-5 h-px bg-emas/60" />

      <nav className="flex-1 space-y-6 overflow-y-auto px-3 py-5">
        {grup.map((g) => (
          <div key={g.label}>
            <p className="mb-2 px-3 text-[11px] font-medium uppercase tracking-wider text-white/40">{g.label}</p>
            <ul className="space-y-1">
              {g.item.map(({ to, label, icon: Icon, end }) => (
                <li key={to}>
                  <NavLink
                    to={to}
                    end={end}
                    onClick={onNavigate}
                    className={({ isActive }) =>
                      `flex items-center gap-3 rounded-lg px-3 py-2 text-sm transition ${
                        isActive ? "bg-white/15 font-medium text-white" : "text-white/70 hover:bg-white/5 hover:text-white"
                      }`
                    }
                  >
                    <Icon className="size-5" />
                    {label}
                  </NavLink>
                </li>
              ))}
            </ul>
          </div>
        ))}
        {isAdmin && (
          <a href="/docs" target="_blank" rel="noreferrer" className="flex items-center gap-3 rounded-lg px-3 py-2 text-sm text-white/70 hover:bg-white/5 hover:text-white">
            <IconBook className="size-5" />
            Dokumentasi API
          </a>
        )}
      </nav>

      <div className="border-t border-white/10 p-3">
        <div className="flex items-center gap-3 rounded-lg px-2 py-2">
          <div className="grid size-9 place-items-center rounded-full bg-white/15 text-xs font-semibold">
            {(user?.username ?? "?").slice(0, 2).toUpperCase()}
          </div>
          <div className="min-w-0 flex-1 leading-tight">
            <p className="truncate text-sm font-medium">{user?.username}</p>
            <p className="text-xs text-white/60">{NAMA_ROLE[user?.role] ?? user?.role}</p>
          </div>
          <button onClick={logout} title="Keluar" className="rounded-md p-2 text-white/70 hover:bg-white/10 hover:text-white">
            <IconLogout className="size-5" />
          </button>
        </div>
      </div>
    </div>
  );
}

export default function AppLayout() {
  const [menuHp, setMenuHp] = useState(false);
  const [cari, setCari] = useState("");
  const { pathname } = useLocation();
  const navigate = useNavigate();

  function submitCari(e) {
    e.preventDefault();
    const q = cari.trim();
    navigate(q ? `/personel?q=${encodeURIComponent(q)}` : "/personel");
    setCari("");
  }

  return (
    <div className="min-h-screen">
      {/* Sidebar desktop */}
      <aside className="fixed inset-y-0 left-0 z-40 hidden w-64 lg:block">
        <Sidebar />
      </aside>

      {/* Sidebar HP */}
      {menuHp && (
        <div className="fixed inset-0 z-50 lg:hidden">
          <div className="absolute inset-0 bg-slate-900/50" onClick={() => setMenuHp(false)} />
          <aside className="absolute inset-y-0 left-0 w-64">
            <Sidebar onNavigate={() => setMenuHp(false)} />
            <button onClick={() => setMenuHp(false)} className="absolute right-3 top-5 rounded-md p-1 text-white/70" aria-label="Tutup menu">
              <IconX className="size-5" />
            </button>
          </aside>
        </div>
      )}

      <div className="lg:pl-64">
        <header className="sticky top-0 z-30 flex items-center gap-3 border-b border-slate-200 bg-white/90 px-4 py-3 backdrop-blur sm:px-6">
          <button onClick={() => setMenuHp(true)} className="rounded-md p-1.5 text-slate-600 hover:bg-slate-100 lg:hidden" aria-label="Buka menu">
            <IconMenu />
          </button>
          <h1 className="text-base font-semibold">{judul(pathname)}</h1>
          <form onSubmit={submitCari} className="ml-auto w-full max-w-xs">
            <div className="relative">
              <IconSearch className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-slate-400" />
              <input
                value={cari}
                onChange={(e) => setCari(e.target.value)}
                placeholder="Cari personel (nama / NRP)"
                className="w-full rounded-lg border border-slate-200 bg-slate-50 py-2 pl-9 pr-3 text-sm focus:border-dongker focus:bg-white focus:outline-none"
              />
            </div>
          </form>
        </header>

        <main className="mx-auto max-w-7xl px-4 py-6 sm:px-6">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
