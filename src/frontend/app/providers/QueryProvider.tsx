import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { useState, type PropsWithChildren } from 'react'
import { isApiError } from '#/api'

const createQueryClient = () =>
  new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 30_000,
        retry: (failureCount, error) => {
          if (
            isApiError(error) &&
            error.status !== null &&
            error.status < 500
          ) {
            return false
          }
          return failureCount < 2
        },
      },
    },
  })

export function QueryProvider({ children }: PropsWithChildren) {
  const [queryClient] = useState(createQueryClient)

  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  )
}
