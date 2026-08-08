import { Button } from '#/shared/components/ui/button'
import type { Pet } from '#/shared/types/pet.type'

interface PetStageProps {
  pet: Pet
  onPet?: () => void
}

function ActionBtn({
  icon,
  label,
  sub,
  onClick,
}: {
  icon: string
  label: string
  sub: string
  onClick?: () => void
}) {
  return (
    <Button
      onClick={onClick}
      variant="secondary"
      className="flex flex-col items-center gap-1 rounded-3xl border border-white/70 bg-white/75 px-8 py-10 shadow-lg backdrop-blur active:scale-95"
    >
      <span className="text-2xl">{icon}</span>
      <span className="text-xs font-bold text-slate-600">{label}</span>
      <span className="text-[10px] text-slate-400">{sub}</span>
    </Button>
  )
}

export function PetStage({ pet, onPet }: PetStageProps) {
  return (
    <div className="relative z-10 flex flex-1 flex-col items-center justify-center gap-5 px-6">
      <div className="w-full max-w-md gap-3">
        <ActionBtn
          icon="🤚"
          label="Погладить"
          sub={`настроение: ${pet.happiness}`}
          onClick={onPet}
        />
      </div>
    </div>
  )
}

interface StreakTrackerProps {
  streakDays: number
  targetStreak?: number
}

export function StreakTracker({
  streakDays,
  targetStreak = 7,
}: StreakTrackerProps) {
  return (
    <div className="hidden rounded-3xl border border-white/70 bg-white/70 px-4 py-3 shadow backdrop-blur lg:block">
      <p className="mb-2 text-center text-xs font-bold text-slate-500">Стрик</p>
      <div className="flex items-center gap-2">
        {Array.from({ length: targetStreak }).map((_, i) => {
          const day = i + 1
          const isActive = day <= streakDays
          return (
            <div
              key={i}
              className={`grid h-9 w-9 place-items-center rounded-full text-xs font-black ${
                isActive
                  ? 'bg-orange-400 text-white'
                  : 'bg-slate-200 text-slate-400'
              }`}
            >
              {day}
            </div>
          )
        })}
      </div>
    </div>
  )
}
