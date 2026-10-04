import { createContext, useCallback, useContext, useRef, useState } from "react";
import { Button } from "./ui";
import { IconAlert, IconCheck } from "./icons";

// Notifikasi kecil (toast) + dialog konfirmasi, dipakai dari halaman mana pun:
//   const toast = useToast();   toast("Data disimpan")  /  toast("Gagal", "error")
//   const confirm = useConfirm(); if (await confirm({ title, message })) { ... }
const Ctx = createContext(null);

export function FeedbackProvider({ children }) {
  const [toasts, setToasts] = useState([]);
  const [dialog, setDialog] = useState(null);
  const resolver = useRef(null);

  const toast = useCallback((message, type = "success") => {
    const id = Math.random();
    setToasts((t) => [...t, { id, message, type }]);
    setTimeout(() => setToasts((t) => t.filter((x) => x.id !== id)), 3500);
  }, []);

  const confirm = useCallback((opts) => {
    setDialog(opts);
    return new Promise((resolve) => (resolver.current = resolve));
  }, []);

  function tutup(hasil) {
    resolver.current?.(hasil);
    setDialog(null);
  }

  return (
    <Ctx.Provider value={{ toast, confirm }}>
      {children}

      <div className="fixed bottom-4 right-4 z-[60] space-y-2">
        {toasts.map((t) => (
          <div
            key={t.id}
            className={`flex items-center gap-2 rounded-lg px-4 py-3 text-sm text-white shadow-lg ${t.type === "error" ? "bg-red-600" : "bg-slate-900"}`}
          >
            {t.type === "error" ? <IconAlert className="size-4" /> : <IconCheck className="size-4 text-emerald-400" />}
            {t.message}
          </div>
        ))}
      </div>

      {dialog && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/40 p-4">
          <div className="w-full max-w-sm rounded-xl bg-white p-5 shadow-xl">
            <h2 className="font-semibold">{dialog.title}</h2>
            {dialog.message && <p className="mt-2 text-sm text-slate-600">{dialog.message}</p>}
            <div className="mt-5 flex justify-end gap-2">
              <Button variant="outline" onClick={() => tutup(false)}>Batal</Button>
              <Button variant={dialog.danger ? "dangerSolid" : "primary"} onClick={() => tutup(true)}>
                {dialog.confirmLabel ?? "Ya, lanjutkan"}
              </Button>
            </div>
          </div>
        </div>
      )}
    </Ctx.Provider>
  );
}

export const useToast = () => useContext(Ctx).toast;
export const useConfirm = () => useContext(Ctx).confirm;
