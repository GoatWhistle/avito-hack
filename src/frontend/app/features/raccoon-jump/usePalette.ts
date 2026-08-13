import { useEffect, useRef, useState } from 'react'
import { readGamePalette, type GamePalette } from './palette'

export const usePalette = (): { current: GamePalette | null } => {
  const paletteRef = useRef<GamePalette | null>(null)
  const [, setVersion] = useState(0)

  useEffect(() => {
    const root = document.documentElement

    const refresh = () => {
      paletteRef.current = readGamePalette(root)
      setVersion((value) => value + 1)
    }

    refresh()

    const observer = new MutationObserver(refresh)
    observer.observe(root, { attributes: true, attributeFilter: ['class'] })

    return () => observer.disconnect()
  }, [])

  return paletteRef
}
