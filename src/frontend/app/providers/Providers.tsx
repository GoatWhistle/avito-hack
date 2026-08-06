import { QueryProvider } from './QueryProvider'
import type { PropsWithChildren } from 'react'

export function Providers({ children }: PropsWithChildren) {
  return <QueryProvider>{children}</QueryProvider>
}
