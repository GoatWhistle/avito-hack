import { useEffect, type RefObject } from 'react'

export const useDismissable = (
  open: boolean,
  ref: RefObject<HTMLElement | null>,
  onDismiss: () => void,
) => {
  useEffect(() => {
    if (!open) return

    const onPointerDown = (event: MouseEvent) => {
      if (!ref.current?.contains(event.target as Node)) onDismiss()
    }
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') onDismiss()
    }

    document.addEventListener('mousedown', onPointerDown)
    document.addEventListener('keydown', onKeyDown)

    return () => {
      document.removeEventListener('mousedown', onPointerDown)
      document.removeEventListener('keydown', onKeyDown)
    }
  }, [open, ref, onDismiss])
}
