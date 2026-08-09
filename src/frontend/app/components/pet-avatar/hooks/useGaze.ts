import { useEffect, useRef, useState } from 'react'
import type { RefObject } from 'react'

export const GAZE_RADIUS_PX = 6
export const GAZE_TILT_DEG = 1.4
export const GAZE_FALLOFF_PX = 420

export interface GazeOffset {
  x: number
  y: number
  tilt: number
}

const ZERO: GazeOffset = { x: 0, y: 0, tilt: 0 }

const isTouchOnly = (): boolean => {
  if (typeof window.matchMedia !== 'function') return false

  return (
    window.matchMedia('(hover: none)').matches ||
    window.matchMedia('(pointer: coarse)').matches
  )
}

const clampUnit = (value: number): number => Math.max(-1, Math.min(1, value))

export const gazeFor = (
  dx: number,
  dy: number,
  falloff: number = GAZE_FALLOFF_PX,
): GazeOffset => {
  const distance = Math.hypot(dx, dy)
  const damping = distance === 0 ? 0 : Math.min(1, falloff / distance) ** 0.5
  const nx = clampUnit(dx / falloff) * damping
  const ny = clampUnit(dy / falloff) * damping

  return {
    x: nx * GAZE_RADIUS_PX,
    y: ny * GAZE_RADIUS_PX,
    tilt: nx * GAZE_TILT_DEG,
  }
}

export const useGaze = (
  ref: RefObject<HTMLElement | null>,
  enabled: boolean,
): GazeOffset => {
  const [offset, setOffset] = useState<GazeOffset>(ZERO)
  const frameRef = useRef<number | null>(null)

  useEffect(() => {
    if (!enabled) {
      setOffset(ZERO)

      return
    }

    if (typeof window === 'undefined' || isTouchOnly()) return

    const onMove = (event: MouseEvent): void => {
      if (frameRef.current !== null) return

      frameRef.current = window.requestAnimationFrame(() => {
        frameRef.current = null
        const node = ref.current
        if (node === null) return

        const box = node.getBoundingClientRect()
        const cx = box.left + box.width / 2
        const cy = box.top + box.height / 2

        setOffset(gazeFor(event.clientX - cx, event.clientY - cy))
      })
    }

    window.addEventListener('mousemove', onMove, { passive: true })

    return () => {
      window.removeEventListener('mousemove', onMove)

      if (frameRef.current !== null) {
        window.cancelAnimationFrame(frameRef.current)
        frameRef.current = null
      }
    }
  }, [enabled, ref])

  return offset
}
