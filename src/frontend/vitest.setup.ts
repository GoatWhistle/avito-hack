import '@testing-library/jest-dom/vitest'
import { afterEach, vi } from 'vitest'
import { cleanup } from '@testing-library/react'

if (typeof window !== 'undefined') {
  Object.defineProperty(window.navigator, 'language', {
    configurable: true,
    get: () => 'ru',
  })
  Object.defineProperty(window.navigator, 'languages', {
    configurable: true,
    get: () => ['ru'],
  })
}

if (typeof window !== 'undefined' && typeof window.matchMedia !== 'function') {
  window.matchMedia = vi.fn().mockImplementation((query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    addListener: vi.fn(),
    removeListener: vi.fn(),
    dispatchEvent: vi.fn(),
  }))
}

afterEach(() => {
  cleanup()
})
