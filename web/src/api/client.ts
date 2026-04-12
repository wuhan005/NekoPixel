import axios from 'axios'
import type { AxiosInstance, AxiosResponse } from 'axios'

export interface HttpResponse<T = unknown> {
  data: T
}

function resolveDefaultBaseURL(): string | undefined {
  if (typeof window !== 'undefined' && window.NEKO_CONFIG) {
    return window.NEKO_CONFIG.pixelBaseURL
  }
  if (import.meta.env?.VITE_API_BASE_URL) {
    return import.meta.env.VITE_API_BASE_URL
  }
  return undefined
}

export const httpClient: AxiosInstance = axios.create({
  baseURL: resolveDefaultBaseURL(),
})

httpClient.interceptors.response.use(
  async (response: AxiosResponse<HttpResponse>) => response.data,
  async (error) => Promise.reject(error)
)

export function setBaseURL(url: string) {
  httpClient.defaults.baseURL = url
}
