import { defineConfig, loadEnv } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, ".", "");
  const apiTarget = env.API_URL || env.VITE_API_PROXY_TARGET || "http://localhost:8080";

  return {
    plugins: [react()],
    server: {
      host: "0.0.0.0",
      port: 5173,
      proxy: {
        "/api": {
          target: apiTarget,
          changeOrigin: true,
        },
        "/health": {
          target: apiTarget,
          changeOrigin: true,
        },
        "/docs": {
          target: apiTarget,
          changeOrigin: true,
        },
        "/openapi": {
          target: apiTarget,
          changeOrigin: true,
        },
      },
    },
  };
});
