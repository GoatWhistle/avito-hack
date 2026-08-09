import { AchievementCard } from './AchievementCard'
import type { BadgeItem } from '#/features/rewards/types'

export interface AchievementsGroupProps {
  title: string
  badges: BadgeItem[]
}

export function AchievementsGroup({ title, badges }: AchievementsGroupProps) {
  if (badges.length === 0) return null

  return (
    <div className="flex flex-col gap-2.5">
      <h3 className="px-1 text-xs font-bold tracking-wider text-muted-foreground uppercase">
        {title}
      </h3>
      <ul className="grid grid-cols-2 gap-3">
        {badges.map((badge) => (
          <AchievementCard key={badge.id} badge={badge} />
        ))}
      </ul>
    </div>
  )
}
