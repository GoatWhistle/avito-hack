import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'

import {
  defaultNamespace,
  fallbackLocale,
  isLocale,
  namespaces,
  ogLocales,
  resources,
  supportedLocales,
  type Locale,
} from './resources'

const STORAGE_KEY = 'avito-hack.locale'

export const readStoredLocale = (): Locale => {
  if (typeof window === 'undefined') return fallbackLocale
  const stored = window.localStorage.getItem(STORAGE_KEY)
  return stored && isLocale(stored) ? stored : fallbackLocale
}

const applyDocumentLocale = (locale: Locale) => {
  if (typeof document === 'undefined') return

  document.documentElement.lang = locale

  const { app } = resources[locale].common
  document.title = app.metaTitle

  const ogLocale = document.querySelector('meta[property="og:locale"]')
  if (ogLocale) ogLocale.setAttribute('content', ogLocales[locale])
}

export const persistLocale = (locale: Locale) => {
  if (typeof window === 'undefined') return
  window.localStorage.setItem(STORAGE_KEY, locale)
}

export const changeLocale = async (locale: Locale) => {
  persistLocale(locale)
  await i18n.changeLanguage(locale)
  applyDocumentLocale(locale)
}

export const initI18n = () => {
  if (i18n.isInitialized) return i18n

  const locale = readStoredLocale()

  void i18n.use(initReactI18next).init({
    resources,
    lng: locale,
    fallbackLng: fallbackLocale,
    supportedLngs: supportedLocales,
    ns: namespaces,
    defaultNS: defaultNamespace,
    interpolation: { escapeValue: false },
    returnNull: false,
  })

  applyDocumentLocale(locale)

  return i18n
}

export { i18n }
