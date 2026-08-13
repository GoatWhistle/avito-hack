import { GamesScreen } from '#/features/games'
import { WeeklyLotteryCard } from '#/features/weekly-lottery'

export default function PetGamesRoute() {
  return <GamesScreen featured={<WeeklyLotteryCard />} />
}
