import { useGetMyPet } from '#/features/tamagotchi/hooks/useGetMyPet'
import { Badge } from '#/shared/components/ui/badge'
import { Card, CardContent, CardHeader } from '#/shared/components/ui/card'
import { Progress } from '#/shared/components/ui/progress'

export function DigestNextReward() {
  const { data: pet, isLoading, error } = useGetMyPet()

  if (isLoading) {
    return <div></div>
  }

  if (error) {
    return <div></div>
  }

  if (!pet) {
    return <div></div>
  }

  const totalXp = pet.xp + pet.next_level_xp
  const progress = totalXp > 0 ? Math.round((pet.xp / totalXp) * 100) : 100
  const remainingXp = pet.next_level_xp

  return (
    <Card className="rounded-2xl border-emerald-500/20 bg-emerald-50/70 dark:bg-emerald-950/20 shadow-xs relative overflow-hidden">
      <div className="absolute -right-4 -top-4 size-20 rounded-full bg-emerald-500/10 blur-xl pointer-events-none" />

      <CardHeader className="pb-2 pt-4 px-4 flex flex-row items-center justify-between space-y-0">
        <Badge
          variant="outline"
          className="rounded-lg border-emerald-500/30 bg-emerald-100/80 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300 text-xs font-semibold px-2.5 py-0.5"
        >
          Следующий уровень
        </Badge>
        <span className="text-xl">🚀</span>
      </CardHeader>

      <CardContent className="pt-0 px-4 pb-4 space-y-2.5">
        <div>
          <p className="text-sm font-bold text-foreground">
            Уровень {pet.level + 1}
          </p>
          <p className="text-xs text-muted-foreground font-medium mt-0.5">
            Осталось набрать {remainingXp} XP
          </p>
        </div>

        <div className="space-y-1">
          <Progress
            value={progress}
            className="h-2 rounded-full bg-emerald-200/50 dark:bg-emerald-900/50 [&_[data-slot=progress-indicator]]:bg-emerald-600 dark:[&_[data-slot=progress-indicator]]:bg-emerald-400"
          />
          <div className="flex justify-between text-[11px] font-medium text-muted-foreground">
            <span>
              {pet.xp} / {totalXp} XP
            </span>
            <span>{progress}%</span>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
