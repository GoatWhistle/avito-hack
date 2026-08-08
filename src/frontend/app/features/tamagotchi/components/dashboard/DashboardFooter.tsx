import type { DashboardFooterItem } from '#/features/tamagotchi/components/dashboard/left-sidebar-footer-item.types'
import { Progress } from '#/shared/components/ui/progress'
import { SidebarFooter } from '#/shared/components/ui/sidebar'
import { cn } from '#/shared/lib/utils'
import type { Pet } from '#/shared/types/pet.type'
import { useMemo } from 'react'

const COLOR_MAP = {
  energy: {
    icon: 'bg-blue-500/10 text-blue-600 dark:bg-blue-500/20 dark:text-blue-400',
    progress: '[&_[data-slot=progress-indicator]]:bg-blue-500',
  },
  happiness: {
    icon: 'bg-yellow-500/10 text-yellow-600 dark:bg-yellow-500/20 dark:text-yellow-400',
    progress: '[&_[data-slot=progress-indicator]]:bg-yellow-500',
  },
  satiety: {
    icon: 'bg-green-500/10 text-green-600 dark:bg-green-500/20 dark:text-green-400',
    progress: '[&_[data-slot=progress-indicator]]:bg-green-500',
  },
  streak: {
    icon: 'bg-orange-500/10 text-orange-600 dark:bg-orange-500/20 dark:text-orange-400',
    progress: '[&_[data-slot=progress-indicator]]:bg-orange-500',
  },
  xp: {
    icon: 'bg-purple-500/10 text-purple-600 dark:bg-purple-500/20 dark:text-purple-400',
    progress: '[&_[data-slot=progress-indicator]]:bg-purple-500',
  },
} as const

interface Props {
  pet: Pet
}

export function DashboardFooter({ pet }: Props) {
  const stats: DashboardFooterItem[] = useMemo(() => {
    const totalXp = pet.xp + pet.next_level_xp

    return [
      {
        key: 'energy' as const,
        icon: '⚡',
        name: 'Энергия',
        value: pet.energy,
        percent: pet.energy,
      },
      {
        key: 'happiness' as const,
        icon: '😊',
        name: 'Счастье',
        value: pet.happiness,
        percent: pet.happiness,
      },
      {
        key: 'satiety' as const,
        icon: '🍎',
        name: 'Сытость',
        value: pet.satiety,
        percent: pet.satiety,
      },
      {
        key: 'streak' as const,
        icon: '🔥',
        name: 'Стрик (дней)',
        value: pet.streak_days,
        percent: Math.min((pet.streak_days / 30) * 100, 100),
      },
      {
        key: 'xp' as const,
        icon: '⭐',
        name: 'Опыт',
        value: `${pet.xp} / ${totalXp}`,
        percent: totalXp > 0 ? (pet.xp / totalXp) * 100 : 0,
      },
    ]
  }, [
    pet.energy,
    pet.happiness,
    pet.satiety,
    pet.streak_days,
    pet.xp,
    pet.next_level_xp,
  ])

  return (
    <SidebarFooter className="border-t border-sidebar-border/60 p-3">
      <div className="space-y-1.5">
        {stats.map(stat => {
          const colors = COLOR_MAP[stat.key as keyof typeof COLOR_MAP]
          return (
            <div
              key={stat.key}
              className="flex items-center gap-3 rounded-xl px-3 py-2.5 transition-colors hover:bg-sidebar-accent/80"
            >
              <div
                className={cn(
                  'flex size-9 shrink-0 items-center justify-center rounded-xl text-base',
                  colors.icon,
                )}
              >
                {stat.icon}
              </div>
              <div className="min-w-0 flex-1">
                <div className="flex items-baseline justify-between gap-2">
                  <span className="truncate text-xs text-muted-foreground font-medium">
                    {stat.name}
                  </span>
                  <span className="text-sm font-bold tabular-nums text-foreground">
                    {stat.value}
                  </span>
                </div>
                <div className="mt-2">
                  <Progress
                    value={stat.percent}
                    className={cn(
                      'h-1.5 rounded-full bg-muted/70',
                      colors.progress,
                    )}
                  />
                </div>
              </div>
            </div>
          )
        })}
      </div>
    </SidebarFooter>
  )
}
