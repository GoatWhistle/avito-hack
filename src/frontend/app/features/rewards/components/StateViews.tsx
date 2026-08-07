import { useTranslation } from 'react-i18next'
import { Button } from '#/components/ui'
import { cn } from '#/lib/utils'

export function SkeletonList({
  count = 3,
  className,
}: {
  count?: number
  className?: string
}) {
  return (
    <div className={cn('flex flex-col gap-3', className)} aria-hidden="true">
      {Array.from({ length: count }).map((_, index) => (
        <div key={index} className="h-24 animate-pulse rounded-xl bg-muted" />
      ))}
    </div>
  )
}

export function LoadingState({ label }: { label: string }) {
  return (
    <div role="status" aria-live="polite" aria-busy="true">
      <span className="sr-only">{label}</span>
      <SkeletonList />
    </div>
  )
}

export function ErrorState({
  message,
  onRetry,
}: {
  message: string
  onRetry?: () => void
}) {
  const { t } = useTranslation('common')

  return (
    <div
      role="alert"
      className="flex flex-col items-start gap-3 rounded-xl bg-destructive-subtle p-4 text-destructive-subtle-foreground"
    >
      <p className="text-sm">{message}</p>
      {onRetry && (
        <Button variant="outline" size="sm" onClick={onRetry}>
          {t('actions.retry', { defaultValue: 'Повторить' })}
        </Button>
      )}
    </div>
  )
}

export function EmptyState({
  title,
  hint,
  action,
}: {
  title: string
  hint?: string
  action?: React.ReactNode
}) {
  return (
    <div className="flex flex-col items-center gap-2 rounded-xl bg-muted/50 px-4 py-8 text-center">
      <p className="text-sm font-medium text-foreground">{title}</p>
      {hint && <p className="text-sm text-muted-foreground">{hint}</p>}
      {action}
    </div>
  )
}
