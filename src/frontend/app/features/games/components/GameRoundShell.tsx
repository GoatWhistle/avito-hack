import { useTranslation } from 'react-i18next'
import { Check, RotateCcw, Sparkles, X } from 'lucide-react'
import { Button } from '#/components/ui'
import { translateApiError } from '#/api'
import { cn } from '#/lib/utils'
import type { GameRoundApi } from '#/features/games/hooks'
import type { GameState } from '#/features/games/types'
import type { ReactNode } from 'react'

interface GameRoundShellProps {
  state?: GameState
  round: Pick<
    GameRoundApi<unknown, never, unknown>,
    | 'round'
    | 'streak'
    | 'targetStreak'
    | 'attemptsUsed'
    | 'maxAttempts'
    | 'state'
    | 'lastResult'
    | 'isStarting'
    | 'error'
    | 'start'
  >
  isLoading: boolean
  hideError?: boolean
  showGuessFeedback?: boolean
  children: ReactNode
}

export function GameRoundShell({
  state,
  round,
  isLoading,
  hideError = false,
  showGuessFeedback = true,
  children,
}: GameRoundShellProps) {
  const { t } = useTranslation(['games', 'errors'])

  if (isLoading) {
    return (
      <p className="py-16 text-center text-sm text-muted-foreground">
        {t('games:loading')}
      </p>
    )
  }

  const target = round.targetStreak || state?.target_streak || 0
  const attemptBased = target === 0 && round.maxAttempts > 0
  const filled = attemptBased ? round.attemptsUsed : round.streak
  const total = attemptBased ? round.maxAttempts : target
  const started = Boolean(round.round)
  const finished = round.state === 'won' || round.state === 'lost'
  const feedback = round.lastResult

  return (
    <div className="flex flex-col gap-4">
      {started && total > 0 && (
        <div className="flex items-center justify-between gap-3">
          <ScoreDots
            streak={filled}
            target={total}
            tone={attemptBased ? 'attempts' : 'streak'}
          />
          <p className="shrink-0 text-sm font-semibold tabular-nums">
            {filled}
            <span className="text-muted-foreground"> / {total}</span>
          </p>
        </div>
      )}

      {feedback && showGuessFeedback && !finished && (
        <p
          role="status"
          className={cn(
            'game-card-enter flex items-center justify-center gap-1.5 rounded-xl py-2.5 text-sm font-semibold',
            feedback.correct
              ? 'bg-primary-subtle text-primary-subtle-foreground'
              : 'bg-destructive-subtle text-destructive',
          )}
        >
          {feedback.correct ? (
            <Check className="size-4" aria-hidden="true" />
          ) : (
            <X className="size-4" aria-hidden="true" />
          )}
          {feedback.correct ? t('games:round.correct') : t('games:round.wrong')}
        </p>
      )}

      {round.state === 'won' && (
        <GameOutcome
          tone="win"
          title={t('games:round.won')}
          hint={
            attemptBased
              ? t('games:round.wonHintAttempts', {
                  count: round.attemptsUsed,
                  total: round.maxAttempts,
                })
              : t('games:round.wonHint', { count: round.streak })
          }
        />
      )}

      {round.state === 'lost' && (
        <GameOutcome
          tone="lose"
          title={t('games:round.lost')}
          hint={
            attemptBased
              ? t('games:round.lostHintAttempts', { total: round.maxAttempts })
              : t('games:round.lostHint', { count: round.streak })
          }
        />
      )}

      {started && children}

      {Boolean(round.error) && !hideError && (
        <p role="alert" className="text-sm text-destructive">
          {translateApiError(round.error, t)}
        </p>
      )}

      {!round.round && (
        <p className="py-16 text-center text-sm text-muted-foreground">
          {t('games:loading')}
        </p>
      )}

      {finished && (
        <Button
          type="button"
          size="lg"
          disabled={round.isStarting}
          onClick={() => void round.start()}
        >
          <RotateCcw className="size-4" aria-hidden="true" />
          {t('games:playAgain')}
        </Button>
      )}
    </div>
  )
}

interface ScoreDotsProps {
  streak: number
  target: number
  tone?: 'streak' | 'attempts'
}

function ScoreDots({ streak, target, tone = 'streak' }: ScoreDotsProps) {
  if (target <= 0) return null

  return (
    <div aria-hidden="true" className="flex min-w-0 flex-1 gap-1.5">
      {Array.from({ length: target }, (_, index) => (
        <span
          key={index}
          className={cn(
            'h-1.5 flex-1 rounded-full transition-colors duration-300',
            index < streak
              ? tone === 'attempts'
                ? 'bg-muted-foreground/50'
                : 'bg-primary'
              : tone === 'attempts'
                ? 'bg-primary/70'
                : 'bg-muted',
          )}
        />
      ))}
    </div>
  )
}

interface GameOutcomeProps {
  tone: 'win' | 'lose'
  title: string
  hint: string
}

function GameOutcome({ tone, title, hint }: GameOutcomeProps) {
  const win = tone === 'win'

  return (
    <div
      role="status"
      className={cn(
        'game-card-enter relative flex items-center gap-4 overflow-hidden rounded-2xl px-5 py-4',
        win
          ? 'bg-primary-subtle ring-1 ring-primary/30'
          : 'bg-destructive-subtle ring-1 ring-destructive/25',
      )}
    >
      <span
        className={cn(
          'flex size-11 shrink-0 items-center justify-center rounded-full',
          win
            ? 'bg-primary text-primary-foreground'
            : 'bg-destructive/15 text-destructive',
        )}
      >
        {win ? (
          <Sparkles className="size-5" aria-hidden="true" />
        ) : (
          <X className="size-5" aria-hidden="true" />
        )}
      </span>

      <div className="flex min-w-0 flex-col">
        <p
          className={cn(
            'text-base font-semibold',
            win ? 'text-primary-subtle-foreground' : 'text-destructive',
          )}
        >
          {title}
        </p>
        <p className="text-sm text-muted-foreground">{hint}</p>
      </div>
    </div>
  )
}
