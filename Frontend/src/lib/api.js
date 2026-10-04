export const BASE_URL = import.meta.env.VITE_API_URL ?? "/api";

export class ApiError extends Error {
  constructor(message, status, details) {
    super(message);
    this.status = status;
    this.details = details; // object field → pesan (validasi), atau null
  }

  // Error per field untuk ditampilkan di bawah input form
  get fields() {
    return this.details && typeof this.details === "object" ? this.details : {};
  }
}

// Panggil backend. Mengembalikan isi `data` dari response { success, data }.
export async function api(path, { method = "GET", body, token } = {}) {
  const headers = { "Content-Type": "application/json" };
  if (token) headers.Authorization = `Bearer ${token}`;

  let res;
  try {
    res = await fetch(BASE_URL + path, {
      method,
      headers,
      body: body ? JSON.stringify(body) : undefined,
    });
  } catch {
    throw new ApiError("Server tidak bisa dihubungi. Pastikan backend sudah berjalan.", 0);
  }

  const json = await res.json().catch(() => ({}));
  if (!res.ok) {
    throw new ApiError(json.error ?? `Permintaan gagal (${res.status})`, res.status, json.details);
  }
  return json.data;
}
