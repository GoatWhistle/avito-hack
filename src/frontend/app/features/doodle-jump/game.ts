export const GAME_WIDTH = 400
export const GAME_HEIGHT = 640

export const PLAYER_W = 60
export const PLAYER_H = 60
const MOVE_SPEED = 6
const GRAVITY = 0.45
const JUMP_V = -13
const SPRING_V = -21

const PLATFORM_W = 70
const PLATFORM_H = 16
const GAP_MIN = 50
const GAP_MAX = 110
const MOVING_SPEED = 1.2
const START_Y = GAME_HEIGHT - 60

export type PlatformType = 'normal' | 'moving' | 'breakable'
export type GameStatus = 'idle' | 'playing' | 'gameover'

export interface Platform {
  x: number
  y: number
  type: PlatformType
  hasSpring: boolean
  vx: number
  broken: boolean
  breakTimer: number
}

const rand = (min: number, max: number) => Math.random() * (max - min) + min

export class Game {
  status: GameStatus = 'idle'
  score = 0
  bestScore = 0

  player = { x: 0, y: 0, vx: 0, vy: 0 }
  platforms: Platform[] = []
  input = { left: false, right: false }

  cameraY = 0
  private highestY = 0

  reset(): void {
    this.platforms = [
      {
        x: GAME_WIDTH / 2 - PLATFORM_W / 2,
        y: START_Y,
        type: 'normal',
        hasSpring: false,
        vx: 0,
        broken: false,
        breakTimer: 0,
      },
    ]

    this.player = {
      x: GAME_WIDTH / 2 - PLAYER_W / 2,
      y: START_Y - PLAYER_H,
      vx: 0,
      vy: JUMP_V,
    }

    this.cameraY = 0
    this.highestY = this.player.y
    this.score = 0
    this.status = 'playing'

    let y = START_Y
    while (y > -GAME_HEIGHT) {
      y -= rand(GAP_MIN, GAP_MAX)
      this.platforms.push(this.createPlatform(y))
    }
  }

  update(): void {
    if (this.status !== 'playing') return

    this.updatePlayer()
    this.updatePlatforms()
    this.checkCollisions()
    this.updateCamera()
    this.generatePlatforms()
    this.checkGameOver()
  }

  render(ctx: CanvasRenderingContext2D): void {
    ctx.fillStyle = '#f0f8ff'
    ctx.fillRect(0, 0, GAME_WIDTH, GAME_HEIGHT)

    ctx.save()
    ctx.translate(0, -this.cameraY)

    for (const p of this.platforms) this.renderPlatform(ctx, p)

    ctx.restore()

    if (this.status === 'playing') {
      ctx.fillStyle = '#333'
      ctx.font = 'bold 24px system-ui, sans-serif'
      ctx.textAlign = 'left'
      ctx.textBaseline = 'top'
      ctx.fillText(String(this.score), 16, 12)
    }
  }

  private updatePlayer(): void {
    this.player.vx = 0
    if (this.input.left) this.player.vx -= MOVE_SPEED
    if (this.input.right) this.player.vx += MOVE_SPEED

    this.player.x += this.player.vx

    if (this.player.x > GAME_WIDTH) this.player.x = -PLAYER_W
    if (this.player.x + PLAYER_W < 0) this.player.x = GAME_WIDTH

    this.player.vy += GRAVITY
    this.player.y += this.player.vy

    if (this.player.y < this.highestY) {
      this.highestY = this.player.y
      this.score = Math.max(
        0,
        Math.floor((START_Y - PLAYER_H - this.highestY) / 10)
      )
    }
  }

  private updatePlatforms(): void {
    for (const p of this.platforms) {
      if (p.type === 'moving' && !p.broken) {
        p.x += p.vx
        if (p.x <= 0 || p.x + PLATFORM_W >= GAME_WIDTH) p.vx *= -1
      }

      if (p.broken) {
        p.breakTimer += 1
        p.y += 3 + p.breakTimer * 0.5
      }
    }

    this.platforms = this.platforms.filter(
      (p) => p.y < this.cameraY + GAME_HEIGHT + 100
    )
  }

