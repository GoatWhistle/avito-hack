import { useEffect, useRef, useState } from 'react'

const IDLE_TIMEOUT_MS = 45_000
const WAVE_MS = 1_600
const ACTIVITY_EVENTS = [
  'pointerdown',
  'pointermove',
  'keydown',
  'wheel',
] as const

export interface IdleSleep {
  isAsleep: boolean
  isWaving: boolean
}

export function useIdleSleep(enabled: boolean): IdleSleep {
  const [isAsleep, setIsAsleep] = useState(false)
  const [isWaving, setIsWaving] = useState(false)
  const asleepRef = useRef(false)

  useEffect(() => {
    if (!enabled) {
      asleepRef.current = false
      setIsAsleep(false)
      setIsWaving(false)

      return
    }

    let idleTimer: ReturnType<typeof setTimeout> | null = null
    let waveTimer: ReturnType<typeof setTimeout> | null = null

    const schedule = (): void => {
      if (idleTimer !== null) {
        clearTimeout(idleTimer)
      }
      idleTimer = setTimeout(() => {
        asleepRef.current = true
        setIsAsleep(true)
      }, IDLE_TIMEOUT_MS)
    }

    const wake = (): void => {
      if (asleepRef.current) {
        asleepRef.current = false
        setIsAsleep(false)
        setIsWaving(true)

        if (waveTimer !== null) {
          clearTimeout(waveTimer)
        }
        waveTimer = setTimeout(() => {
          setIsWaving(false)
        }, WAVE_MS)
      }
      schedule()
    }

    const onVisibility = (): void => {
      if (document.visibilityState === 'visible') {
        wake()
      }
    }

    for (const event of ACTIVITY_EVENTS) {
      window.addEventListener(event, wake, { passive: true })
    }
    document.addEventListener('visibilitychange', onVisibility)
    schedule()

    return () => {
      if (idleTimer !== null) {
        clearTimeout(idleTimer)
      }
      if (waveTimer !== null) {
        clearTimeout(waveTimer)
      }
      for (const event of ACTIVITY_EVENTS) {
        window.removeEventListener(event, wake)
      }
      document.removeEventListener('visibilitychange', onVisibility)
    }
  }, [enabled])

  return { isAsleep, isWaving }
}
