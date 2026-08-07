import type { PropsWithChildren } from 'react'
import { SessionProvider } from '#/features/auth/session'
import { ThemeProvider } from '#/features/layout/theme'
import { I18nProvider } from '#/providers/I18nProvider'
import { QueryProvider } from '#/providers/QueryProvider'

export function Providers({ children }: PropsWithChildren) {
  return (
    <I18nProvider>
      <ThemeProvider>
        <QueryProvider>
          <SessionProvider>{children}</SessionProvider>
        </QueryProvider>
      </ThemeProvider>
    </I18nProvider>
  )
}
