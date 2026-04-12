import type { App as VueApp, Plugin } from 'vue'
import NekoPixelApp from './App.vue'
import './style.css'
import { setBaseURL } from './api/client'

export type NekoPixelOptions = {
  baseURL?: string
}

export const NekoPixelPlugin: Plugin = {
  install(app: VueApp, options: NekoPixelOptions = {}) {
    if (options.baseURL) {
      setBaseURL(options.baseURL)
    }
    app.component('NekoPixelApp', NekoPixelApp)
  },
}

export { NekoPixelApp }
export type { pixelStatus, canvasColor, pixel, setPixelData } from './api/pixel'
export default NekoPixelPlugin
