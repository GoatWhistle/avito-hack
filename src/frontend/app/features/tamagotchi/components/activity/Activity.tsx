import { stageMap } from '#/features/tamagotchi/components/dashboard/DashboardHeader'
import { Raccoon } from '#/features/tamagotchi/components/raccoon/Raccoon'
import {
  PetStage,
  StreakTracker,
} from '#/features/tamagotchi/components/top-hud/Stage'
import { TopHud } from '#/features/tamagotchi/components/top-hud/TopHud'
import { useCheckinMyPet } from '#/features/tamagotchi/hooks/useCheckinMyPet'
import { useGetMyPet } from '#/features/tamagotchi/hooks/useGetMyPet'
import { useStrokeMyPet } from '#/features/tamagotchi/hooks/useStrokeMyPet'

export function Activity() {
  const { data, isLoading, error } = useGetMyPet()
  const { mutateAsync: stroke } = useStrokeMyPet()
  const { mutateAsync: checkin } = useCheckinMyPet()

  if (!data || isLoading || error) return null

  const handlePet = async () => {
    await stroke()
    await checkin()
  }

  return (
    <div className="flex flex-col flex-1 items-center justify-start h-full w-full gap-4">
      <TopHud />
      <div className="max-w-xs rounded-3xl border border-white/70 bg-white/80 px-4 py-3 text-center text-sm font-medium shadow-lg backdrop-blur">
        Привет! Я Енотик Ноти.
        <br />
        Моя стадия: <b>{stageMap[data.stage]}</b>
      </div>
      <Raccoon />
      <PetStage pet={data} onPet={handlePet} />
      <StreakTracker streakDays={data.streak_days} />
    </div>
  )
}
