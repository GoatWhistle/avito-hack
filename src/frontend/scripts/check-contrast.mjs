import { palette, semanticDark, semanticLight } from '../app/tokens/colors.ts'

const AA_BODY = 4.5
const AA_LARGE = 3
const AA_UI = 3

function resolve(scope, token, seen = new Set()) {
  if (seen.has(token)) return null
  seen.add(token)
  const map = scope === 'dark' ? semanticDark : semanticLight
  const raw = map[token]
  if (raw == null) return null
  const paletteRef = raw.match(/^var\(--color-([\w-]+)\)$/)
  if (paletteRef) return palette[paletteRef[1]] ?? null
  const semanticRef = raw.match(/^var\(--([\w-]+)\)$/)
  if (semanticRef) return resolve(scope, semanticRef[1], seen)
  return raw
}

function oklchToLinearSrgb(str) {
  const m = str.match(
    /oklch\(\s*([\d.]+)\s+([\d.]+)\s+([\d.]+)\s*(?:\/\s*([\d.]+)\s*)?\)/,
  )
  if (!m) return null
  const L = parseFloat(m[1])
  const C = parseFloat(m[2])
  const H = (parseFloat(m[3]) * Math.PI) / 180
  const alpha = m[4] === undefined ? 1 : parseFloat(m[4])
  const a = C * Math.cos(H)
  const b = C * Math.sin(H)
  const l = (L + 0.3963377774 * a + 0.2158037573 * b) ** 3
  const mm = (L - 0.1055613458 * a - 0.0638541728 * b) ** 3
  const s = (L - 0.0894841775 * a - 1.291485548 * b) ** 3
  return {
    r: 4.0767416621 * l - 3.3077115913 * mm + 0.2309699292 * s,
    g: -1.2684380046 * l + 2.6097574011 * mm - 0.3413193965 * s,
    b: -0.0041960863 * l - 0.7034186147 * mm + 1.707614701 * s,
    alpha,
  }
}

const clamp = (x) => Math.min(1, Math.max(0, x))

const relLuminance = ({ r, g, b }) =>
  0.2126 * clamp(r) + 0.7152 * clamp(g) + 0.0722 * clamp(b)

function composite(fg, bg) {
  if (fg.alpha >= 1) return fg
  const a = fg.alpha
  return {
    r: fg.r * a + bg.r * (1 - a),
    g: fg.g * a + bg.g * (1 - a),
    b: fg.b * a + bg.b * (1 - a),
    alpha: 1,
  }
}

function contrast(fgRaw, bg) {
  const fg = composite(fgRaw, bg)
  const l1 = relLuminance(fg)
  const l2 = relLuminance(bg)
  const [hi, lo] = l1 > l2 ? [l1, l2] : [l2, l1]
  return (hi + 0.05) / (lo + 0.05)
}

const PAIRS = [
  ['foreground', 'background', AA_BODY, 'body text'],
  ['card-foreground', 'card', AA_BODY, 'card text'],
  ['popover-foreground', 'popover', AA_BODY, 'popover text'],
  ['muted-foreground', 'background', AA_BODY, 'muted on bg'],
  ['muted-foreground', 'muted', AA_BODY, 'muted on muted'],
  ['primary-foreground', 'primary', AA_BODY, 'primary button'],
  ['secondary-foreground', 'secondary', AA_BODY, 'secondary button'],
  ['accent-foreground', 'accent', AA_BODY, 'accent button'],
  ['destructive-foreground', 'destructive', AA_BODY, 'destructive button'],
  ['success-foreground', 'success', AA_BODY, 'success button'],
  ['warning-foreground', 'warning', AA_BODY, 'warning button'],
  ['info-foreground', 'info', AA_BODY, 'info button'],
  ['primary-subtle-foreground', 'primary-subtle', AA_BODY, 'primary subtle'],
  ['accent-subtle-foreground', 'accent-subtle', AA_BODY, 'accent subtle'],
  [
    'destructive-subtle-foreground',
    'destructive-subtle',
    AA_BODY,
    'destructive subtle',
  ],
  ['success-subtle-foreground', 'success-subtle', AA_BODY, 'success subtle'],
  ['warning-subtle-foreground', 'warning-subtle', AA_BODY, 'warning subtle'],
  ['info-subtle-foreground', 'info-subtle', AA_BODY, 'info subtle'],
  ['sidebar-foreground', 'sidebar', AA_BODY, 'sidebar text'],
  ['sidebar-primary-foreground', 'sidebar-primary', AA_BODY, 'sidebar primary'],
  ['sidebar-accent-foreground', 'sidebar-accent', AA_BODY, 'sidebar accent'],
  ['primary', 'background', AA_LARGE, 'primary as large text'],
  ['destructive', 'background', AA_LARGE, 'destructive as text'],
  ['ring', 'background', AA_UI, 'focus ring'],
  ['xp-fill', 'xp-track', AA_UI, 'xp bar'],
  ['input', 'background', AA_UI, 'input border'],
]

const failures = []
const rows = []

for (const scope of ['light', 'dark']) {
  for (const [fgToken, bgToken, threshold, label] of PAIRS) {
    const fgRaw = resolve(scope, fgToken)
    const bgRaw = resolve(scope, bgToken)
    if (!fgRaw || !bgRaw) {
      failures.push(`${scope} ${label}: unresolved (${fgToken} / ${bgToken})`)
      continue
    }
    const fg = oklchToLinearSrgb(fgRaw)
    const bg = oklchToLinearSrgb(bgRaw)
    if (!fg || !bg) {
      failures.push(`${scope} ${label}: unparsed color`)
      continue
    }
    const ratio = contrast(fg, bg)
    const ok = ratio >= threshold
    if (!ok) {
      failures.push(
        `${scope} ${label}: ${ratio.toFixed(2)}:1 < ${threshold}:1 (${fgToken} on ${bgToken})`,
      )
    }
    rows.push(
      `${ok ? 'PASS' : 'FAIL'}  ${scope.padEnd(5)} ${label.padEnd(24)} ${ratio.toFixed(2)}:1 (min ${threshold})`,
    )
  }
}

for (const row of rows) console.log(row)
console.log(`\n${rows.length} pairs checked, ${failures.length} failing`)

if (failures.length > 0) {
  console.error('\nContrast failures:')
  for (const f of failures) console.error(`  ${f}`)
  process.exit(1)
}
