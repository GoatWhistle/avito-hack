import { useTranslation } from 'react-i18next'
import { Gamepad2, RotateCcw, Trophy } from 'lucide-react'
import { Button } from '#/components/ui'
import type { GameStatus } from './game'

interface DoodleJumpOverlaysProps {
  status: GameStatus
  finalScore: number
  bestScore: number
  onStart: () => void
}

export function DoodleJumpOverlays({
  status,
  finalScore,
  bestScore,
  onStart,
}: DoodleJumpOverlaysProps) {
  const { t } = useTranslation('games')

  if (status === 'idle') {
    return (
      <div className="absolute inset-0 flex flex-col items-center justify-center rounded-xl bg-background/70 p-6 backdrop-blur-sm">
        <div className="game-card-enter flex w-full max-w-[17rem] flex-col items-center gap-4 rounded-2xl bg-card p-6 text-center shadow-lg ring-1 ring-border">
          <span className="flex size-12 items-center justify-center rounded-full bg-primary-subtle text-primary-subtle-foreground">
            <Gamepad2 className="size-6" aria-hidden="true" />
          </span>

          <h2 className="text-xl font-semibold text-card-foreground">
            {t('raccoonjump.name')}
          </h2>

          <Button type="button" size="lg" onClick={onStart} className="w-full">
            {t('raccoonjump.play')}
          </Button>
        </div>
      </div>
    )
  }

  if (status === 'gameover') {
    return (
      <div className="absolute inset-0 flex flex-col items-center justify-center rounded-xl bg-background/70 p-6 backdrop-blur-sm">
        <div className="game-card-enter flex w-full max-w-[17rem] flex-col items-center gap-4 rounded-2xl bg-card p-6 text-center shadow-lg ring-1 ring-border">
          <h2 className="text-xl font-semibold text-card-foreground">
            {t('raccoonjump.gameOver')}
          </h2>

          <div className="flex w-full flex-col gap-2">
            <p className="rounded-xl bg-primary-subtle py-3 text-2xl font-semibold text-primary-subtle-foreground tabular-nums">
              {t('raccoonjump.score', { score: finalScore })}
            </p>

            <p className="flex items-center justify-center gap-1.5 text-sm text-muted-foreground tabular-nums">
              <Trophy className="size-4" aria-hidden="true" />
              {t('raccoonjump.best', { score: bestScore })}
            </p>
          </div>

          <Button type="button" size="lg" onClick={onStart} className="w-full">
            <RotateCcw className="size-4" aria-hidden="true" />
            {t('raccoonjump.playAgain')}
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
