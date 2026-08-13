export interface GamePalette {
  background: string
  backgroundAlt: string
  scoreText: string
  scorePill: string
  scorePillRing: string
  normal: string
  normalEdge: string
  moving: string
  movingEdge: string
  breakable: string
  breakableEdge: string
  platformGloss: string
  springBody: string
  springCap: string
  grid: string
  isDark: boolean
}

type Token =
  | '--background'
  | '--foreground'
  | '--card'
  | '--muted'
  | '--muted-foreground'
  | '--primary'
  | '--accent'
  | '--warning'
  | '--destructive'

const mix = (color: string, other: string, percent: number): string =>
  `color-mix(in oklch, ${color}, ${other} ${percent}%)`

const alpha = (color: string, percent: number): string =>
  `color-mix(in oklch, ${color} ${percent}%, transparent)`

export const readGamePalette = (root: HTMLElement): GamePalette => {
  const styles = getComputedStyle(root)
  const isDark = root.classList.contains('dark')

  const read = (token: Token): string => styles.getPropertyValue(token).trim()

  const background = read('--background')
  const card = read('--card')
  const foreground = read('--foreground')
  const muted = read('--muted')
  const mutedForeground = read('--muted-foreground')
  const primary = read('--primary')
  const accent = read('--accent')
  const breakable = isDark ? read('--warning') : read('--destructive')

  const shade = isDark ? foreground : background
  const deepen = isDark ? background : foreground

  return {
    background,
    backgroundAlt: mix(background, card, 60),
    scoreText: foreground,
    scorePill: alpha(card, isDark ? 55 : 70),
    scorePillRing: alpha(foreground, 12),
    normal: accent,
    normalEdge: mix(accent, deepen, 30),
    moving: primary,
    movingEdge: mix(primary, deepen, 30),
    breakable,
    breakableEdge: mix(breakable, deepen, 30),
    platformGloss: alpha(shade, isDark ? 16 : 45),
    springBody: mutedForeground,
    springCap: mix(muted, mutedForeground, 45),
    grid: alpha(foreground, isDark ? 6 : 5),
    isDark,
  }
}
