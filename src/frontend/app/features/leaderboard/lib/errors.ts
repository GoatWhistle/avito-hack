import axios from 'axios'

export const leaderboardErrorKey = (error: unknown): string => {
  if (axios.isAxiosError(error)) {
    if (!error.response) return 'errors.network'
    if (error.response.status === 401) return 'errors.unauthorized'
  }

  return 'errors.loadFailed'
}
