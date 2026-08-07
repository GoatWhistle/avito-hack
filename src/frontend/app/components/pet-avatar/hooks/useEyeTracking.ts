import { useEffect, useRef, useState, type RefObject } from 'react'
import type { EyeOffset } from '../types'

const MAX_OFFSET_PX = 3
const REACH_PX = 320
const ZERO: EyeOffset = { x: 0, y: 0 }

function clampOffset(value: number): number {
  return Math.max(-MAX_OFFSET_PX, Math.min(MAX_OFFSET_PX, value))
}

export interface EyeTracking<T extends Element> {
  ref: RefObject<T | null>
  offset: EyeOffset
}

export function useEyeTracking<T extends Element>(
  enabled: boolean,
): EyeTracking<T> {
  const ref = useRef<T>(null)
  const [offset, setOffset] = useState<EyeOffset>(ZERO)

  useEffect(() => {
    if (!enabled) {
      setOffset(ZERO)

      return
    }

    let frame: number | null = null

    const onMove = (event: PointerEvent): void => {
      if (frame !== null) {
        return
      }

      frame = requestAnimationFrame(() => {
        frame = null
        const node = ref.current

        if (node === null) {
          return
        }

        const rect = node.getBoundingClientRect()
        const centerX = rect.left + rect.width / 2
        const centerY = rect.top + rect.height / 2

        setOffset({
          x: clampOffset(
            ((event.clientX - centerX) / REACH_PX) * MAX_OFFSET_PX,
          ),
          y: clampOffset(
            ((event.clientY - centerY) / REACH_PX) * MAX_OFFSET_PX,
          ),
        })
      })
    }

    window.addEventListener('pointermove', onMove, { passive: true })

    return () => {
      if (frame !== null) {
        cancelAnimationFrame(frame)
      }
      window.removeEventListener('pointermove', onMove)
    }
  }, [enabled])

  return { ref, offset }
}
