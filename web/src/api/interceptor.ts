import axios from 'axios';
import type { AxiosResponse } from 'axios';

export interface HttpResponse<T = unknown> {
  data: T;
}

if (import.meta.env.VITE_API_BASE_URL) {
  axios.defaults.baseURL = import.meta.env.VITE_API_BASE_URL;
}
axios.interceptors.response.use(
    async (response: AxiosResponse<HttpResponse>) => {
    return response.data
  },
  async (error) => {
    return Promise.reject(error);
  }
);
