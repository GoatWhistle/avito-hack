import type { Badge } from '#/shared/types/badge.type'
import { AchievementsCard } from './AchievementsCard'

interface Props {
  badges: Badge[]
  title?: string
}

export function AchievementsGroupSection({ badges, title }: Props) {
  if (!badges || badges.length === 0) return null

  return (
    <div className="space-y-2.5">
      {title && (
        <h3 className="px-1 text-xs font-bold uppercase tracking-wider text-muted-foreground/80">
          {title}
        </h3>
      )}
      <div className="grid grid-cols-2 gap-3">
        {badges.map(badge => (
          <AchievementsCard key={badge.id} badge={badge} />
        ))}
      </div>
    </div>
  )
}
