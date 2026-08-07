import { useCallback, useEffect, useRef, useState } from 'react'
import type { PetEmotion } from '#/components/pet-avatar'
import type { PetEvent } from '#/features/pet/types'

const XP_TOAST_MS = 1_800

export interface XpToast {
  id: number
  amount: number
}

export interface CelebrationBanner {
  kind: 'levelUp' | 'hatched' | 'reward' | 'streak'
  level?: number
  title?: string
  days?: number
}

export interface PetCelebration {
  emotion: PetEmotion | null
  xpToasts: XpToast[]
  banner: CelebrationBanner | null
  handleEvent: (event: PetEvent) => void
  clearEmotion: () => void
  dismissBanner: () => void
  celebrate: (emotion: PetEmotion) => void
}

export const usePetCelebration = (): PetCelebration => {
  const [emotion, setEmotion] = useState<PetEmotion | null>(null)
  const [xpToasts, setXpToasts] = useState<XpToast[]>([])
  const [banner, setBanner] = useState<CelebrationBanner | null>(null)
  const idRef = useRef(0)
  const timersRef = useRef<ReturnType<typeof setTimeout>[]>([])

  useEffect(
    () => () => {
      for (const timer of timersRef.current) clearTimeout(timer)
    },
    [],
  )

  const pushXp = useCallback((amount: number) => {
    idRef.current += 1
    const id = idRef.current
    setXpToasts((current) => [...current, { id, amount }])

    const timer = setTimeout(() => {
      setXpToasts((current) => current.filter((toast) => toast.id !== id))
    }, XP_TOAST_MS)
    timersRef.current.push(timer)
  }, [])

  const handleEvent = useCallback(
    (event: PetEvent) => {
      switch (event.type) {
        case 'xp.gained':
          pushXp(event.payload.amount)
          setEmotion('celebrate')

          return
        case 'level.up':
          setEmotion('levelup')
          setBanner({ kind: 'levelUp', level: event.payload.level })

          return
        case 'pet.hatched':
          setEmotion('hatching')
          setBanner({ kind: 'hatched' })

          return
        case 'reward.granted':
          setBanner({ kind: 'reward', title: event.payload.title })

          return
        case 'streak.updated':
          if (event.payload.milestone === true) {
            setBanner({ kind: 'streak', days: event.payload.days })
          }

          return
        default:
          return
      }
    },
    [pushXp],
  )

  const clearEmotion = useCallback(() => setEmotion(null), [])
  const dismissBanner = useCallback(() => setBanner(null), [])
  const celebrate = useCallback((next: PetEmotion) => setEmotion(next), [])

  return {
    emotion,
    xpToasts,
    banner,
    handleEvent,
    clearEmotion,
    dismissBanner,
    celebrate,
  }
}
