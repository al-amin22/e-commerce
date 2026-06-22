import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  appType: "spa",
  plugins: [react()],
  server: {
    host: "0.0.0.0",
    strictPort: true,
    port: 5173,
  },
  preview: {
    host: "0.0.0.0",
    strictPort: true,
    port: 4173,
  },
});
