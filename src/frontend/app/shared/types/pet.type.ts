export type PetStage = 'egg' | 'baby' | 'teen' | 'adult' | 'legend'
export type PetState = 'happy' | 'neutral' | 'sad' | 'sleeping'

export interface Pet {
  id: string
  name: string
  level: number
  xp: number
  next_level_xp: number
  stage: PetStage
  state: PetState
  streak_days: number
  energy: number
  happiness: number
  satiety: number
  is_hatched: boolean
  hatched_at: Date
  last_checkin_date: Date
  last_decay_time: Date
  updated_at: Date
  freezes: number
}
