import {
  GAME_HEIGHT,
  GAME_WIDTH,
  PLATFORM_H,
  PLATFORM_W,
  type Game,
  type Platform,
} from './game'
import type { GamePalette } from './palette'

const GRID_STEP = 64
const SCORE_FONT = 'bold 24px "Avito Sans", ui-sans-serif, system-ui, sans-serif'
const PLATFORM_RADIUS = PLATFORM_H / 2
const SPRING_W = 22
const SPRING_H = 12

export class GameRenderer {
  private backgroundGradient: CanvasGradient | null = null
  private backgroundPalette: GamePalette | null = null

  render(
    ctx: CanvasRenderingContext2D,
    game: Game,
    palette: GamePalette,
  ): void {
    this.renderBackground(ctx, game.cameraY, palette)

    ctx.save()
    ctx.translate(0, -game.cameraY)

    for (const platform of game.platforms) {
      this.renderPlatform(ctx, platform, palette)
    }

    ctx.restore()

    if (game.status === 'playing') this.renderScore(ctx, game.score, palette)
  }

  private renderBackground(
    ctx: CanvasRenderingContext2D,
    cameraY: number,
    palette: GamePalette,
  ): void {
    if (this.backgroundPalette !== palette || !this.backgroundGradient) {
      const gradient = ctx.createLinearGradient(0, 0, 0, GAME_HEIGHT)
      gradient.addColorStop(0, palette.backgroundAlt)
      gradient.addColorStop(1, palette.background)
      this.backgroundGradient = gradient
      this.backgroundPalette = palette
    }

    ctx.fillStyle = this.backgroundGradient
    ctx.fillRect(0, 0, GAME_WIDTH, GAME_HEIGHT)

    ctx.fillStyle = palette.grid
    const offset = ((-cameraY * 0.35) % GRID_STEP) - GRID_STEP
    for (let y = offset; y < GAME_HEIGHT; y += GRID_STEP) {
      ctx.fillRect(0, y, GAME_WIDTH, 1)
    }
  }

  private renderScore(
    ctx: CanvasRenderingContext2D,
    score: number,
    palette: GamePalette,
  ): void {
    const label = String(score)

    ctx.font = SCORE_FONT
    ctx.textAlign = 'left'
    ctx.textBaseline = 'middle'

    const pillW = Math.max(56, ctx.measureText(label).width + 28)

    ctx.fillStyle = palette.scorePill
    ctx.beginPath()
    ctx.roundRect(14, 14, pillW, 38, 19)
    ctx.fill()

    ctx.strokeStyle = palette.scorePillRing
    ctx.lineWidth = 1
    ctx.stroke()

    ctx.fillStyle = palette.scoreText
    ctx.fillText(label, 28, 34)
  }

  private renderPlatform(
    ctx: CanvasRenderingContext2D,
    p: Platform,
    palette: GamePalette,
  ): void {
    let fill = palette.normal
    let edge = palette.normalEdge

    if (p.type === 'moving') {
      fill = palette.moving
      edge = palette.movingEdge
    } else if (p.type === 'breakable') {
      fill = palette.breakable
      edge = palette.breakableEdge
    }

    ctx.save()

    if (p.broken) {
      ctx.globalAlpha = Math.max(0, 1 - p.breakTimer / 30)
      ctx.translate(p.x + PLATFORM_W / 2, p.y + PLATFORM_H / 2)
      ctx.rotate(p.breakTimer * 0.05)
      ctx.translate(-(p.x + PLATFORM_W / 2), -(p.y + PLATFORM_H / 2))
    }

    if (p.hasSpring) this.renderSpring(ctx, p, palette)

    ctx.fillStyle = edge
    ctx.beginPath()
    ctx.roundRect(p.x, p.y + 3, PLATFORM_W, PLATFORM_H - 1, PLATFORM_RADIUS)
    ctx.fill()

    ctx.fillStyle = fill
    ctx.beginPath()
    ctx.roundRect(p.x, p.y, PLATFORM_W, PLATFORM_H, PLATFORM_RADIUS)
    ctx.fill()

    ctx.fillStyle = palette.platformGloss
    ctx.beginPath()
    ctx.roundRect(p.x + 5, p.y + 3, PLATFORM_W - 10, 4, 2)
    ctx.fill()

    if (p.type === 'moving') this.renderMovingMarks(ctx, p, palette)
    else if (p.type === 'breakable') this.renderCracks(ctx, p, edge)

    ctx.restore()
  }

  private renderMovingMarks(
    ctx: CanvasRenderingContext2D,
    p: Platform,
    palette: GamePalette,
  ): void {
    const cy = p.y + PLATFORM_H / 2
    const cx = p.x + PLATFORM_W / 2
    const dir = p.vx >= 0 ? 1 : -1

    ctx.strokeStyle = palette.platformGloss
    ctx.lineWidth = 2
    ctx.lineCap = 'round'
    ctx.lineJoin = 'round'

    ctx.beginPath()
    for (let i = 0; i < 2; i++) {
      const ox = cx + dir * (i * 7 - 3)
      ctx.moveTo(ox - dir * 3, cy - 4)
      ctx.lineTo(ox + dir * 3, cy)
      ctx.lineTo(ox - dir * 3, cy + 4)
    }
    ctx.stroke()
  }

  private renderCracks(
    ctx: CanvasRenderingContext2D,
    p: Platform,
    edge: string,
  ): void {
    ctx.strokeStyle = edge
    ctx.lineWidth = 1.5
    ctx.lineCap = 'round'
    ctx.lineJoin = 'round'

    ctx.beginPath()
    ctx.moveTo(p.x + PLATFORM_W * 0.26, p.y + 3)
    ctx.lineTo(p.x + PLATFORM_W * 0.4, p.y + PLATFORM_H - 4)
    ctx.lineTo(p.x + PLATFORM_W * 0.54, p.y + 4)
    ctx.lineTo(p.x + PLATFORM_W * 0.68, p.y + PLATFORM_H - 3)
    ctx.stroke()
  }

  private renderSpring(
    ctx: CanvasRenderingContext2D,
    p: Platform,
    palette: GamePalette,
  ): void {
    const sx = p.x + PLATFORM_W / 2 - SPRING_W / 2
    const sy = p.y - SPRING_H

    ctx.strokeStyle = palette.springBody
    ctx.lineWidth = 2.5
    ctx.lineCap = 'round'

    ctx.beginPath()
    for (let i = 0; i < 3; i++) {
      const y = sy + 3 + i * 3.5
      ctx.moveTo(sx + 3, y)
      ctx.lineTo(sx + SPRING_W - 3, y + 1.5)
    }
    ctx.stroke()

    ctx.fillStyle = palette.springCap
    ctx.beginPath()
    ctx.roundRect(sx, sy - 2, SPRING_W, 5, 2.5)
    ctx.fill()
  }
}
