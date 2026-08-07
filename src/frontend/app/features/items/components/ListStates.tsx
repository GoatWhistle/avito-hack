import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { Button } from '#/components/ui'
import { cn } from '#/lib/utils'
import type { ReactNode } from 'react'

interface ItemsSkeletonProps {
  count?: number
  className?: string
}

export function ItemsSkeleton({ count = 6, className }: ItemsSkeletonProps) {
  const { t } = useTranslation()

  return (
    <div
      role="status"
      aria-live="polite"
      aria-busy="true"
      aria-label={t('status.loading')}
      className={cn(
        'grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3',
        className,
      )}
    >
      {Array.from({ length: count }, (_, index) => (
        <div
          key={index}
          className="flex flex-col gap-3 rounded-xl bg-card p-4 ring-1 ring-foreground/10"
        >
          <div className="aspect-4/3 w-full animate-pulse rounded-lg bg-muted" />
          <div className="h-4 w-3/4 animate-pulse rounded bg-muted" />
          <div className="h-4 w-1/3 animate-pulse rounded bg-muted" />
        </div>
      ))}
    </div>
  )
}

interface ErrorStateProps {
  message?: string
  onRetry?: () => void
}

export function ErrorState({ message, onRetry }: ErrorStateProps) {
  const { t } = useTranslation()

  return (
    <div
      role="alert"
      className="flex flex-col items-center gap-3 rounded-xl bg-card px-6 py-10 text-center ring-1 ring-foreground/10"
    >
      <p className="text-sm font-medium">{t('status.error')}</p>
      {message && (
        <p className="max-w-prose text-sm text-muted-foreground">{message}</p>
      )}
      {onRetry && (
        <Button variant="outline" onClick={onRetry}>
          {t('actions.retry')}
        </Button>
      )}
    </div>
  )
}

interface EmptyStateProps {
  title: string
  hint?: string
  actionLabel?: string
  actionTo?: string
  icon?: ReactNode
}

export function EmptyState({
  title,
  hint,
  actionLabel,
  actionTo,
  icon,
}: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center gap-3 rounded-xl bg-card px-6 py-12 text-center ring-1 ring-foreground/10">
      {icon && <div className="text-muted-foreground">{icon}</div>}
      <p className="text-base font-medium">{title}</p>
      {hint && (
        <p className="max-w-prose text-sm text-muted-foreground">{hint}</p>
      )}
      {actionLabel && actionTo && (
        <Button render={<Link to={actionTo} />}>{actionLabel}</Button>
      )}
    </div>
  )
}
