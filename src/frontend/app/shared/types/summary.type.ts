import type { PetStage, PetState } from './pet.type'

export interface SummaryAction {
  action: string
  count: number
  amount: number
}

export interface SummaryFacts {
  total_xp: number
  actions: SummaryAction[]
  level: number
  previous_level: number
  leveled_up: boolean
  xp: number
  next_level_xp: number
  xp_to_next_level: number
  rewards: string[]
  badges: string[]
  stage: PetStage
  state: PetState
  satiety: number
  happiness: number
  energy: number
  streak_days: number
  streak_broken: boolean
  issues_count: number
  leaderboard_rank: number
}

export interface SummaryAdvice {
  text: string
  action: string
  item_id?: string
}

export type SummarySource = 'template' | 'llm'

export interface DailySummary {
  id: string
  date: Date
  message: string
  generated_by: SummarySource
  facts: SummaryFacts
  created_at: Date
  advice: SummaryAdvice
}
