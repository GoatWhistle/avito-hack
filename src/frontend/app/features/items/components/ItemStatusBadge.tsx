import { useTranslation } from 'react-i18next'
import { cn } from '#/lib/utils'
import { statusTone } from '#/features/items/lib'
import type { ItemStatus } from '#/features/items/types'

const toneClasses = {
  neutral: 'bg-muted text-muted-foreground',
  info: 'bg-primary/10 text-primary',
  success: 'bg-primary text-primary-foreground',
  muted: 'bg-muted/60 text-muted-foreground',
} as const

interface ItemStatusBadgeProps {
  status: ItemStatus
  className?: string
}

export function ItemStatusBadge({ status, className }: ItemStatusBadgeProps) {
  const { t } = useTranslation('items')

  return (
    <span
      className={cn(
        'inline-flex shrink-0 items-center rounded-md px-2 py-0.5 text-xs font-medium',
        toneClasses[statusTone(status)],
        className,
      )}
    >
      {t(`status.${status}`)}
    </span>
  )
}
