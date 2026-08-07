import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { DesktopNav } from './DesktopNav'
import { LocaleToggle } from './LocaleToggle'
import { ThemeToggle } from './ThemeToggle'
import { navItems } from '#/features/layout/nav-items'
import { renderWithShell } from '#/features/layout/test-utils'
import { themeStorageKey } from '#/features/layout/theme'
import { changeLocale, i18n, initI18n } from '#/i18n'

beforeEach(() => {
  window.localStorage.clear()
  document.documentElement.classList.remove('dark')
})

afterEach(async () => {
  initI18n()
  await changeLocale('ru')
})

describe('DesktopNav', () => {
  it('renders a link per navigation item', () => {
    renderWithShell(<DesktopNav />)

    expect(screen.getByRole('navigation')).toBeInTheDocument()
    expect(screen.getAllByRole('link')).toHaveLength(navItems.length)
  })

  it('marks the icon of the active route', () => {
    const { container } = renderWithShell(<DesktopNav />, {
      initialEntries: ['/favorites'],
    })

    const activeLink = container.querySelector('a[aria-current="page"]')
    expect(activeLink).toHaveAttribute('href', '/favorites')
    expect(container.querySelectorAll('svg[aria-current="page"]')).toHaveLength(
      1,
    )
  })

  it('marks no icon when the route is outside the nav', () => {
    const { container } = renderWithShell(<DesktopNav />, {
      initialEntries: ['/profile'],
    })

    expect(container.querySelectorAll('svg[aria-current="page"]')).toHaveLength(
      0,
    )
  })
})

describe('ThemeToggle', () => {
  it('offers the dark theme while light is resolved', () => {
    renderWithShell(<ThemeToggle />)

    expect(screen.getByRole('button')).toHaveAccessibleName('Тема: Тёмная')
  })

  it('switches to dark and then offers light back', async () => {
    renderWithShell(<ThemeToggle />)
    const button = screen.getByRole('button')

    await userEvent.click(button)

    expect(button).toHaveAccessibleName('Тема: Светлая')
    expect(document.documentElement).toHaveClass('dark')
    expect(window.localStorage.getItem(themeStorageKey)).toBe('dark')

    await userEvent.click(button)

    expect(button).toHaveAccessibleName('Тема: Тёмная')
    expect(document.documentElement).not.toHaveClass('dark')
    expect(window.localStorage.getItem(themeStorageKey)).toBe('light')
  })

  it('starts dark when the stored preference says so', () => {
    window.localStorage.setItem(themeStorageKey, 'dark')

    renderWithShell(<ThemeToggle />)

    expect(screen.getByRole('button')).toHaveAccessibleName('Тема: Светлая')
  })
})

describe('LocaleToggle', () => {
  it('shows the current locale and offers the other one', () => {
    renderWithShell(<LocaleToggle />)

    const button = screen.getByRole('button')
    expect(button).toHaveTextContent('ru')
    expect(button).toHaveAccessibleName('Язык: English')
  })

  it('switches the language and the document lang attribute', async () => {
    renderWithShell(<LocaleToggle />)

    await userEvent.click(screen.getByRole('button'))

    expect(i18n.language).toBe('en')
    expect(document.documentElement.lang).toBe('en')
    expect(screen.getByRole('button')).toHaveTextContent('en')
  })

  it('falls back to ru for an unsupported active language', async () => {
    initI18n()
    await i18n.changeLanguage('fr')

    renderWithShell(<LocaleToggle />)

    const button = screen.getByRole('button')
    expect(button).toHaveTextContent('ru')
    expect(button).toHaveAccessibleName('Язык: English')
  })
})

describe('theme system preference', () => {
  it('follows the media query when no preference is stored', async () => {
    const listeners = new Set<() => void>()
    let dark = false
    const spy = vi.spyOn(window, 'matchMedia').mockImplementation(
      (query: string) =>
        ({
          get matches() {
            return dark
          },
          media: query,
          onchange: null,
          addEventListener: (_: string, fn: () => void) => listeners.add(fn),
          removeEventListener: (_: string, fn: () => void) =>
            listeners.delete(fn),
          addListener: () => {},
          removeListener: () => {},
          dispatchEvent: () => true,
        }) as unknown as MediaQueryList,
    )

    renderWithShell(<ThemeToggle />)

    expect(screen.getByRole('button')).toHaveAccessibleName('Тема: Тёмная')

    dark = true
    for (const fn of listeners) fn()

    expect(await screen.findByRole('button')).toHaveAccessibleName(
      'Тема: Светлая',
    )
    expect(document.documentElement).toHaveClass('dark')

    spy.mockRestore()
  })
})
