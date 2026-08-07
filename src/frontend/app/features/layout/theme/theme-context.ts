import { createContext, use } from 'react'

export const themeModes = ['light', 'dark', 'system'] as const

export type ThemeMode = (typeof themeModes)[number]

export const themeStorageKey = 'avito-hack.theme'

export interface ThemeValue {
  mode: ThemeMode
  resolved: 'light' | 'dark'
  setMode: (mode: ThemeMode) => void
  toggle: () => void
}

export const ThemeContext = createContext<ThemeValue | null>(null)

export const useTheme = (): ThemeValue => {
  const value = use(ThemeContext)

  if (!value) {
    throw new Error('useTheme must be used inside ThemeProvider')
  }

  return value
}

export const isThemeMode = (value: string): value is ThemeMode =>
  (themeModes as readonly string[]).includes(value)
