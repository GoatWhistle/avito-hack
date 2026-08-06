import { QueryProvider } from '#/providers/QueryProvider'
import type { PropsWithChildren } from 'react'

export function Providers({ children }: PropsWithChildren) {
  return <QueryProvider>{children}</QueryProvider>
}
