import { useCallback, useEffect, useRef, useState } from 'react'

const focusableSelector = [
  'a[href]',
  'button:not([disabled])',
  'input:not([disabled])',
  'select:not([disabled])',
  'textarea:not([disabled])',
  '[tabindex]:not([tabindex="-1"])',
].join(', ')

const focusableIn = (container: HTMLElement): HTMLElement[] =>
  Array.from(container.querySelectorAll<HTMLElement>(focusableSelector)).filter(
    (node) =>
      !node.hasAttribute('disabled') &&
      node.getAttribute('aria-hidden') !== 'true',
  )

export const useFocusTrap = (onEscape: () => void) => {
  const escapeRef = useRef(onEscape)
  escapeRef.current = onEscape
  const [container, setContainer] = useState<HTMLElement | null>(null)

  const containerRef = useCallback((node: HTMLElement | null) => {
    setContainer(node)
  }, [])

  useEffect(() => {
    if (!container) return

    const previous = document.activeElement as HTMLElement | null
    const initial = focusableIn(container)[0]
    if (initial) initial.focus()
    else container.focus()

    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        event.preventDefault()
        escapeRef.current()

        return
      }

      if (event.key !== 'Tab') return

      const nodes = focusableIn(container)
      if (nodes.length === 0) {
        event.preventDefault()
        container.focus()

        return
      }

      const first = nodes[0]
      const last = nodes[nodes.length - 1]
      const active = document.activeElement

      if (event.shiftKey && (active === first || !container.contains(active))) {
        event.preventDefault()
        last.focus()

        return
      }

      if (!event.shiftKey && active === last) {
        event.preventDefault()
        first.focus()
      }
    }

    document.addEventListener('keydown', onKeyDown)

    return () => {
      document.removeEventListener('keydown', onKeyDown)
      previous?.focus?.()
    }
  }, [container])

  return containerRef
}
