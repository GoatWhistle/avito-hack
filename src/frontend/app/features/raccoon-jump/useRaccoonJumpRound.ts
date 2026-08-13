import { useCallback, useRef, useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { gameRepository } from '#/features/games/repository'
import { gameKeys } from '#/features/games/hooks'
import type {
  GameRound,
  RaccoonJumpMove,
  RaccoonJumpPrompt,
  RaccoonJumpReveal,
} from '#/features/games/types'

export const RACCOON_JUMP_SLUG = 'raccoonjump'

export interface RaccoonJumpRoundApi {
  bestScore: number
  collectibles: RaccoonJumpPrompt['collectibles']
  reveal: RaccoonJumpReveal | null
  minStreakScore: number
  isStarting: boolean
  isSubmitting: boolean
  startError: unknown
  submitError: unknown
  startRound: () => Promise<RaccoonJumpPrompt | null>
  submitScore: (move: RaccoonJumpMove) => Promise<void>
  retrySubmit: () => Promise<void>
  adoptState: (best: number) => void
}

const DEFAULT_MIN_STREAK_SCORE = 25

export const useRaccoonJumpRound = (): RaccoonJumpRoundApi => {
  const queryClient = useQueryClient()

  const [bestScore, setBestScore] = useState(0)
  const [collectibles, setCollectibles] =
    useState<RaccoonJumpPrompt['collectibles']>(undefined)
  const [minStreakScore, setMinStreakScore] = useState(DEFAULT_MIN_STREAK_SCORE)
  const [reveal, setReveal] = useState<RaccoonJumpReveal | null>(null)
  const [isStarting, setIsStarting] = useState(false)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [startError, setStartError] = useState<unknown>(null)
  const [submitError, setSubmitError] = useState<unknown>(null)

  const roundIdRef = useRef<string | null>(null)
  const pendingMoveRef = useRef<RaccoonJumpMove | null>(null)

  const adoptState = useCallback((best: number) => {
    setBestScore((current) => Math.max(current, best))
  }, [])

  const startRound = useCallback(async () => {
    setIsStarting(true)
    setStartError(null)
    setSubmitError(null)
    setReveal(null)
    pendingMoveRef.current = null

    try {
      const started: GameRound = await gameRepository.startRound(
        RACCOON_JUMP_SLUG,
      )
      const prompt = started.prompt as RaccoonJumpPrompt

      roundIdRef.current = started.round_id
      setCollectibles(prompt?.collectibles)
      setBestScore((current) => Math.max(current, prompt?.best_score ?? 0))
      if (prompt?.min_streak_score) setMinStreakScore(prompt.min_streak_score)

      return prompt ?? null
    } catch (cause) {
      roundIdRef.current = null
      setStartError(cause)

      return null
    } finally {
      setIsStarting(false)
    }
  }, [])

  const send = useCallback(
    async (move: RaccoonJumpMove) => {
      const roundId = roundIdRef.current
      if (!roundId) return

      setIsSubmitting(true)
      setSubmitError(null)

      try {
        const result = await gameRepository.guess(
          RACCOON_JUMP_SLUG,
          roundId,
          move,
        )
        const revealed = result.reveal as RaccoonJumpReveal

        pendingMoveRef.current = null
        roundIdRef.current = null
        setReveal(revealed)
        setBestScore(revealed?.best_score ?? move.score)

        void queryClient.invalidateQueries({
          queryKey: gameKeys.state(RACCOON_JUMP_SLUG),
        })
        void queryClient.invalidateQueries({ queryKey: gameKeys.list() })
      } catch (cause) {
        setSubmitError(cause)
      } finally {
        setIsSubmitting(false)
      }
    },
    [queryClient],
  )

  const submitScore = useCallback(
    async (move: RaccoonJumpMove) => {
      pendingMoveRef.current = move
      await send(move)
    },
    [send],
  )

  const retrySubmit = useCallback(async () => {
    const pending = pendingMoveRef.current
    if (!pending) return

    await send(pending)
  }, [send])

  return {
    bestScore,
    collectibles,
    reveal,
    minStreakScore,
    isStarting,
    isSubmitting,
    startError,
    submitError,
    startRound,
    submitScore,
    retrySubmit,
    adoptState,
  }
}
