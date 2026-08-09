import { useEffect, useState, type ReactNode } from 'react'
import { createPortal } from 'react-dom'

export interface OverlayPortalProps {
  children: ReactNode
}

export function OverlayPortal({ children }: OverlayPortalProps) {
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    setMounted(true)
  }, [])

  if (!mounted) return null

  return createPortal(children, document.body)
}
