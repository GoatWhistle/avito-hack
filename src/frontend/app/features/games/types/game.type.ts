export type GameRoundState = 'active' | 'won' | 'lost'

export interface GameStreak {
  current_days: number
  best_days: number
  reward_ready: boolean
}

export interface GameDailyProgress {
  attempts: number
  best_streak: number
}

export interface GameSummary {
  slug: string
  target_streak: number
  daily_done: boolean
  streak: GameStreak
}

export interface GameState {
  slug: string
  target_streak: number
  streak: GameStreak
  daily: GameDailyProgress
  active_round?: ActiveGameRound
}

export interface ActiveGameRound {
  round_id: string
  streak: number
  prompt: unknown
}

export interface GameRound {
  round_id: string
  streak: number
  target_streak: number
  state?: GameRoundState
  prompt: unknown
}

export interface GameGuessResult {
  correct: boolean
  reveal: unknown
  streak: number
  state: GameRoundState
  prompt?: unknown
  attempt_completed: boolean
  streak_after?: GameStreak
}

export interface GameRewardCode {
  code: string
}

export interface MoreLessItem {
  item_id: string
  title: string
  photo_url: string
  price?: number
}

export interface MoreLessPrompt {
  left: MoreLessItem
  right: MoreLessItem
}

export interface MoreLessReveal {
  right_price: number
}

export type MoreLessChoice = 'higher' | 'lower'
