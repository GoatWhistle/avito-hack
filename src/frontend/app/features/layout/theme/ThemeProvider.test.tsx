import { render, renderHook, screen, act } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { ThemeProvider } from './ThemeProvider'
import { isThemeMode, themeStorageKey, useTheme } from './theme-context'

function Probe() {
  const { mode, resolved, setMode } = useTheme()

  return (
    <div>
      <span data-testid="mode">{mode}</span>
      <span data-testid="resolved">{resolved}</span>
      <button type="button" onClick={() => setMode('system')}>
        system
      </button>
      <button type="button" onClick={() => setMode('dark')}>
        dark
      </button>
    </div>
  )
}

const renderProbe = () =>
  render(
    <ThemeProvider>
      <Probe />
    </ThemeProvider>,
  )

beforeEach(() => {
  window.localStorage.clear()
  document.documentElement.classList.remove('dark')
  document.documentElement.style.colorScheme = ''
})

describe('isThemeMode', () => {
  it('accepts the known modes and rejects anything else', () => {
    expect(isThemeMode('light')).toBe(true)
    expect(isThemeMode('dark')).toBe(true)
    expect(isThemeMode('system')).toBe(true)
    expect(isThemeMode('sepia')).toBe(false)
  })
})

describe('useTheme', () => {
  it('throws outside of a provider', () => {
    expect(() => renderHook(() => useTheme())).toThrow(
      /useTheme must be used inside ThemeProvider/,
    )
  })
})

describe('ThemeProvider', () => {
  it('defaults to the system mode resolved as light', () => {
    renderProbe()

    expect(screen.getByTestId('mode')).toHaveTextContent('system')
    expect(screen.getByTestId('resolved')).toHaveTextContent('light')
    expect(document.documentElement.style.colorScheme).toBe('light')
  })

  it('restores a stored mode', () => {
    window.localStorage.setItem(themeStorageKey, 'dark')

    renderProbe()

    expect(screen.getByTestId('mode')).toHaveTextContent('dark')
    expect(screen.getByTestId('resolved')).toHaveTextContent('dark')
    expect(document.documentElement).toHaveClass('dark')
  })

  it('ignores a corrupted stored mode', () => {
    window.localStorage.setItem(themeStorageKey, 'neon')

    renderProbe()

    expect(screen.getByTestId('mode')).toHaveTextContent('system')
  })

  it('persists an explicit mode change', async () => {
    renderProbe()

    await userEvent.click(screen.getByRole('button', { name: 'dark' }))

    expect(screen.getByTestId('mode')).toHaveTextContent('dark')
    expect(screen.getByTestId('resolved')).toHaveTextContent('dark')
    expect(window.localStorage.getItem(themeStorageKey)).toBe('dark')
    expect(document.documentElement.style.colorScheme).toBe('dark')
  })

  it('resolves system mode against the media query', async () => {
    const spy = vi.spyOn(window, 'matchMedia').mockReturnValue({
      matches: true,
      addEventListener: () => {},
      removeEventListener: () => {},
    } as unknown as MediaQueryList)

    renderProbe()
    await userEvent.click(screen.getByRole('button', { name: 'system' }))

    expect(screen.getByTestId('resolved')).toHaveTextContent('dark')
    spy.mockRestore()
  })

  it('does not subscribe to media changes for an explicit mode', async () => {
    const removeEventListener = vi.fn()
    const spy = vi.spyOn(window, 'matchMedia').mockReturnValue({
      matches: false,
      addEventListener: () => {},
      removeEventListener,
    } as unknown as MediaQueryList)

    renderProbe()
    await userEvent.click(screen.getByRole('button', { name: 'dark' }))

    expect(screen.getByTestId('mode')).toHaveTextContent('dark')
    spy.mockRestore()
  })

  it('survives an environment without matchMedia', () => {
    const original = window.matchMedia
    Reflect.deleteProperty(window, 'matchMedia')

    renderProbe()

    expect(screen.getByTestId('resolved')).toHaveTextContent('light')
    window.matchMedia = original
  })

  it('unsubscribes from the media query on unmount', () => {
    const removeEventListener = vi.fn()
    const spy = vi.spyOn(window, 'matchMedia').mockReturnValue({
      matches: false,
      addEventListener: () => {},
      removeEventListener,
    } as unknown as MediaQueryList)

    const view = renderProbe()
    act(() => view.unmount())

    expect(removeEventListener).toHaveBeenCalled()
    spy.mockRestore()
  })
})
