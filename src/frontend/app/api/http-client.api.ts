import axios, { type CreateAxiosDefaults } from 'axios'

const config: CreateAxiosDefaults = {
  baseURL: import.meta.env.VITE_API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true,
}

export const httpClient = axios.create(config)
