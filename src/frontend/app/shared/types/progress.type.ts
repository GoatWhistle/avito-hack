import type { PetStage } from './pet.type'

export interface Progress {
  level: number
  xp: number
  next_level_xp: number
  xp_to_next_level: number
  stage: PetStage
  streak_days: number
  freezes: number
  is_max_level: boolean
}
