import { Apple } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Button } from '#/components/ui'
import { useFeedCooldown } from '#/features/pet/hooks'

export interface FeedButtonProps {
  value: number
  isPending: boolean
  error: string | null
  onFeed: () => void
  availableAt?: string | null
}

const MAX_SATIETY = 100

export function FeedButton({
  value,
  isPending,
  error,
  onFeed,
  availableAt = null,
}: FeedButtonProps) {
  const { t } = useTranslation('pet')
  const cooldown = useFeedCooldown(availableAt)
  const isFull = value >= MAX_SATIETY
  const isDisabled = isPending || isFull || cooldown.isLocked

  const status =
    error ??
    (cooldown.isLocked
      ? t('actions.feedCooldown', { time: cooldown.label })
      : isFull
        ? t('actions.feedFull')
        : '')

  const hint = isPending
    ? t('actions.feeding')
    : status === ''
      ? t('actions.feed')
      : status

  return (
    <>
      {cooldown.isLocked && (
        <span
          aria-hidden="true"
          data-testid="feed-countdown"
          className="shrink-0 font-mono text-[10px] leading-none tabular-nums text-muted-foreground"
        >
          {cooldown.label}
        </span>
      )}

      <Button
        type="button"
        size="icon"
        variant="secondary"
        data-testid="feed-button"
        className="size-7 shrink-0 rounded-lg"
        onClick={onFeed}
        disabled={isDisabled}
        aria-label={`${t('actions.feedAria')}. ${hint}`}
        aria-busy={isPending}
        title={hint}
      >
        <Apple aria-hidden="true" className="size-3.5" />
      </Button>

      <span
        role="status"
        aria-live="polite"
        data-testid="feed-status"
        className="sr-only"
      >
        {isPending ? t('actions.feeding') : status}
      </span>
    </>
  )
}
