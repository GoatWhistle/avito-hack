import { Badge } from '#/shared/components/ui/badge'
import { Button } from '#/shared/components/ui/button'
import { Card, CardContent } from '#/shared/components/ui/card'
import { cn } from '#/shared/lib/utils'
import type {
  Reward,
  RewardConditionType,
  RewardKind,
} from '#/shared/types/reward.type'
import {
  CheckCircle2,
  Gift,
  LockIcon,
  Palette,
  Percent,
  Wrench,
} from 'lucide-react'
import type { ComponentType } from 'react'

interface AwardsLevelCardProps {
  reward: Reward
  onActivate: (rewardId: string) => void
  isActivating?: boolean
}

const CONDITION_LABEL: Record<RewardConditionType, string> = {
  level: 'ур.',
  streak: 'дн.',
  achievement: 'дост.',
}

const KIND_ICON: Record<RewardKind, ComponentType<{ className?: string }>> = {
  promo: Percent,
  cosmetic: Palette,
  utility: Wrench,
}

const KIND_COLOR: Record<RewardKind, string> = {
  promo: 'text-purple-500',
  cosmetic: 'text-pink-500',
  utility: 'text-blue-500',
}

function getCardState(reward: Reward): 'locked' | 'current' | 'completed' {
  if (reward.claimed) return 'completed'
  if (reward.unlocked || reward.progress_current >= reward.progress_target)
    return 'current'
  return 'locked'
}

export function AwardsLevelCard({
  reward,
  onActivate,
  isActivating,
}: AwardsLevelCardProps) {
  const state = getCardState(reward)
  const isCompleted = state === 'completed'
  const isLocked = state === 'locked'
  const canActivate = state === 'current' && !isActivating

  const unitLabel = CONDITION_LABEL[reward.condition_type] ?? ''
  const KindIcon = KIND_ICON[reward.kind]

  return (
    <Card
      className={cn(
        'transition-all',
        state === 'current' &&
          'border-emerald-300 ring-2 ring-emerald-200 dark:bg-emerald-950',
        isCompleted &&
          'border-emerald-200 bg-emerald-50 dark:bg-emerald-950/40',
        isLocked && 'opacity-70',
      )}
    >
      <CardContent className="flex items-start justify-between gap-3 p-4">
        <div className="flex min-w-0 flex-1 items-start gap-3">
          <div
            className={cn(
              'grid size-10 shrink-0 place-items-center rounded-xl text-sm font-black mt-0.5',
              isLocked
                ? 'bg-muted text-muted-foreground'
                : 'bg-emerald-500 text-white',
            )}
          >
            {isCompleted ? (
              <CheckCircle2 className="size-5" />
            ) : (
              reward.progress_target
            )}
          </div>

          <div className="min-w-0 flex-1">
            <div className="flex items-center gap-1.5">
              <KindIcon
                className={cn('size-3.5 shrink-0', KIND_COLOR[reward.kind])}
              />
              <p className="flex items-center gap-1 text-sm font-bold leading-tight break-words">
                {isLocked && <LockIcon className="size-3 shrink-0" />}
                <span className="break-words">{reward.title}</span>
              </p>
            </div>
            <p className="mt-1 text-xs text-muted-foreground leading-tight line-clamp-2 break-words">
              {reward.description}
            </p>
          </div>
        </div>

        <div className="flex shrink-0 items-center self-center">
          {isCompleted ? (
            <Badge variant="default" className="font-bold whitespace-nowrap">
              Получено
            </Badge>
          ) : isLocked ? (
            <Badge variant="outline" className="font-bold whitespace-nowrap">
              {reward.progress_current}/{reward.progress_target} {unitLabel}
            </Badge>
          ) : (
            <Button
              size="sm"
              className="h-8 gap-1.5 whitespace-nowrap"
              disabled={!canActivate}
              onClick={() => onActivate(reward.id)}
            >
              {isActivating ? (
                <span className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
              ) : (
                <Gift className="size-3.5" />
              )}
              Получить
            </Button>
          )}
        </div>
      </CardContent>
    </Card>
  )
}
