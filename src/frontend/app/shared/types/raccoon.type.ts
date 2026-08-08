import type { Badge } from './badge.type'
import type { PetStage, PetState } from './pet.type'

export interface Raccoon {
  id: string
  name: string
  level: number
  xp: number
  xp_to_next_level: number
  stage: PetStage
  state: PetState
  current_streak: number
  badges: Badge[]
}
