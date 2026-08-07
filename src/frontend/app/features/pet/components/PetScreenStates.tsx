import { useTranslation } from 'react-i18next'
import { Button } from '#/components/ui'

export function PetScreenSkeleton() {
  const { t } = useTranslation('pet')

  return (
    <div
      role="status"
      aria-busy="true"
      aria-label={t('states.loading')}
      className="flex w-full flex-col items-center gap-6 py-6"
    >
      <div className="h-6 w-32 animate-pulse rounded-full bg-muted" />
      <div className="size-56 animate-pulse rounded-full bg-muted" />
      <div className="flex w-full flex-col gap-3">
        {[0, 1, 2].map((row) => (
          <div key={row} className="flex flex-col gap-1.5">
            <div className="h-3 w-24 animate-pulse rounded-full bg-muted" />
            <div className="h-2 w-full animate-pulse rounded-full bg-muted" />
          </div>
        ))}
      </div>
      <div className="grid w-full grid-cols-2 gap-2">
        <div className="h-9 animate-pulse rounded-lg bg-muted" />
        <div className="h-9 animate-pulse rounded-lg bg-muted" />
      </div>
    </div>
  )
}

export interface PetScreenErrorProps {
  message: string
  onRetry: () => void
  isRetrying: boolean
}

export function PetScreenError({
  message,
  onRetry,
  isRetrying,
}: PetScreenErrorProps) {
  const { t } = useTranslation('pet')

  return (
    <div
      role="alert"
      className="flex flex-col items-center gap-4 rounded-xl bg-destructive-subtle px-4 py-8 text-center"
    >
      <span aria-hidden="true" className="text-3xl">
        😿
      </span>
      <div className="flex flex-col gap-1">
        <h2 className="text-base font-semibold text-destructive-subtle-foreground">
          {t('states.errorTitle')}
        </h2>
        <p className="text-sm text-destructive-subtle-foreground/80">
          {message}
        </p>
      </div>
      <Button variant="outline" onClick={onRetry} disabled={isRetrying}>
        {isRetrying ? t('states.retrying') : t('states.retry')}
      </Button>
    </div>
  )
}

export interface PetScreenEmptyProps {
  onCheckIn: () => void
  isCheckingIn: boolean
}

export function PetScreenEmpty({
  onCheckIn,
  isCheckingIn,
}: PetScreenEmptyProps) {
  const { t } = useTranslation('pet')

  return (
    <div className="flex flex-col items-center gap-5 py-10 text-center">
      <span aria-hidden="true" className="animate-pulse text-6xl">
        🥚
      </span>
      <div className="flex flex-col gap-1.5">
        <h2 className="text-xl font-semibold text-foreground">
          {t('hatching.title')}
        </h2>
        <p className="max-w-xs text-sm text-muted-foreground">
          {t('hatching.hint')}
        </p>
      </div>
      <Button size="lg" onClick={onCheckIn} disabled={isCheckingIn}>
        {isCheckingIn ? t('actions.checkingIn') : t('actions.checkIn')}
      </Button>
    </div>
  )
}
