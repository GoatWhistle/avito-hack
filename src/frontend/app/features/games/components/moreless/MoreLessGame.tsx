import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { ArrowRight, ChevronDown, ChevronUp } from 'lucide-react'
import { Button } from '#/components/ui'
import { useGameRound, useGameStateQuery } from '#/features/games/hooks'
import { GameRoundShell } from '../GameRoundShell'
import { MoreLessCard } from './MoreLessCard'
import type {
  MoreLessChoice,
  MoreLessPrompt,
  MoreLessReveal,
} from '#/features/games/types'

const SLUG = 'moreless'

export function MoreLessGame() {
  const { t } = useTranslation('games')
  const stateQuery = useGameStateQuery(SLUG)
  const active = stateQuery.data?.active_round
  const round = useGameRound<
    MoreLessPrompt,
    { choice: MoreLessChoice },
    MoreLessReveal
  >(
    SLUG,
    active && {
      round_id: active.round_id,
      streak: active.streak,
      target_streak: stateQuery.data?.target_streak ?? 0,
      state: 'active',
      prompt: active.prompt,
    },
  )
  const [frozen, setFrozen] = useState<MoreLessPrompt | null>(null)

  const { lastResult, dismissFeedback, prompt, guess } = round
  const view = frozen ?? prompt
  const revealed = lastResult?.reveal.right_price
  const settled = Boolean(lastResult)

  const submit = async (choice: MoreLessChoice) => {
    setFrozen(prompt)
    await guess({ choice })
  }

  const advance = () => {
    setFrozen(null)
    dismissFeedback()
  }

  return (
    <GameRoundShell
      state={stateQuery.data}
      round={round}
      isLoading={stateQuery.isPending}
    >
      {view && (
        <div className="flex flex-col gap-4">
          <div className="relative grid gap-3 sm:grid-cols-2 sm:gap-4">
            <MoreLessCard
              key={`left-${view.left.item_id}`}
              item={view.left}
              linkable={settled}
            />

            <span
              aria-hidden="true"
              className="pointer-events-none absolute top-1/2 left-1/2 z-10 hidden size-16 -translate-x-1/2 -translate-y-1/2 items-center justify-center rounded-full bg-background text-[0.6875rem] font-bold tracking-wide text-muted-foreground uppercase shadow-lg ring-1 ring-border sm:flex"
            >
              {t('moreless.vs')}
            </span>

            <MoreLessCard
              key={`right-${view.right.item_id}`}
              item={view.right}
              revealedPrice={revealed}
              hint={t('moreless.thisItem')}
              highlight={
                lastResult ? (lastResult.correct ? 'correct' : 'wrong') : null
              }
              linkable={settled}
            />
          </div>

          {settled ? (
            round.state === 'active' && (
              <Button type="button" size="lg" onClick={advance}>
                {t('round.next')}
                <ArrowRight className="size-4" aria-hidden="true" />
              </Button>
            )
          ) : (
            <div className="flex flex-col gap-2">
              <p className="text-center text-xs text-muted-foreground">
                {t('moreless.rules')}
              </p>
              <div className="grid grid-cols-2 gap-3">
                <Button
                  type="button"
                  size="lg"
                  variant="outline"
                  disabled={round.isGuessing}
                  onClick={() => void submit('higher')}
                >
                  <ChevronUp className="size-4" aria-hidden="true" />
                  {t('moreless.higher')}
                </Button>
                <Button
                  type="button"
                  size="lg"
                  variant="outline"
                  disabled={round.isGuessing}
                  onClick={() => void submit('lower')}
                >
                  <ChevronDown className="size-4" aria-hidden="true" />
                  {t('moreless.lower')}
                </Button>
              </div>
            </div>
          )}
        </div>
      )}
    </GameRoundShell>
  )
}
