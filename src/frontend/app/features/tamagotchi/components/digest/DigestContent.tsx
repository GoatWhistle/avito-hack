import { TamagotchiDashboardContent } from '#/features/tamagotchi/components/sidebar/TamagotchiDashboardContent'
import { DigestDailyTasks } from './DigestDailyTasks'
import { DigestNextReward } from './DigestNextReward'

export function DigestContent() {
  return (
    <TamagotchiDashboardContent
      children={
        <>
          <DigestNextReward />
          <DigestDailyTasks />
        </>
      }
    />
  )
}
