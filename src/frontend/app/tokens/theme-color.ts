import { palette, type PaletteToken } from './palette'

const OKLCH_PATTERN = /oklch\(\s*([\d.]+)\s+([\d.]+)\s+([\d.]+)\s*\)/

const linearToSrgb = (channel: number) => {
  const clamped = Math.min(1, Math.max(0, channel))
  const encoded =
    clamped <= 0.0031308
      ? clamped * 12.92
      : 1.055 * clamped ** (1 / 2.4) - 0.055
  return Math.round(encoded * 255)
}

const toHexPair = (value: number) => value.toString(16).padStart(2, '0')

export const oklchToHex = (value: string): string => {
  const parsed = OKLCH_PATTERN.exec(value)
  if (!parsed) return value

  const lightness = Number(parsed[1])
  const chroma = Number(parsed[2])
  const hue = (Number(parsed[3]) * Math.PI) / 180
  const a = chroma * Math.cos(hue)
  const b = chroma * Math.sin(hue)

  const l = (lightness + 0.3963377774 * a + 0.2158037573 * b) ** 3
  const m = (lightness - 0.1055613458 * a - 0.0638541728 * b) ** 3
  const s = (lightness - 0.0894841775 * a - 1.291485548 * b) ** 3

  const red = 4.0767416621 * l - 3.3077115913 * m + 0.2309699292 * s
  const green = -1.2684380046 * l + 2.6097574011 * m - 0.3413193965 * s
  const blue = -0.0041960863 * l - 0.7034186147 * m + 1.707614701 * s

  return `#${toHexPair(linearToSrgb(red))}${toHexPair(linearToSrgb(green))}${toHexPair(linearToSrgb(blue))}`
}

const hexFor = (token: PaletteToken) => oklchToHex(palette[token])

export const themeColors = {
  light: hexFor('neutral-0'),
  dark: hexFor('neutral-950'),
} as const
