import { localTokenStorage } from '#/shared/storage/local-token.storage'
import axios, { type CreateAxiosDefaults } from 'axios'

const config: CreateAxiosDefaults = {
  baseURL: import.meta.env.VITE_API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true,
}

export const httpClient = axios.create(config)

httpClient.interceptors.request.use(config => {
  const token = localTokenStorage.get()

  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }

  return config
})

httpClient.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      localTokenStorage.remove()

      if (window.location.pathname !== '/') {
        window.location.href = '/'
      }
    }
    return Promise.reject(error)
  },
)
