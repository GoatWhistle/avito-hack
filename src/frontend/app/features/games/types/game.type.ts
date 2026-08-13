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
  max_attempts?: number
}

export interface GameList {
  games: GameSummary[]
  streak: GameStreak
  daily_done: boolean
}

export interface GameState {
  slug: string
  target_streak: number
  streak: GameStreak
  daily: GameDailyProgress
  active_round?: ActiveGameRound
  max_attempts?: number
  best_score?: number
}

export interface ActiveGameRound {
  round_id: string
  streak: number
  prompt: unknown
  attempts_used?: number
  max_attempts?: number
}

export interface GameRound {
  round_id: string
  streak: number
  target_streak: number
  state?: GameRoundState
  prompt: unknown
  attempts_used?: number
  max_attempts?: number
}

export type GameProgress = 'continue' | 'advance' | 'win' | 'lose'

export interface GameGuessResult {
  correct: boolean
  reveal: unknown
  streak: number
  state: GameRoundState
  prompt?: unknown
  attempt_completed: boolean
  streak_after?: GameStreak
  progress?: GameProgress
  attempts_used?: number
  max_attempts?: number
  best_score?: number
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
  right_item_id?: string
  right_title?: string
}

export type MoreLessChoice = 'higher' | 'lower'

export type BukovkiLetterStatus = 'correct' | 'present' | 'absent'

export interface BukovkiLetter {
  char: string
  status: BukovkiLetterStatus
}

export interface BukovkiPrompt {
  word_length: number
  max_tries: number
  history: BukovkiLetter[][]
}

export interface BukovkiListing {
  display_id: string
  title: string
  price_kopeks: number
  photo_url: string
}

export interface BukovkiReveal {
  feedback: BukovkiLetter[]
  game_over: boolean
  win: boolean
  secret?: string
  listings?: BukovkiListing[]
}

export interface RaccoonJumpCollectible {
  index: number
  title: string
  photo_url: string
}

export interface RaccoonJumpPrompt {
  seed: number
  max_score: number
  min_streak_score: number
  max_score_per_second: number
  collectibles?: RaccoonJumpCollectible[]
  best_score: number
  started_at_unix_milli: number
}

export interface RaccoonJumpMove {
  score: number
  collected: number[]
}

export interface RaccoonJumpReveal {
  score: number
  best_score: number
  new_best: boolean
  counts_toward_streak: boolean
  listings?: BukovkiListing[]
}
