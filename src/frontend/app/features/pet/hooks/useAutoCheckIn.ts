import { useEffect, useRef } from 'react'
import type { Pet } from '#/features/pet/types'
import type { PetCelebration } from './usePetCelebration'

export const useAutoCheckIn = (
  pet: Pet | undefined,
  celebration: PetCelebration,
): void => {
  const seenRef = useRef<string | null>(null)
  const celebrationRef = useRef(celebration)
  celebrationRef.current = celebration

  useEffect(() => {
    if (pet === undefined) return
    if (pet.checkin_applied !== true) return

    const info = pet.checkin
    if (info === null || info === undefined) return

    const marker = `${pet.id}:${info.streak.days}:${info.xp_granted}`
    if (seenRef.current === marker) return
    seenRef.current = marker

    celebrationRef.current.showCheckIn({
      xp: info.xp_granted,
      days: info.streak.days,
    })

    if (info.level > info.previous_level) {
      celebrationRef.current.celebrate('levelup')
    }
  }, [pet])
}
