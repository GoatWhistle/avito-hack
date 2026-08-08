import { AwardsLevelCard } from '#/features/tamagotchi/components/awards/AwardsLevelCard'
import { useActivateReward } from '#/features/tamagotchi/hooks/useActivateReward'
import { useFindRewards } from '#/features/tamagotchi/hooks/useFindRewards'
import { Button } from '#/shared/components/ui/button'
import { Loader2 } from 'lucide-react'

export function AwardsLevelList() {
  const { data: rewards, isLoading, error, refetch } = useFindRewards()
  const { mutate: activateReward, isPending: isActivating } =
    useActivateReward()

  if (isLoading) {
    return (
      <div className="flex justify-center p-4">
        <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
      </div>
    )
  }

  if (error || !rewards) {
    return (
      <div className="space-y-2 p-4 text-center">
        <p className="text-red-500">Ошибка загрузки наград</p>
        <Button variant="outline" size="sm" onClick={() => refetch()}>
          Повторить
        </Button>
      </div>
    )
  }

  if (rewards.items.length === 0) {
    return (
      <div className="p-4 text-center text-muted-foreground">
        У вас пока нет доступных наград
      </div>
    )
  }

  const handleActivate = (rewardId: string) => {
    activateReward(rewardId)
  }

  return (
    <div className="space-y-2.5">
      {rewards.items.map(reward => (
        <AwardsLevelCard
          key={reward.id}
          reward={reward}
          onActivate={handleActivate}
          isActivating={isActivating}
        />
      ))}
    </div>
  )
}
