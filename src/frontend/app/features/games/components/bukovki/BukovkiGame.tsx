import { useCallback, useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useGameRound, useGameStateQuery } from '#/features/games/hooks'
import { GameRoundShell } from '../GameRoundShell'
import { BukovkiGrid } from './BukovkiGrid'
import { BukovkiKeyboard } from './BukovkiKeyboard'
import { BukovkiListings } from './BukovkiListings'
import type { BukovkiPrompt, BukovkiReveal } from '#/features/games/types'

const SLUG = 'bukovki'
const LETTER = /^[а-яё]$/i

const guessErrorKey = (error: unknown) => {
  const message =
    typeof error === 'object' && error !== null
      ? ((error as { response?: { data?: { error?: { message?: string } } } })
          .response?.data?.error?.message ?? '')
      : ''

  if (message.includes('already guessed')) return 'bukovki.alreadyTried'
  if (message.includes('not valid')) return 'bukovki.notAWord'
  if (message.includes('word length')) return 'bukovki.tooShort'

  return null
}

export function BukovkiGame() {
  const { t } = useTranslation('games')
  const stateQuery = useGameStateQuery(SLUG)
  const active = stateQuery.data?.active_round
  const round = useGameRound<BukovkiPrompt, { guess: string }, BukovkiReveal>(
    SLUG,
    active && {
      round_id: active.round_id,
      streak: active.streak,
      target_streak: stateQuery.data?.target_streak ?? 0,
      state: 'active',
      prompt: active.prompt,
      attempts_used: active.attempts_used,
      max_attempts: active.max_attempts ?? stateQuery.data?.max_attempts,
    },
    stateQuery.isSuccess,
  )
  const [draft, setDraft] = useState('')

  const { prompt, lastResult, guess, isGuessing, state } = round
  const wordLength = prompt?.word_length ?? 0
  const playing = state === 'active'
  const ready = draft.length === wordLength && !isGuessing
  const guessError = guessErrorKey(round.error)

  const submit = useCallback(async () => {
    if (draft.length !== wordLength || isGuessing) return
    await guess({ guess: draft })
  }, [draft, wordLength, isGuessing, guess])

  const append = useCallback(
    (char: string) => {
      setDraft((current) =>
        current.length >= wordLength
          ? current
          : current + char.toLowerCase().replace('ё', 'е'),
      )
    },
    [wordLength],
  )

  const backspace = useCallback(() => {
    setDraft((current) => [...current].slice(0, -1).join(''))
  }, [])

  useEffect(() => {
    if (prompt) setDraft('')
  }, [round.round?.round_id, prompt?.history.length])

  useEffect(() => {
    if (!playing) return

    const onKeyDown = (event: KeyboardEvent) => {
      if (event.metaKey || event.ctrlKey || event.altKey) return

      if (event.key === 'Enter') {
        event.preventDefault()
        void submit()
        return
      }
      if (event.key === 'Backspace') {
        event.preventDefault()
        backspace()
        return
      }
      if (LETTER.test(event.key)) {
        event.preventDefault()
        append(event.key)
      }
    }

    window.addEventListener('keydown', onKeyDown)

    return () => window.removeEventListener('keydown', onKeyDown)
  }, [playing, submit, backspace, append])

  return (
    <GameRoundShell
      state={stateQuery.data}
      round={round}
      isLoading={stateQuery.isPending}
      hideError={Boolean(guessError)}
      showGuessFeedback={false}
    >
      {prompt && (
        <div className="relative flex flex-col gap-6">
          <div className="flex min-w-0 flex-col gap-6 sm:gap-8">
            <BukovkiGrid
              wordLength={prompt.word_length}
              maxTries={prompt.max_tries}
              history={prompt.history}
              draft={draft}
              active={playing}
            />

            {guessError && (
              <p
                role="alert"
                className="text-center text-sm font-medium text-destructive"
              >
                {t(guessError as 'bukovki.notAWord')}
              </p>
            )}

            {lastResult?.reveal.secret && (
              <p className="text-center text-sm font-medium">
                {t('bukovki.secretWas', {
                  word: lastResult.reveal.secret.toUpperCase(),
                })}
              </p>
            )}

            {playing && (
              <div className="mt-2 sm:mt-4">
                <BukovkiKeyboard
                  history={prompt.history}
                  disabled={isGuessing}
                  canSubmit={ready}
                  onKey={append}
                  onBackspace={backspace}
                  onSubmit={() => void submit()}
                  submitLabel={t('bukovki.submit')}
                  backspaceLabel={t('bukovki.backspace')}
                />
              </div>
            )}

            {lastResult?.reveal.game_over && (
              <BukovkiListings listings={lastResult.reveal.listings ?? []} />
            )}
          </div>
        </div>
      )}
    </GameRoundShell>
  )
}
