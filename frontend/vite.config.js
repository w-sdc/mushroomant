import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { defineConfig } from 'vite'
import  closure  from '@ampproject/rollup-plugin-closure-compiler'


const __dirname = dirname(fileURLToPath(import.meta.url))

export default defineConfig({
  build: {
    rollupOptions: {
      plugins: [
        closure()
      ],
      input: {
        main: resolve(__dirname, 'index.html'),
        ctrplane: resolve(__dirname, 'subpages/ctrplane.html'),
      },
    },
  },
})
