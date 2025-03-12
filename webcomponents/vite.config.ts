import { defineConfig, resolveConfig } from "vite";

export default defineConfig({
    build: {
        lib: {
            entry: "./src/main.ts",
            name: "webcomponents",
            fileName: "webcomponents"
        },
        sourcemap: true,
        rollupOptions: {
            output: {
                format: "es",
                entryFileNames: "bundle.js",
            }
        }
    }
})