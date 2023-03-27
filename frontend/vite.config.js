import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

export default defineConfig({
  base: "./",
  plugins: [vue()],
  build: {
    chunkSizeWarningLimit: 1500
  }
})
