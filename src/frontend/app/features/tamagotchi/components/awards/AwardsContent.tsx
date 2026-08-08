import { TamagotchiDashboardContent } from '#/features/tamagotchi/components/sidebar/TamagotchiDashboardContent'
import { AwardsLevelList } from './AwardsLevelList'

export function AwardsContent() {
  return <TamagotchiDashboardContent children={<AwardsLevelList />} />
}
