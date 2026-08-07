import { useCallback, useEffect, useRef, useState } from 'react'

const STROKE_FEEDBACK_MS = 1_200

export interface StrokeState {
  isPetted: boolean
  strokeKey: number
  trigger: () => void
}

export function useStroke(onStroke?: () => void): StrokeState {
  const [isPetted, setIsPetted] = useState(false)
  const [strokeKey, setStrokeKey] = useState(0)
  const strokeRef = useRef(onStroke)
  strokeRef.current = onStroke

  useEffect(() => {
    if (!isPetted) {
      return
    }

    const timer = setTimeout(() => {
      setIsPetted(false)
    }, STROKE_FEEDBACK_MS)

    return () => {
      clearTimeout(timer)
    }
  }, [isPetted, strokeKey])

  const trigger = useCallback(() => {
    setIsPetted(true)
    setStrokeKey((value) => value + 1)
    strokeRef.current?.()
  }, [])

  return { isPetted, strokeKey, trigger }
}
