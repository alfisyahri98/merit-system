import { useState } from "react";
import { Navigate, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { Alert, Button, Field, Input } from "../components/ui";
import logo from "../assets/logo-ssdm.png";

export default function LoginPage() {
  const { token, login } = useAuth();
  const navigate = useNavigate();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  if (token) return <Navigate to="/" replace />;

  async function submit(e) {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      await login(username.trim(), password);
      navigate("/", { replace: true });
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-slate-50 p-6">
      <div className="w-full max-w-sm">
        <div className="mb-8 text-center">
          <img src={logo} alt="Logo SSDM Polri" className="mx-auto h-20 w-auto" />
          <h1 className="mt-5 text-2xl font-semibold tracking-tight">Merit System Personel Polri</h1>
          <p className="mt-1 text-sm text-slate-500">Staf Sumber Daya Manusia Polri</p>
        </div>

        <form onSubmit={submit} className="space-y-5 rounded-xl border border-slate-200 bg-white p-6 shadow-sm">
          <Alert>{error}</Alert>
          <Field label="Username">
            <Input value={username} onChange={(e) => setUsername(e.target.value)} autoComplete="username" required autoFocus />
          </Field>
          <Field label="Password">
            <Input type="password" value={password} onChange={(e) => setPassword(e.target.value)} autoComplete="current-password" required />
          </Field>
          <Button type="submit" disabled={loading} className="w-full py-2.5">
            {loading ? "Memproses..." : "Masuk"}
          </Button>
        </form>

        <p className="mt-6 text-center text-xs text-slate-500">Lupa password? Hubungi Admin SSDM.</p>
      </div>
    </div>
  );
}
