import { useEffect, useRef, useState } from 'react'
import type { PetEmotion } from '../types'

export const EMOTION_RESET_MS = 2_500

export function useTransientEmotion(
  emotion: PetEmotion | null | undefined,
  resetMs: number,
  onEmotionEnd?: () => void,
): PetEmotion | null {
  const [active, setActive] = useState<PetEmotion | null>(emotion ?? null)
  const endRef = useRef(onEmotionEnd)
  endRef.current = onEmotionEnd

  useEffect(() => {
    const next = emotion ?? null
    setActive(next)

    if (next === null) {
      return
    }

    const timer = setTimeout(() => {
      setActive(null)
      endRef.current?.()
    }, resetMs)

    return () => {
      clearTimeout(timer)
    }
  }, [emotion, resetMs])

  return active
}
