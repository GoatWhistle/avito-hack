import { useEffect, type RefObject } from 'react'
import type { Game } from './game'

const MOVE_LEFT_KEYS = ['ArrowLeft', 'a', 'A', 'ф', 'Ф']
const MOVE_RIGHT_KEYS = ['ArrowRight', 'd', 'D', 'в', 'В']

const INTERACTIVE_SELECTOR =
  'input, textarea, select, button, a[href], [contenteditable=""], [contenteditable="true"], [role="dialog"], [role="textbox"]'

export const useGameKeyboard = (
  gameRef: RefObject<Game | null>,
  surfaceRef: RefObject<HTMLElement | null>,
): void => {
  useEffect(() => {
    const surface = surfaceRef.current

    const shouldIgnore = (node: EventTarget | null) => {
      if (!(node instanceof Element)) return false
      if (surface?.contains(node)) return false
      return Boolean(node.closest(INTERACTIVE_SELECTOR))
    }

    const onKeyDown = (e: KeyboardEvent) => {
      const game = gameRef.current
      if (!game) return
      if (e.altKey || e.ctrlKey || e.metaKey) return
      if (shouldIgnore(e.target) || shouldIgnore(document.activeElement)) return

      if (MOVE_LEFT_KEYS.includes(e.key)) {
        game.input.left = true
        e.preventDefault()
      }

      if (MOVE_RIGHT_KEYS.includes(e.key)) {
        game.input.right = true
        e.preventDefault()
      }

      if (e.key === ' ') {
        e.preventDefault()
        if (game.status !== 'playing') game.reset()
      }
    }

    const onKeyUp = (e: KeyboardEvent) => {
      const game = gameRef.current
      if (!game) return
      if (MOVE_LEFT_KEYS.includes(e.key)) game.input.left = false
      if (MOVE_RIGHT_KEYS.includes(e.key)) game.input.right = false
    }

    const onBlur = () => {
      const game = gameRef.current
      if (!game) return
      game.input.left = false
      game.input.right = false
    }

    window.addEventListener('keydown', onKeyDown)
    window.addEventListener('keyup', onKeyUp)
    window.addEventListener('blur', onBlur)

    return () => {
      window.removeEventListener('keydown', onKeyDown)
      window.removeEventListener('keyup', onKeyUp)
      window.removeEventListener('blur', onBlur)
    }
  }, [gameRef, surfaceRef])
}
