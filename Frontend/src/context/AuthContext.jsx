import { createContext, useCallback, useContext, useState } from "react";
import { ApiError, BASE_URL, api } from "../lib/api";

const AuthContext = createContext(null);
const KEY = "merit_sesi";

function tokenMasihBerlaku(token) {
  try {
    return JSON.parse(atob(token.split(".")[1])).exp * 1000 > Date.now();
  } catch {
    return false;
  }
}

function bacaSesi() {
  try {
    const sesi = JSON.parse(localStorage.getItem(KEY));
    if (sesi?.token && tokenMasihBerlaku(sesi.token)) return sesi;
  } catch {
    /* sesi rusak */
  }
  localStorage.removeItem(KEY);
  return null;
}

export function AuthProvider({ children }) {
  const [sesi, setSesi] = useState(bacaSesi);

  async function login(username, password) {
    const data = await api("/login", { method: "POST", body: { username, password } });
    const baru = { token: data.token, user: data.user };
    localStorage.setItem(KEY, JSON.stringify(baru));
    setSesi(baru);
  }

  const logout = useCallback(() => {
    localStorage.removeItem(KEY);
    setSesi(null);
  }, []);

  // Request dengan token. 401 = token habis / akun dinonaktifkan → keluar.
  const request = useCallback(
    async (path, opts = {}) => {
      try {
        return await api(path, { ...opts, token: sesi?.token });
      } catch (err) {
        if (err.status === 401) logout();
        throw err;
      }
    },
    [sesi?.token, logout]
  );

  // Unduh file (mis. Excel) dari endpoint yang butuh token
  const download = useCallback(
    async (path, namaCadangan) => {
      const res = await fetch(BASE_URL + path, { headers: { Authorization: `Bearer ${sesi?.token}` } });
      if (!res.ok) {
        if (res.status === 401) logout();
        const json = await res.json().catch(() => ({}));
        throw new ApiError(json.error ?? "Gagal mengunduh file", res.status);
      }
      const nama = /filename="([^"]+)"/.exec(res.headers.get("Content-Disposition") ?? "")?.[1] ?? namaCadangan;
      const url = URL.createObjectURL(await res.blob());
      const a = document.createElement("a");
      a.href = url;
      a.download = nama;
      a.click();
      URL.revokeObjectURL(url);
    },
    [sesi?.token, logout]
  );

  const user = sesi?.user;
  const isAdmin = user?.role === "ADMIN_SSDM";

  return (
    <AuthContext.Provider value={{ token: sesi?.token, user, isAdmin, login, logout, request, download }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}
