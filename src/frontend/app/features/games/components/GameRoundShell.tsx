import { useTranslation } from 'react-i18next'
import { Check, Play, RotateCcw, Sparkles, X } from 'lucide-react'
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
    | 'state'
    | 'lastResult'
    | 'isStarting'
    | 'error'
    | 'start'
  >
  isLoading: boolean
  children: ReactNode
}

export function GameRoundShell({
  state,
  round,
  isLoading,
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
  const started = Boolean(round.round)
  const finished = round.state === 'won' || round.state === 'lost'
  const feedback = round.lastResult

  return (
    <div className="flex flex-col gap-4">
      {started && (
        <div className="flex items-center justify-between gap-3">
          <ScoreDots streak={round.streak} target={target} />
          <p className="shrink-0 text-sm font-semibold tabular-nums">
            {round.streak}
            <span className="text-muted-foreground"> / {target}</span>
          </p>
        </div>
      )}

      {feedback && !finished && (
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
          hint={t('games:round.wonHint', { count: round.streak })}
        />
      )}

      {round.state === 'lost' && (
        <GameOutcome
          tone="lose"
          title={t('games:round.lost')}
          hint={t('games:round.lostHint', { count: round.streak })}
        />
      )}

      {started && children}

      {Boolean(round.error) && (
        <p role="alert" className="text-sm text-destructive">
          {translateApiError(round.error, t)}
        </p>
      )}

      {!round.round && (
        <Button
          type="button"
          size="lg"
          disabled={round.isStarting}
          onClick={() => void round.start()}
        >
          <Play className="size-4" aria-hidden="true" />
          {t('games:play')}
        </Button>
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
}

function ScoreDots({ streak, target }: ScoreDotsProps) {
  if (target <= 0) return null

  return (
    <div aria-hidden="true" className="flex min-w-0 flex-1 gap-1.5">
      {Array.from({ length: target }, (_, index) => (
        <span
          key={index}
          className={cn(
            'h-1.5 flex-1 rounded-full transition-colors duration-300',
            index < streak ? 'bg-primary' : 'bg-muted',
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
