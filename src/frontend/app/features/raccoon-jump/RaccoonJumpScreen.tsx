import type {
  TouchEvent as ReactTouchEvent,
  TouchList as ReactTouchList,
} from 'react'
import { useCallback, useEffect, useRef, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useGameStateQuery } from '#/features/games/hooks'
import { RaccoonJumpOverlays } from './RaccoonJumpOverlays'
import { RaccoonJumpPlayer } from './RaccoonJumpPlayer'
import { GAME_HEIGHT, GAME_WIDTH, type Game, type GameStatus } from './game'
import { useGameKeyboard } from './useGameKeyboard'
import { useGameLoop } from './useGameLoop'
import { usePalette } from './usePalette'
import { RACCOON_JUMP_SLUG, useRaccoonJumpRound } from './useRaccoonJumpRound'

const MAX_SCALE = 2
const VERTICAL_MARGIN = 16

export function RaccoonJumpScreen() {
  const { t } = useTranslation('games')
  const canvasRef = useRef<HTMLCanvasElement>(null)
  const gameRef = useRef<Game | null>(null)
  const lottieContainerRef = useRef<HTMLDivElement>(null)
  const frameRef = useRef<HTMLDivElement>(null)
  const surfaceRef = useRef<HTMLDivElement>(null)
  const disposedRef = useRef(false)
  const paletteRef = usePalette()

  const [status, setStatus] = useState<GameStatus>('idle')
  const [finalScore, setFinalScore] = useState(0)
  const [scale, setScale] = useState(1)
  const scaleRef = useRef(1)

  const round = useRaccoonJumpRound()
  const stateQuery = useGameStateQuery(RACCOON_JUMP_SLUG)
  const { adoptState, submitScore } = round
  const serverBest = stateQuery.data?.best_score

  useEffect(() => {
    if (typeof serverBest === 'number') adoptState(serverBest)
  }, [serverBest, adoptState])

  const onGameOver = useCallback(
    (score: number, collected: number[]) => {
      setFinalScore(score)
      void submitScore({ score, collected })
    },
    [submitScore],
  )

  const onGameOverRef = useRef(onGameOver)
  useEffect(() => {
    onGameOverRef.current = onGameOver
  }, [onGameOver])

  useEffect(() => {
    const frame = frameRef.current
    if (!frame) return

    const measure = () => {
      const availableWidth = frame.clientWidth
      if (!availableWidth) return

      const viewportHeight = window.innerHeight
      const offsetTop = frame.getBoundingClientRect().top
      const availableHeight = viewportHeight - offsetTop - VERTICAL_MARGIN

      const widthScale = availableWidth / GAME_WIDTH
      const heightScale =
        availableHeight > 0 ? availableHeight / GAME_HEIGHT : widthScale

      const next = Math.min(MAX_SCALE, widthScale, heightScale)
      scaleRef.current = next
      setScale(next)
    }

    measure()

    const observer = new ResizeObserver(measure)
    observer.observe(frame)
    window.addEventListener('resize', measure)
    return () => {
      observer.disconnect()
      window.removeEventListener('resize', measure)
    }
  }, [])

  useGameLoop({
    canvasRef,
    gameRef,
    lottieContainerRef,
    disposedRef,
    paletteRef,
    scaleRef,
    onStatusChange: setStatus,
    onGameOverRef,
  })

  const { isStarting, startRound } = round

  const startGame = useCallback(async () => {
    const game = gameRef.current
    if (!game || isStarting) return

    const prompt = await startRound()
    if (!prompt || disposedRef.current) return

    game.collectibleIndexes = (prompt.collectibles ?? []).map(
      (item) => item.index,
    )
    game.reset(prompt.seed)
  }, [isStarting, startRound])

  const requestStart = useCallback(() => {
    void startGame()
  }, [startGame])

  useGameKeyboard(gameRef, surfaceRef, requestStart)

  const updateTouchInput = useCallback(
    (touches: ReactTouchList, rect: DOMRect) => {
      const game = gameRef.current
      if (!game) return

      game.input.left = false
      game.input.right = false

      const mid = rect.width / 2
      for (let i = 0; i < touches.length; i++) {
        const touch = touches.item(i)
        if (!touch) continue
        const x = touch.clientX - rect.left
        if (x < mid) game.input.left = true
        else game.input.right = true
      }
    },
    [],
  )

  const onTouchStart = useCallback(
    (e: ReactTouchEvent<HTMLDivElement>) => {
      updateTouchInput(e.touches, e.currentTarget.getBoundingClientRect())
    },
    [updateTouchInput],
  )

  const onTouchEnd = useCallback(
    (e: ReactTouchEvent<HTMLDivElement>) => {
      updateTouchInput(e.touches, e.currentTarget.getBoundingClientRect())
    },
    [updateTouchInput],
  )

  return (
    <div
      ref={frameRef}
      className="flex w-full items-start justify-center"
      style={{ height: GAME_HEIGHT * scale }}
    >
      <div
        className="relative"
        style={{ width: GAME_WIDTH * scale, height: GAME_HEIGHT * scale }}
        onTouchStart={onTouchStart}
        onTouchMove={onTouchStart}
        onTouchEnd={onTouchEnd}
        onTouchCancel={onTouchEnd}
      >
        <div
          ref={surfaceRef}
          className="absolute left-0 top-0 origin-top-left"
          style={{
            width: GAME_WIDTH,
            height: GAME_HEIGHT,
            transform: `scale(${scale})`,
          }}
        >
          <canvas
            ref={canvasRef}
            role="img"
            aria-label={t('raccoonjump.name')}
            className="absolute left-0 top-0 h-full w-full rounded-xl ring-1 ring-border"
          >
            {t('raccoonjump.description')}
          </canvas>

          <RaccoonJumpPlayer containerRef={lottieContainerRef} status={status} />
        </div>

        <p role="status" aria-live="polite" className="sr-only">
          {status === 'gameover' ? t('raccoonjump.score', { score: finalScore }) : ''}
        </p>

        <RaccoonJumpOverlays
          status={status}
          finalScore={finalScore}
          bestScore={round.bestScore}
          reveal={round.reveal}
          minStreakScore={round.minStreakScore}
          isStarting={round.isStarting}
          isSubmitting={round.isSubmitting}
          startError={round.startError}
          submitError={round.submitError}
          onStart={requestStart}
          onRetrySubmit={() => void round.retrySubmit()}
        />
      </div>
    </div>
  )
}
