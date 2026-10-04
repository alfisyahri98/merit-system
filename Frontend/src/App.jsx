import { Navigate, Route, Routes } from "react-router-dom";
import ProtectedRoute from "./components/ProtectedRoute";
import AppLayout from "./components/AppLayout";
import LoginPage from "./pages/LoginPage";
import BerandaPage from "./pages/BerandaPage";
import PersonelPage from "./pages/PersonelPage";
import PersonelDetailPage from "./pages/PersonelDetailPage";
import PersonelFormPage from "./pages/PersonelFormPage";
import PenggunaPage from "./pages/PenggunaPage";
import AplikasiPage from "./pages/AplikasiPage";

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route
        element={
          <ProtectedRoute>
            <AppLayout />
          </ProtectedRoute>
        }
      >
        <Route index element={<BerandaPage />} />
        <Route path="personel" element={<PersonelPage />} />
        <Route path="personel/baru" element={<PersonelFormPage />} />
        <Route path="personel/:id" element={<PersonelDetailPage />} />
        <Route path="personel/:id/ubah" element={<PersonelFormPage />} />
        <Route path="pengguna" element={<ProtectedRoute adminOnly><PenggunaPage /></ProtectedRoute>} />
        <Route path="aplikasi" element={<ProtectedRoute adminOnly><AplikasiPage /></ProtectedRoute>} />
      </Route>
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
