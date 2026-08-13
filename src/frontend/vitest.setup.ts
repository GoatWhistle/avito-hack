import '@testing-library/jest-dom/vitest'
import { afterEach, vi } from 'vitest'
import { cleanup } from '@testing-library/react'
import { createElement } from 'react'

vi.mock('@lottiefiles/dotlottie-react', () => ({
  DotLottieReact: ({
    src,
    className,
    dotLottieRefCallback,
  }: {
    src?: string
    className?: string
    dotLottieRefCallback?: (instance: unknown) => void
  }) => {
    dotLottieRefCallback?.({
      addEventListener: (event: string, listener: () => void) => {
        if (event === 'load') listener()
      },
      removeEventListener: () => {},
    })

    return createElement('canvas', {
      'data-testid': 'dotlottie-canvas',
      'data-src': src,
      className,
    })
  },
  setWasmUrl: vi.fn(),
}))

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

if (
  typeof globalThis !== 'undefined' &&
  typeof globalThis.ResizeObserver !== 'function'
) {
  globalThis.ResizeObserver = class {
    observe() {}
    unobserve() {}
    disconnect() {}
  }
}

afterEach(() => {
  cleanup()
})
