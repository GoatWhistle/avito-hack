import { useFindMyBadges } from '#/features/tamagotchi/hooks/useFindMyBadges'
import { TamagotchiDashboardContent } from '../sidebar/TamagotchiDashboardContent'
import { AchievementsGroupSection } from './AchievementsGroupSection'

export function AchievementsContent() {
  const { data, isLoading, error } = useFindMyBadges()

  if (isLoading) {
    return <TamagotchiDashboardContent>Загрузка...</TamagotchiDashboardContent>
  }

  if (error || !data) {
    return (
      <TamagotchiDashboardContent>Ошибка загрузки</TamagotchiDashboardContent>
    )
  }

  if (data.length === 0) {
    return (
      <TamagotchiDashboardContent>Нет достижений</TamagotchiDashboardContent>
    )
  }

  return (
    <TamagotchiDashboardContent>
      {data.map(badge => (
        <AchievementsGroupSection key={badge.id} badges={data} />
      ))}
    </TamagotchiDashboardContent>
  )
}
