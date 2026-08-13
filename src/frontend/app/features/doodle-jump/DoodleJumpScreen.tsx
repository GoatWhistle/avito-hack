import type {
  TouchEvent as ReactTouchEvent,
  TouchList as ReactTouchList,
} from 'react'
import { useCallback, useEffect, useRef, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { DoodleJumpOverlays } from './DoodleJumpOverlays'
import { DoodleJumpPlayer } from './DoodleJumpPlayer'
import { Game, GAME_HEIGHT, GAME_WIDTH, type GameStatus } from './game'
import { GameRenderer } from './render'
import { useGameKeyboard } from './useGameKeyboard'
import { usePalette } from './usePalette'

const LS_KEY = 'doodle-jump-best'

export function DoodleJumpScreen() {
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
  const [bestScore, setBestScore] = useState(0)
  const [scale, setScale] = useState(1)

  useEffect(() => {
    const frame = frameRef.current
    if (!frame) return

    const measure = () => {
      const available = frame.clientWidth
      if (!available) return
      setScale(Math.min(1, available / GAME_WIDTH))
    }

    measure()

    const observer = new ResizeObserver(measure)
    observer.observe(frame)
    return () => observer.disconnect()
  }, [])

  useEffect(() => {
    const canvas = canvasRef.current
    if (!canvas) return
    const ctx = canvas.getContext('2d')
    if (!ctx) return

    const game = new Game()
    const renderer = new GameRenderer()
    gameRef.current = game

    if (typeof window !== 'undefined') {
      const saved = localStorage.getItem(LS_KEY)
      if (saved) {
        game.bestScore = parseInt(saved, 10) || 0
        setBestScore(game.bestScore)
      }
    }

    disposedRef.current = false

    let dpr = 0
    let mediaQuery: MediaQueryList | null = null

    const onDprChange = () => {
      if (disposedRef.current) return
      syncCanvasSize()
    }

    const syncCanvasSize = () => {
      const next = window.devicePixelRatio || 1
      if (next === dpr) return

      dpr = next
      canvas.width = Math.round(GAME_WIDTH * dpr)
      canvas.height = Math.round(GAME_HEIGHT * dpr)

      mediaQuery?.removeEventListener('change', onDprChange)
      mediaQuery = window.matchMedia(`(resolution: ${dpr}dppx)`)
      mediaQuery.addEventListener('change', onDprChange, { once: true })
    }

    syncCanvasSize()
    window.addEventListener('resize', onDprChange)

    const FIXED_DT = 1000 / 60
    let accumulator = 0
    let lastTime = performance.now()
    let animId = 0
    let lastStatus: GameStatus = 'idle'

    const loop = (now: number) => {
      const delta = Math.min(now - lastTime, 100)
      lastTime = now
      accumulator += delta

      while (accumulator >= FIXED_DT) {
        game.update()
        accumulator -= FIXED_DT
      }

      const palette = paletteRef.current
      if (palette) {
        ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
        renderer.render(ctx, game, palette)
      }

      if (lottieContainerRef.current) {
        const { player, cameraY } = game
        const screenY = player.y - cameraY
        lottieContainerRef.current.style.transform = `translate(${player.x}px, ${screenY}px)`
      }

      if (game.status !== lastStatus && !disposedRef.current) {
        lastStatus = game.status
        setStatus(game.status)
        if (game.status === 'gameover') {
          setFinalScore(game.score)
          setBestScore(game.bestScore)
          localStorage.setItem(LS_KEY, String(game.bestScore))
        }
      }

      animId = requestAnimationFrame(loop)
    }

    animId = requestAnimationFrame(loop)
    return () => {
      disposedRef.current = true
      cancelAnimationFrame(animId)
      window.removeEventListener('resize', onDprChange)
      mediaQuery?.removeEventListener('change', onDprChange)
    }
  }, [paletteRef])

  useGameKeyboard(gameRef, surfaceRef)

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
      if (gameRef.current?.status !== 'playing') {
        gameRef.current?.reset()
      }
    },
    [updateTouchInput],
  )

  const onTouchEnd = useCallback(
    (e: ReactTouchEvent<HTMLDivElement>) => {
      updateTouchInput(e.touches, e.currentTarget.getBoundingClientRect())
    },
    [updateTouchInput],
  )

  const startGame = useCallback(() => {
    gameRef.current?.reset()
  }, [])

  return (
    <div
      ref={frameRef}
      className="flex w-full items-start justify-center"
      style={{ height: GAME_HEIGHT * scale }}
    >
      <div
        ref={surfaceRef}
        className="relative origin-top"
        style={{
          width: GAME_WIDTH,
          height: GAME_HEIGHT,
          transform: `scale(${scale})`,
        }}
        onTouchStart={onTouchStart}
        onTouchMove={onTouchStart}
        onTouchEnd={onTouchEnd}
        onTouchCancel={onTouchEnd}
      >
        <canvas
          ref={canvasRef}
          role="img"
          aria-label={t('raccoonjump.name')}
          className="absolute left-0 top-0 h-full w-full rounded-xl ring-1 ring-border"
        >
          {t('raccoonjump.description')}
        </canvas>

        <p role="status" aria-live="polite" className="sr-only">
          {status === 'gameover' ? t('raccoonjump.score', { score: finalScore }) : ''}
        </p>

        <DoodleJumpPlayer containerRef={lottieContainerRef} status={status} />

        <DoodleJumpOverlays
          status={status}
          finalScore={finalScore}
          bestScore={bestScore}
          onStart={startGame}
        />
      </div>
    </div>
  )
}