  private checkCollisions(): void {
    if (this.player.vy <= 0) return

    const bottom = this.player.y + PLAYER_H
    const prevBottom = bottom - this.player.vy

    for (const p of this.platforms) {
      if (p.broken) continue

      const overlapsX =
        this.player.x + PLAYER_W > p.x && this.player.x < p.x + PLATFORM_W

      const crossesY = prevBottom <= p.y && bottom >= p.y

      if (!overlapsX || !crossesY) continue

      this.player.y = p.y - PLAYER_H
      this.player.vy = p.hasSpring ? SPRING_V : JUMP_V

      if (p.type === 'breakable') {
        p.broken = true
      }

      break
    }
  }

  private updateCamera(): void {
    const target = this.player.y - GAME_HEIGHT * 0.4
    if (target < this.cameraY) this.cameraY = target
  }

  private generatePlatforms(): void {
    let topY = Infinity
    for (const p of this.platforms) {
      if (p.y < topY) topY = p.y
    }

    while (topY > this.cameraY - GAP_MAX) {
      topY -= rand(GAP_MIN, GAP_MAX)
      this.platforms.push(this.createPlatform(topY))
    }
  }

  private checkGameOver(): void {
    if (this.player.y > this.cameraY + GAME_HEIGHT + PLAYER_H) {
      this.status = 'gameover'
      this.bestScore = Math.max(this.bestScore, this.score)
    }
  }

  private createPlatform(y: number): Platform {
    const difficulty = Math.min(1, this.score / 500)

    const roll = Math.random()
    let type: PlatformType = 'normal'
    if (roll < 0.08 + difficulty * 0.15) type = 'breakable'
    else if (roll < 0.25 + difficulty * 0.2) type = 'moving'

    const hasSpring = type === 'normal' && Math.random() < 0.06

    return {
      x: rand(0, GAME_WIDTH - PLATFORM_W),
      y,
      type,
      hasSpring,
      vx:
        type === 'moving'
          ? (Math.random() > 0.5 ? 1 : -1) * MOVING_SPEED
          : 0,
      broken: false,
      breakTimer: 0,
    }
  }

  private renderPlatform(ctx: CanvasRenderingContext2D, p: Platform): void {
    const colors: Record<PlatformType, string> = {
      normal: '#4ade80',
      moving: '#60a5fa',
      breakable: '#ef4444',
    }

    ctx.save()

    if (p.broken) {
      const alpha = Math.max(0, 1 - p.breakTimer / 30)
      ctx.globalAlpha = alpha
      ctx.translate(p.x + PLATFORM_W / 2, p.y + PLATFORM_H / 2)
      ctx.rotate(p.breakTimer * 0.05)
      ctx.translate(-(p.x + PLATFORM_W / 2), -(p.y + PLATFORM_H / 2))
    }

    ctx.fillStyle = colors[p.type]
    ctx.beginPath()
    ctx.roundRect(p.x, p.y, PLATFORM_W, PLATFORM_H, 8)
    ctx.fill()

    if (p.type === 'breakable') {
      ctx.strokeStyle = 'rgba(0,0,0,0.35)'
      ctx.lineWidth = 1.5
      ctx.beginPath()
      ctx.moveTo(p.x + PLATFORM_W * 0.3, p.y + 2)
      ctx.lineTo(p.x + PLATFORM_W * 0.45, p.y + PLATFORM_H - 2)
      ctx.lineTo(p.x + PLATFORM_W * 0.6, p.y + 3)
      ctx.stroke()
    }

    if (p.hasSpring) {
      const sw = 20
      const sh = 10
      ctx.fillStyle = '#9ca3af'
      ctx.fillRect(p.x + PLATFORM_W / 2 - sw / 2, p.y - sh, sw, sh)
      ctx.fillStyle = '#6b7280'
      ctx.fillRect(p.x + PLATFORM_W / 2 - sw / 2, p.y - sh, sw, 3)
    }

    ctx.restore()
  }
}
