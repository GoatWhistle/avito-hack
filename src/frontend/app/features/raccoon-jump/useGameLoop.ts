import { useEffect, type RefObject } from 'react'
import { Game, GAME_HEIGHT, GAME_WIDTH, type GameStatus } from './game'
import { GameRenderer } from './render'
import type { GamePalette } from './palette'

interface GameLoopOptions {
  canvasRef: RefObject<HTMLCanvasElement | null>
  gameRef: RefObject<Game | null>
  lottieContainerRef: RefObject<HTMLDivElement | null>
  disposedRef: RefObject<boolean>
  paletteRef: RefObject<GamePalette | null>
  scaleRef: RefObject<number>
  onStatusChange: (status: GameStatus) => void
  onGameOverRef: RefObject<(score: number, collected: number[]) => void>
}

export const useGameLoop = ({
  canvasRef,
  gameRef,
  lottieContainerRef,
  disposedRef,
  paletteRef,
  scaleRef,
  onStatusChange,
  onGameOverRef,
}: GameLoopOptions): void => {
  useEffect(() => {
    const canvas = canvasRef.current
    if (!canvas) return
    const ctx = canvas.getContext('2d')
    if (!ctx) return

    const game = new Game()
    const renderer = new GameRenderer()
    gameRef.current = game

    disposedRef.current = false

    let dpr = 0
    let renderScale = 0
    let appliedScale = 0
    let mediaQuery: MediaQueryList | null = null

    const onDprChange = () => {
      if (disposedRef.current) return
      dpr = 0
      syncCanvasSize()
    }

    const syncCanvasSize = () => {
      const nextDpr = window.devicePixelRatio || 1
      const nextScale = scaleRef.current || 1
      if (nextDpr === dpr && nextScale === appliedScale) return

      dpr = nextDpr
      appliedScale = nextScale
      renderScale = dpr * nextScale
      canvas.width = Math.round(GAME_WIDTH * renderScale)
      canvas.height = Math.round(GAME_HEIGHT * renderScale)

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

      if ((scaleRef.current || 1) !== appliedScale) syncCanvasSize()

      const palette = paletteRef.current
      if (palette) {
        ctx.setTransform(renderScale, 0, 0, renderScale, 0, 0)
        renderer.render(ctx, game, palette)
      }

      if (lottieContainerRef.current) {
        const { player, cameraY } = game
        const screenY = player.y - cameraY
        lottieContainerRef.current.style.transform = `translate(${player.x}px, ${screenY}px)`
      }

      if (game.status !== lastStatus && !disposedRef.current) {
        lastStatus = game.status
        onStatusChange(game.status)
        if (game.status === 'gameover') {
          onGameOverRef.current(game.score, [...game.collected])
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
  }, [
    canvasRef,
    gameRef,
    lottieContainerRef,
    disposedRef,
    paletteRef,
    scaleRef,
    onStatusChange,
    onGameOverRef,
  ])
}
