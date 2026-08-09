import { CheckCircle2, Flame, Gift, PartyPopper, X } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Button } from '#/components/ui'
import type { CelebrationBanner as BannerData } from '#/features/pet/hooks'

const ICON = {
  levelUp: PartyPopper,
  reward: Gift,
  streak: Flame,
  checkIn: CheckCircle2,
} as const

export interface CelebrationBannerProps {
  banner: BannerData
  onDismiss: () => void
}

export function CelebrationBanner({
  banner,
  onDismiss,
}: CelebrationBannerProps) {
  const { t } = useTranslation('pet')
  const Icon = ICON[banner.kind]

  const text =
    banner.kind === 'levelUp'
      ? t('events.levelUp', { level: banner.level ?? 0 })
      : banner.kind === 'checkIn'
        ? t('events.checkInApplied', {
            count: banner.days ?? 0,
            xp: banner.xp ?? 0,
          })
        : banner.kind === 'reward'
          ? t('events.rewardGranted', { title: banner.title ?? '' })
          : t('streak.milestone', { count: banner.days ?? 0 })

  return (
    <div
      role="status"
      data-testid="celebration-banner"
      className="animate-in fade-in slide-in-from-top-2 flex items-center gap-3 rounded-xl bg-linear-to-r from-primary to-accent px-4 py-3 text-primary-foreground shadow-lg duration-300"
    >
      <Icon aria-hidden="true" className="size-5 shrink-0" />
      <p className="min-w-0 flex-1 text-sm font-semibold">{text}</p>
      <Button
        size="icon-sm"
        variant="ghost"
        onClick={onDismiss}
        aria-label={t('actions.dismiss')}
        className="shrink-0 text-primary-foreground hover:bg-primary-foreground/20"
      >
        <X aria-hidden="true" className="size-4" />
      </Button>
    </div>
  )
}
