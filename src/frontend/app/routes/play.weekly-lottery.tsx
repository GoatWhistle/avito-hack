import { RequireAuth } from '#/features/auth/components'
import { WeeklyLotteryScreen } from '#/features/weekly-lottery'

export default function WeeklyLotteryRoute() {
  return (
    <RequireAuth>
      <WeeklyLotteryScreen />
    </RequireAuth>
  )
}
