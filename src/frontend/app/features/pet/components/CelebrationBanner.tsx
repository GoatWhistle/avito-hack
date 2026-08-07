import { useTranslation } from 'react-i18next'
import { Button } from '#/components/ui'
import type { CelebrationBanner as BannerData } from '#/features/pet/hooks'

const ICON: Record<BannerData['kind'], string> = {
  levelUp: '🎉',
  hatched: '🐣',
  reward: '🎁',
  streak: '🔥',
}

export interface CelebrationBannerProps {
  banner: BannerData
  onDismiss: () => void
}

export function CelebrationBanner({
  banner,
  onDismiss,
}: CelebrationBannerProps) {
  const { t } = useTranslation('pet')

  const text =
    banner.kind === 'levelUp'
      ? t('events.levelUp', { level: banner.level ?? 0 })
      : banner.kind === 'hatched'
        ? t('events.hatched')
        : banner.kind === 'reward'
          ? t('events.rewardGranted', { title: banner.title ?? '' })
          : t('streak.milestone', { count: banner.days ?? 0 })

  return (
    <div
      role="status"
      data-testid="celebration-banner"
      className="animate-in fade-in slide-in-from-top-2 flex items-center gap-3 rounded-xl bg-linear-to-r from-primary to-accent px-4 py-3 text-primary-foreground shadow-lg duration-300"
    >
      <span aria-hidden="true" className="text-xl leading-none">
        {ICON[banner.kind]}
      </span>
      <p className="min-w-0 flex-1 text-sm font-semibold">{text}</p>
      <Button
        size="icon-sm"
        variant="ghost"
        onClick={onDismiss}
        aria-label={t('actions.dismiss')}
        className="shrink-0 text-primary-foreground hover:bg-primary-foreground/20"
      >
        <span aria-hidden="true">✕</span>
      </Button>
    </div>
  )
}
