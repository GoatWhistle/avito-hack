import { useTranslation } from 'react-i18next'
import { Gamepad2, RotateCcw } from 'lucide-react'
import { Button } from '#/components/ui'
import { translateApiError } from '#/api'
import { RaccoonJumpResult } from './RaccoonJumpResult'
import type { RaccoonJumpReveal } from '#/features/games/types'
import type { GameStatus } from './game'

interface RaccoonJumpOverlaysProps {
  status: GameStatus
  finalScore: number
  bestScore: number
  reveal: RaccoonJumpReveal | null
  minStreakScore: number
  isStarting: boolean
  isSubmitting: boolean
  startError: unknown
  submitError: unknown
  onStart: () => void
  onRetrySubmit: () => void
}

export function RaccoonJumpOverlays({
  status,
  finalScore,
  bestScore,
  reveal,
  minStreakScore,
  isStarting,
  isSubmitting,
  startError,
  submitError,
  onStart,
  onRetrySubmit,
}: RaccoonJumpOverlaysProps) {
  const { t } = useTranslation(['games', 'errors'])

  if (status === 'idle') {
    return (
      <div className="absolute inset-0 flex flex-col items-center justify-center overflow-y-auto rounded-xl bg-background/70 p-6 backdrop-blur-sm">
        <div className="game-card-enter flex w-full max-w-[17rem] flex-col items-center gap-4 rounded-2xl bg-card p-6 text-center shadow-lg ring-1 ring-border">
          <span className="flex size-12 items-center justify-center rounded-full bg-primary-subtle text-primary-subtle-foreground">
            <Gamepad2 className="size-6" aria-hidden="true" />
          </span>

          <h2 className="text-xl font-semibold text-card-foreground">
            {t('games:raccoonjump.name')}
          </h2>

          {bestScore > 0 && (
            <p className="text-sm text-muted-foreground tabular-nums">
              {t('games:raccoonjump.best', { score: bestScore })}
            </p>
          )}

          {Boolean(startError) && (
            <p role="alert" className="text-sm text-destructive">
              {t('games:raccoonjump.startFailed')}:{' '}
              {translateApiError(startError, t)}
            </p>
          )}

          <Button
            type="button"
            size="lg"
            disabled={isStarting}
            onClick={onStart}
            className="w-full"
          >
            {t('games:raccoonjump.play')}
          </Button>
        </div>
      </div>
    )
  }

  if (status === 'gameover') {
    return (
      <div className="absolute inset-0 flex flex-col items-center justify-center overflow-y-auto rounded-xl bg-background/70 p-4 backdrop-blur-sm">
        <div className="game-card-enter flex w-full max-w-md flex-col items-center gap-4 rounded-2xl bg-card p-5 text-center shadow-lg ring-1 ring-border">
          <h2 className="text-xl font-semibold text-card-foreground">
            {t('games:raccoonjump.gameOver')}
          </h2>

          <div className="flex w-full flex-col gap-2">
            <p className="rounded-xl bg-primary-subtle py-3 text-2xl font-semibold text-primary-subtle-foreground tabular-nums">
              {t('games:raccoonjump.score', { score: finalScore })}
            </p>

            <RaccoonJumpResult
              reveal={reveal}
              minStreakScore={minStreakScore}
              isSubmitting={isSubmitting}
              submitError={submitError}
              onRetry={onRetrySubmit}
            />
          </div>

          {Boolean(startError) && (
            <p role="alert" className="text-sm text-destructive">
              {t('games:raccoonjump.startFailed')}:{' '}
              {translateApiError(startError, t)}
            </p>
          )}

          <Button
            type="button"
            size="lg"
            disabled={isStarting}
            onClick={onStart}
            className="w-full"
          >
            <RotateCcw className="size-4" aria-hidden="true" />
            {t('games:raccoonjump.playAgain')}
          </Button>
        </div>
      </div>
    )
  }

  return (
    <div
      className="pointer-events-none absolute bottom-4 flex w-full select-none justify-between px-6 text-4xl text-muted-foreground/30"
      aria-hidden
    >
      <div>‹</div>
      <div>›</div>
    </div>
  )
}
