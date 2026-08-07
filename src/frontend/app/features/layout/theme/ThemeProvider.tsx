import {
  useCallback,
  useEffect,
  useMemo,
  useState,
  type PropsWithChildren,
} from 'react'
import {
  isThemeMode,
  ThemeContext,
  themeStorageKey,
  type ThemeMode,
  type ThemeValue,
} from './theme-context'

const prefersDark = () =>
  typeof window !== 'undefined' &&
  typeof window.matchMedia === 'function' &&
  window.matchMedia('(prefers-color-scheme: dark)').matches

const resolveMode = (mode: ThemeMode): 'light' | 'dark' =>
  mode === 'system' ? (prefersDark() ? 'dark' : 'light') : mode

const applyTheme = (resolved: 'light' | 'dark') => {
  if (typeof document === 'undefined') return
  document.documentElement.classList.toggle('dark', resolved === 'dark')
  document.documentElement.style.colorScheme = resolved
}

export function ThemeProvider({ children }: PropsWithChildren) {
  const [mode, setModeState] = useState<ThemeMode>('system')
  const [resolved, setResolved] = useState<'light' | 'dark'>('light')

  useEffect(() => {
    const stored = window.localStorage.getItem(themeStorageKey)
    const initial = stored && isThemeMode(stored) ? stored : 'system'
    setModeState(initial)
    const next = resolveMode(initial)
    setResolved(next)
    applyTheme(next)
  }, [])

  useEffect(() => {
    if (mode !== 'system') return
    if (typeof window.matchMedia !== 'function') return

    const media = window.matchMedia('(prefers-color-scheme: dark)')
    const onChange = () => {
      const next = prefersDark() ? 'dark' : 'light'
      setResolved(next)
      applyTheme(next)
    }

    media.addEventListener('change', onChange)
    return () => media.removeEventListener('change', onChange)
  }, [mode])

  const setMode = useCallback((next: ThemeMode) => {
    setModeState(next)
    window.localStorage.setItem(themeStorageKey, next)
    const resolvedNext = resolveMode(next)
    setResolved(resolvedNext)
    applyTheme(resolvedNext)
  }, [])

  const toggle = useCallback(() => {
    setMode(resolved === 'dark' ? 'light' : 'dark')
  }, [resolved, setMode])

  const value = useMemo<ThemeValue>(
    () => ({ mode, resolved, setMode, toggle }),
    [mode, resolved, setMode, toggle],
  )

  return <ThemeContext value={value}>{children}</ThemeContext>
}
