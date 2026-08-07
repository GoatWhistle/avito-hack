import { useState, type PropsWithChildren } from 'react'
import { I18nextProvider } from 'react-i18next'
import { initI18n } from '#/i18n'

export function I18nProvider({ children }: PropsWithChildren) {
  const [instance] = useState(() => initI18n())

  return <I18nextProvider i18n={instance}>{children}</I18nextProvider>
}
