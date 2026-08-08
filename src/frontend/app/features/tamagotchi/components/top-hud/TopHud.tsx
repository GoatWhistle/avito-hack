import { HudChip } from '#/features/tamagotchi/components/top-hud/HudChip'
import { useGetMyRaccoon } from '#/features/tamagotchi/hooks/useGetMyRaccoon'
import { Card } from '#/shared/components/ui/card'

export function TopHud() {
  const { data: pet, isLoading, error } = useGetMyRaccoon()

  if (!pet || isLoading || error) return null

  const totalNeeded = pet.xp + pet.xp_to_next_level
  const progress =
    totalNeeded > 0 ? Math.round((pet.xp / totalNeeded) * 100) : 100
  const points = pet.badges.length

  return (
    <header className="relative z-20 px-4 pt-4 lg:px-8 lg:pt-2.5">
      <Card className="mx-auto max-w-4xl rounded-3xl p-3 lg:max-w-5xl lg:p-4">
        <div className="flex items-center justify-between gap-3">
          <div className="flex min-w-0 items-center gap-3">
            <div className="grid h-12 w-12 shrink-0 place-items-center rounded-2xl bg-emerald-500 text-lg font-black text-white shadow-md lg:h-14 lg:w-14">
              {pet.level}
            </div>

            <div className="min-w-0">
              <p className="truncate text-sm font-bold lg:text-base">
                Ноти · ур. {pet.level}
              </p>

              <p className="truncate text-xs text-slate-500">
                {pet.xp} XP · до ур. {pet.level + 1}: {pet.xp_to_next_level} XP
              </p>

              <div className="mt-2 h-2.5 w-64 overflow-hidden rounded-full bg-slate-200 lg:w-96">
                <div
                  className="h-full rounded-full bg-emerald-500 transition-all duration-500"
                  style={{ width: `${progress}%` }}
                />
              </div>
            </div>
          </div>

          <div className="flex flex-col items-end gap-1.5">
            <HudChip color="amber">⭐ {points} бейджей</HudChip>
            <HudChip color="orange">🔥 {pet.current_streak} дн.</HudChip>
          </div>
        </div>
      </Card>
    </header>
  )
}
