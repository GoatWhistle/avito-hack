import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Check, Copy, Flame, Gift } from 'lucide-react'
import { Button } from '#/components/ui'
import { cn } from '#/lib/utils'
import { useClaimGameReward } from '#/features/games/hooks'
import type { GameStreak } from '#/features/games/types'

const WEEK_LENGTH = 7

interface GameStreakCardProps {
  slug: string
  streak: GameStreak
  countedToday?: boolean
}

export function GameStreakCard({
  slug,
  streak,
  countedToday = false,
}: GameStreakCardProps) {
  const { t } = useTranslation('games')
  const claim = useClaimGameReward(slug)
  const [copied, setCopied] = useState(false)

  const code = claim.data?.code
  const filled = Math.min(streak.current_days, WEEK_LENGTH)

  const copy = async () => {
    if (!code) return
    await navigator.clipboard.writeText(code)
    setCopied(true)
  }

  return (
    <section className="flex flex-col gap-3 rounded-xl bg-card p-3 ring-1 ring-border">
      <header className="flex items-center justify-between gap-3">
        <div className="flex items-center gap-2">
          <Flame className="size-4 text-primary" aria-hidden="true" />
          <h3 className="text-sm font-semibold">{t('streak.title')}</h3>
        </div>
        <p className="text-xs text-muted-foreground">
          {t('streak.best', { count: streak.best_days })}
        </p>
      </header>

      <div className="flex flex-col gap-2">
        <div
          role="img"
          aria-label={t('streak.progress', {
            current: streak.current_days,
            total: WEEK_LENGTH,
          })}
          className="flex gap-1.5"
        >
          {Array.from({ length: WEEK_LENGTH }, (_, index) => (
            <span
              key={index}
              className={cn(
                'h-1.5 flex-1 rounded-full transition-colors',
                index < filled ? 'bg-primary' : 'bg-muted',
              )}
            />
          ))}
        </div>
        <p className="text-xs text-muted-foreground">
          {t('streak.progress', {
            current: streak.current_days,
            total: WEEK_LENGTH,
          })}
          {countedToday && ` · ${t('streak.today')}`}
        </p>
      </div>

      {code ? (
        <div className="flex flex-col gap-2 rounded-lg bg-primary-subtle p-3">
          <p className="text-xs text-primary-subtle-foreground">
            {t('streak.promoCode')}
          </p>
          <div className="flex items-center justify-between gap-2">
            <code className="font-mono text-sm font-semibold">{code}</code>
            <Button type="button" variant="ghost" size="sm" onClick={copy}>
              {copied ? (
                <Check className="size-3.5" aria-hidden="true" />
              ) : (
                <Copy className="size-3.5" aria-hidden="true" />
              )}
              {copied ? t('streak.copied') : t('streak.copy')}
            </Button>
          </div>
        </div>
      ) : streak.reward_ready ? (
        <Button
          type="button"
          size="sm"
          disabled={claim.isPending}
          onClick={() => claim.mutate()}
        >
          <Gift className="size-4" aria-hidden="true" />
          {t('streak.claim')}
        </Button>
      ) : (
        <p className="text-xs text-muted-foreground">{t('streak.rewardHint')}</p>
      )}
    </section>
  )
}
