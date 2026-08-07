import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { Button } from '#/components/ui'
import { cn } from '#/lib/utils'

export function RowsSkeleton({
  count = 6,
  className,
}: {
  count?: number
  className?: string
}) {
  return (
    <div className={cn('flex flex-col gap-2', className)} aria-hidden="true">
      {Array.from({ length: count }).map((_, index) => (
        <div key={index} className="h-10 animate-pulse rounded-lg bg-muted" />
      ))}
    </div>
  )
}

export function LeaderboardLoading({ label }: { label: string }) {
  return (
    <div role="status" aria-live="polite" aria-busy="true">
      <span className="sr-only">{label}</span>
      <RowsSkeleton />
    </div>
  )
}

export function LeaderboardError({
  message,
  onRetry,
}: {
  message: string
  onRetry?: () => void
}) {
  const { t } = useTranslation('leaderboard')

  return (
    <div
      role="alert"
      className="flex flex-col items-start gap-3 rounded-xl bg-destructive-subtle p-4 text-destructive-subtle-foreground"
    >
      <p className="text-sm">{message}</p>
      {onRetry && (
        <Button variant="outline" size="sm" onClick={onRetry}>
          {t('actions.retry')}
        </Button>
      )}
    </div>
  )
}

export function LeaderboardGuest() {
  const { t } = useTranslation('leaderboard')

  return (
    <section className="flex flex-col items-center gap-3 rounded-xl bg-primary-subtle px-4 py-8 text-center">
      <p className="text-sm font-semibold text-primary-subtle-foreground">
        {t('guest.title')}
      </p>
      <p className="max-w-sm text-sm text-primary-subtle-foreground/90">
        {t('guest.hint')}
      </p>
      <div className="flex flex-wrap items-center justify-center gap-2">
        <Button render={<Link to="/sign-up" />} className="no-underline">
          {t('guest.signUp')}
        </Button>
        <Button
          variant="outline"
          render={<Link to="/sign-in" />}
          className="no-underline"
        >
          {t('guest.signIn')}
        </Button>
      </div>
    </section>
  )
}

export function LeaderboardEmpty({
  title,
  hint,
}: {
  title: string
  hint?: string
}) {
  return (
    <div className="flex flex-col items-center gap-2 rounded-xl bg-muted/50 px-4 py-8 text-center">
      <p className="text-sm font-medium text-foreground">{title}</p>
      {hint && <p className="text-sm text-muted-foreground">{hint}</p>}
    </div>
  )
}
