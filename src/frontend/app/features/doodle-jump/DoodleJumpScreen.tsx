import type {
  TouchEvent as ReactTouchEvent,
  TouchList as ReactTouchList,
} from 'react'
import { useCallback, useEffect, useRef, useState } from 'react'
import { DotLottieReact } from '@lottiefiles/dotlottie-react'
import {
  Game,
  GAME_HEIGHT,
  GAME_WIDTH,
  type GameStatus,
  PLAYER_H,
  PLAYER_W,
} from './game'

const LS_KEY = 'doodle-jump-best'

export function DoodleJumpScreen() {
  const canvasRef = useRef<HTMLCanvasElement>(null)
  const gameRef = useRef<Game | null>(null)
  const lottieContainerRef = useRef<HTMLDivElement>(null)

  const [status, setStatus] = useState<GameStatus>('idle')
  const [finalScore, setFinalScore] = useState(0)
  const [bestScore, setBestScore] = useState(0)

  useEffect(() => {
    const canvas = canvasRef.current
    if (!canvas) return
    const ctx = canvas.getContext('2d')
    if (!ctx) return

    const game = new Game()
    gameRef.current = game

    if (typeof window !== 'undefined') {
      const saved = localStorage.getItem(LS_KEY)
      if (saved) {
        game.bestScore = parseInt(saved, 10) || 0
        setBestScore(game.bestScore)
      }
    }

    const dpr = window.devicePixelRatio || 1
    canvas.width = GAME_WIDTH * dpr
    canvas.height = GAME_HEIGHT * dpr

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

      ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
      game.render(ctx)

      if (lottieContainerRef.current) {
        const { player, cameraY } = game
        const screenY = player.y - cameraY
        lottieContainerRef.current.style.transform = `translate(${player.x}px, ${screenY}px)`
      }

      if (game.status !== lastStatus) {
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
    return () => cancelAnimationFrame(animId)
  }, [])

  useEffect(() => {
    const game = gameRef.current
    if (!game) return

    const keys = {
      left: ['ArrowLeft', 'a', 'A', 'ф', 'Ф'],
      right: ['ArrowRight', 'd', 'D', 'в', 'В'],
    }

    const onKeyDown = (e: KeyboardEvent) => {
      if (keys.left.includes(e.key)) game.input.left = true
      if (keys.right.includes(e.key)) game.input.right = true
      if ((e.key === ' ' || e.key === 'Enter') && game.status !== 'playing') {
        game.reset()
      }
    }

    const onKeyUp = (e: KeyboardEvent) => {
      if (keys.left.includes(e.key)) game.input.left = false
      if (keys.right.includes(e.key)) game.input.right = false
    }

    window.addEventListener('keydown', onKeyDown)
    window.addEventListener('keyup', onKeyUp)
    return () => {
      window.removeEventListener('keydown', onKeyDown)
      window.removeEventListener('keyup', onKeyUp)
    }
  }, [])

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
    <div className="flex items-center justify-center bg-background py-6">
      <div
        className="relative"
        style={{ width: GAME_WIDTH, height: GAME_HEIGHT }}
        onTouchStart={onTouchStart}
        onTouchMove={onTouchStart}
        onTouchEnd={onTouchEnd}
        onTouchCancel={onTouchEnd}
      >
        <canvas
          ref={canvasRef}
          className="absolute left-0 top-0 h-full w-full rounded-xl ring-1 ring-border"
        />

        <div
          ref={lottieContainerRef}
          style={{
            position: 'absolute',
            top: 0,
            left: 0,
            width: PLAYER_W,
            height: PLAYER_H,
            pointerEvents: 'none',
            zIndex: 5,
          }}
        >
          <DotLottieReact
            src="/lottie/raccoon-adult.json"
            loop
            autoplay
            style={{
              width: '100%',
              height: '100%',
              display: 'block',
            }}
            onError={(err) => console.error('Lottie error:', err)}
          />
        </div>

        {status === 'idle' && (
          <div className="absolute inset-0 flex flex-col items-center justify-center gap-4 rounded-xl bg-black/60 p-6 text-center text-white">
            <h1 className="text-4xl font-bold">Doodle Jump</h1>
            <p>ПК: ← → или A / D</p>
            <p>Телефон: касайся левой / правой стороны</p>
            <button
              onClick={startGame}
              className="mt-4 rounded-lg bg-primary px-6 py-3 text-lg font-semibold text-primary-foreground"
            >
              Играть
            </button>
          </div>
        )}

        {status === 'gameover' && (
          <div className="absolute inset-0 flex flex-col items-center justify-center gap-4 rounded-xl bg-black/60 p-6 text-center text-white">
            <h1 className="text-4xl font-bold">Game Over</h1>
            <p className="text-xl">Счёт: {finalScore}</p>
            <p className="text-lg text-muted-foreground">Рекорд: {bestScore}</p>
            <button
              onClick={startGame}
              className="mt-4 rounded-lg bg-primary px-6 py-3 text-lg font-semibold text-primary-foreground"
            >
              Ещё раз
            </button>
          </div>
        )}

        {status === 'playing' && (
          <div
            className="pointer-events-none absolute bottom-4 w-full select-none px-6 text-4xl text-black/20"
            aria-hidden
          >
            <div className="float-left">‹</div>
            <div className="float-right">›</div>
          </div>
        )}
      </div>
    </div>
  )
}
