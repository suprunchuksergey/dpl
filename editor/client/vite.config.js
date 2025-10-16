import { defineConfig } from "vite";
import tailwindcss from "@tailwindcss/vite";
import monacoEditorPlugin from "vite-plugin-monaco-editor";

export default defineConfig({
  plugins: [tailwindcss(), monacoEditorPlugin],
  base: "./",
  build: {
    outDir: "../../docs",
    emptyOutDir: true,
  },
});
