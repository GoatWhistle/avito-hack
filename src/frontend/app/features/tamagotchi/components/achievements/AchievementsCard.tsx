import {
  Avatar,
  AvatarFallback,
  AvatarImage,
} from '#/shared/components/ui/avatar'
import { Card, CardContent } from '#/shared/components/ui/card'
import { cn } from '#/shared/lib/utils'
import type { Badge } from '#/shared/types/badge.type'

interface Props {
  badge: Badge
}

export function AchievementsCard({ badge }: Props) {
  const done = !!badge.earned_at

  return (
    <Card
      className={cn(
        'rounded-2xl border-border/60 bg-card shadow-xs transition-all hover:border-border flex flex-col justify-between',
        done
          ? 'border-emerald-500/30 bg-emerald-50/50 dark:border-emerald-500/30 dark:bg-emerald-950/20'
          : 'opacity-75',
      )}
    >
      <CardContent className="p-4 flex flex-col h-full justify-between gap-4">
        <div>
          <Avatar
            className={cn(
              'size-11 rounded-xl border shadow-none',
              done
                ? 'border-emerald-500/30 bg-white dark:border-emerald-700/50 dark:bg-card'
                : 'border-border/60 bg-muted/50',
            )}
          >
            <AvatarImage src={badge.icon_url} alt={badge.name} />
            <AvatarFallback
              className={cn(
                'rounded-xl bg-transparent text-lg',
                !done && 'grayscale opacity-70',
              )}
            >
              {done ? '🎖️' : '🔒'}
            </AvatarFallback>
          </Avatar>

          <p className="mt-3 text-sm font-semibold text-foreground leading-tight">
            {badge.name}
          </p>

          <p className="mt-1 line-clamp-2 text-xs text-muted-foreground font-medium">
            {badge.description}
          </p>
        </div>
      </CardContent>
    </Card>
  )
}
