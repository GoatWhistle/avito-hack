import { ShieldCheck, Store } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { cn } from '#/lib/utils'

interface ItemSourceBadgeProps {
  isSeed: boolean
  aiVerified: boolean
  className?: string
}

export function ItemSourceBadge({
  isSeed,
  aiVerified,
  className,
}: ItemSourceBadgeProps) {
  const { t } = useTranslation('items')

  if (!isSeed && !aiVerified) return null

  const label = isSeed ? t('badge.platform') : t('badge.aiVerified')
  const hint = isSeed ? t('badge.platformHint') : t('badge.aiVerifiedHint')
  const Icon = isSeed ? Store : ShieldCheck

  return (
    <span
      title={hint}
      className={cn(
        'inline-flex shrink-0 items-center gap-1 rounded-md px-2 py-0.5 text-xs font-medium',
        isSeed
          ? 'bg-muted text-muted-foreground'
          : 'bg-primary/10 text-primary',
        className,
      )}
    >
      <Icon className="size-3" aria-hidden="true" />
      {label}
    </span>
  )
}
