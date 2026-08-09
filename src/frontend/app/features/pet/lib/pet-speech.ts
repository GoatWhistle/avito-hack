import type { Pet } from '#/features/pet/types'
import { MAX_LEVEL, STAT_CRITICAL_BELOW, STAT_WARNING_BELOW } from './pet-stats'

export type PetSpeechKey =
  'hungry' | 'tired' | 'lonely' | 'checkIn' | 'streak' | 'legend' | 'happy'

export interface PetSpeechInput {
  pet: Pet
  canCheckIn: boolean
}

export const petSpeech = ({
  pet,
  canCheckIn,
}: PetSpeechInput): PetSpeechKey => {
  if (pet.satiety < STAT_CRITICAL_BELOW) return 'hungry'
  if (pet.energy < STAT_CRITICAL_BELOW) return 'tired'
  if (pet.happiness < STAT_WARNING_BELOW) return 'lonely'
  if (canCheckIn) return 'checkIn'
  if (pet.streak_days >= 7) return 'streak'
  if (pet.level >= MAX_LEVEL) return 'legend'

  return 'happy'
}
