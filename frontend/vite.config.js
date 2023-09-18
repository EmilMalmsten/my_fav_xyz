import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import path from "path";
import dotenv from "dotenv";

export default defineConfig(({ mode }) => {
    const envFile =
        mode === "development" ? ".env.development" : ".env.production";
    dotenv.config({ path: envFile });
    return {
        plugins: [react()],
        build: {
            outDir: "build",
            assetsDir: "assets",
            emptyOutDir: true,
        },
        resolve: {
            alias: {
                "@": path.resolve(__dirname, "./src"),
            },
        },
        server: {
            proxy: {
                "/api": {
                    target: process.env.VITE_API_URL,
                    changeOrigin: true,
                },
            },
        },
    };
});
