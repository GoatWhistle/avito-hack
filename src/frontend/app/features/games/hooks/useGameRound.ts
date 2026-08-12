import { useCallback, useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { gameRepository } from '#/features/games/repository'
import { gameKeys } from './query-keys'
import type {
  GameGuessResult,
  GameRound,
  GameRoundState,
} from '#/features/games/types'

export interface GameRoundApi<TPrompt, TMove, TReveal> {
  round: GameRound | null
  prompt: TPrompt | null
  streak: number
  targetStreak: number
  state: GameRoundState | null
  lastResult: GuessFeedback<TReveal> | null
  isStarting: boolean
  isGuessing: boolean
  error: unknown
  start: () => Promise<void>
  guess: (move: TMove) => Promise<void>
  dismissFeedback: () => void
}

export interface GuessFeedback<TReveal> {
  correct: boolean
  reveal: TReveal
  attemptCompleted: boolean
}

export const useGameRound = <TPrompt, TMove, TReveal>(
  slug: string,
  initialRound?: GameRound,
): GameRoundApi<TPrompt, TMove, TReveal> => {
  const queryClient = useQueryClient()
  const [round, setRound] = useState<GameRound | null>(initialRound ?? null)
  const [adopted, setAdopted] = useState<string | null>(
    initialRound?.round_id ?? null,
  )
  const [lastResult, setLastResult] = useState<GuessFeedback<TReveal> | null>(
    null,
  )
  const [isStarting, setIsStarting] = useState(false)
  const [isGuessing, setIsGuessing] = useState(false)
  const [error, setError] = useState<unknown>(null)

  if (initialRound && initialRound.round_id !== adopted && !round) {
    setAdopted(initialRound.round_id)
    setRound(initialRound)
  }

  const start = useCallback(async () => {
    setIsStarting(true)
    setError(null)
    setLastResult(null)

    try {
      const started = await gameRepository.startRound(slug)
      setAdopted(started.round_id)
      setRound({ ...started, state: started.state ?? 'active' })
    } catch (cause) {
      setError(cause)
    } finally {
      setIsStarting(false)
    }
  }, [slug])

  const guess = useCallback(
    async (move: TMove) => {
      if (!round || isGuessing) return

      setIsGuessing(true)
      setError(null)

      try {
        const result: GameGuessResult = await gameRepository.guess(
          slug,
          round.round_id,
          move,
        )

        setLastResult({
          correct: result.correct,
          reveal: result.reveal as TReveal,
          attemptCompleted: result.attempt_completed,
        })

        setRound({
          round_id: round.round_id,
          streak: result.streak,
          target_streak: round.target_streak,
          state: result.state,
          prompt: result.prompt ?? round.prompt,
        })

        if (result.state !== 'active') {
          void queryClient.invalidateQueries({ queryKey: gameKeys.state(slug) })
          void queryClient.invalidateQueries({ queryKey: gameKeys.list() })
        }
      } catch (cause) {
        setError(cause)
      } finally {
        setIsGuessing(false)
      }
    },
    [round, isGuessing, slug, queryClient],
  )

  const dismissFeedback = useCallback(() => setLastResult(null), [])

  return {
    round,
    prompt: (round?.prompt ?? null) as TPrompt | null,
    streak: round?.streak ?? 0,
    targetStreak: round?.target_streak ?? 0,
    state: round?.state ?? null,
    lastResult,
    isStarting,
    isGuessing,
    error,
    start,
    guess,
    dismissFeedback,
  }
}
