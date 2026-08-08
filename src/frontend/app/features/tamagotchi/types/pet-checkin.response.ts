import type { Pet } from '#/shared/types/pet.type'
import type { Streak } from '#/shared/types/streak.type'

export interface PetCheckinResponse {
  level: number
  nextLevelXp: number
  previousLevel: number
  pet: Pet
  streak: Streak
  unlockedRewards: string[]
  xpGranted: number
}
