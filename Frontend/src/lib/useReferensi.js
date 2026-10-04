import { useEffect, useState } from "react";
import { useAuth } from "../context/AuthContext";

// Data referensi (satker, pangkat, dll) jarang berubah → disimpan sekali per sesi.
const cache = new Map();

export function useReferensi(path) {
  const { request, token } = useAuth();
  const key = `${token}|${path}`;
  const [data, setData] = useState(() => cache.get(key) ?? []);

  useEffect(() => {
    if (cache.has(key)) {
      setData(cache.get(key));
      return;
    }
    let batal = false;
    request(path)
      .then((d) => {
        cache.set(key, d);
        if (!batal) setData(d);
      })
      .catch(() => {});
    return () => {
      batal = true;
    };
  }, [key, path, request]);

  return data;
}
