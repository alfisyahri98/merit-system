import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    port: 5173,
    // Request /api diteruskan ke backend Go, jadi tidak perlu CORS saat development.
    proxy: { "/api": "http://localhost:8080", "/docs": "http://localhost:8080" },
  },
});
