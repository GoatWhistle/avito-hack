export interface Leaderboard {
  items: {
    level: number
    name: string
    rank: number
    streak_days: number
    user_id: string
    xp: number
  }[]
  my_rank: number
}
