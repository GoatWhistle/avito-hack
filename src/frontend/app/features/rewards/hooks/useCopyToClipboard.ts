import { useCallback, useEffect, useRef, useState } from 'react'

export const useCopyToClipboard = (resetAfterMs = 2000) => {
  const [copied, setCopied] = useState(false)
  const timer = useRef<ReturnType<typeof setTimeout> | null>(null)

  useEffect(
    () => () => {
      if (timer.current) clearTimeout(timer.current)
    },
    [],
  )

  const copy = useCallback(
    async (value: string) => {
      try {
        await navigator.clipboard.writeText(value)
        setCopied(true)
        if (timer.current) clearTimeout(timer.current)
        timer.current = setTimeout(() => setCopied(false), resetAfterMs)

        return true
      } catch {
        setCopied(false)

        return false
      }
    },
    [resetAfterMs],
  )

  return { copied, copy }
}
