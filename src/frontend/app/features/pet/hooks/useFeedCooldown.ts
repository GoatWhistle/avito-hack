import { useEffect, useState } from 'react'

const SECOND_MS = 1_000

export interface FeedCooldown {
  isLocked: boolean
  remainingMs: number
  label: string
}

const pad = (value: number): string => String(value).padStart(2, '0')

export const formatCountdown = (remainingMs: number): string => {
  const total = Math.max(0, Math.ceil(remainingMs / SECOND_MS))
  const hours = Math.floor(total / 3_600)
  const minutes = Math.floor((total % 3_600) / 60)
  const seconds = total % 60

  return `${hours}:${pad(minutes)}:${pad(seconds)}`
}

const remainingFrom = (
  availableAt: string | null | undefined,
  now: number,
): number => {
  if (availableAt === null || availableAt === undefined || availableAt === '') {
    return 0
  }

  const target = new Date(availableAt).getTime()
  if (Number.isNaN(target)) return 0

  return Math.max(0, target - now)
}

export const useFeedCooldown = (
  availableAt: string | null | undefined,
): FeedCooldown => {
  const [remainingMs, setRemainingMs] = useState(() =>
    remainingFrom(availableAt, Date.now()),
  )

  useEffect(() => {
    const initial = remainingFrom(availableAt, Date.now())
    setRemainingMs(initial)

    if (initial <= 0) return

    const timer = setInterval(() => {
      const next = remainingFrom(availableAt, Date.now())
      setRemainingMs(next)

      if (next <= 0) clearInterval(timer)
    }, SECOND_MS)

    return () => {
      clearInterval(timer)
    }
  }, [availableAt])

  return {
    isLocked: remainingMs > 0,
    remainingMs,
    label: formatCountdown(remainingMs),
  }
}
