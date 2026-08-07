import { useEffect, useState } from 'react'

const BLINK_MIN_MS = 3_200
const BLINK_MAX_MS = 6_400
const BLINK_DURATION_MS = 140

function nextDelay(): number {
  return BLINK_MIN_MS + Math.random() * (BLINK_MAX_MS - BLINK_MIN_MS)
}

export function useBlink(enabled: boolean): boolean {
  const [closed, setClosed] = useState(false)

  useEffect(() => {
    if (!enabled) {
      setClosed(false)

      return
    }

    let openTimer: ReturnType<typeof setTimeout> | null = null
    let closeTimer: ReturnType<typeof setTimeout> | null = null

    const schedule = (): void => {
      closeTimer = setTimeout(() => {
        setClosed(true)
        openTimer = setTimeout(() => {
          setClosed(false)
          schedule()
        }, BLINK_DURATION_MS)
      }, nextDelay())
    }

    schedule()

    return () => {
      if (openTimer !== null) {
        clearTimeout(openTimer)
      }
      if (closeTimer !== null) {
        clearTimeout(closeTimer)
      }
      setClosed(false)
    }
  }, [enabled])

  return closed
}
