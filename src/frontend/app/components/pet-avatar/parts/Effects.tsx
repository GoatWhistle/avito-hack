const HEART_PATH = 'M 0 4 q -6 -7 -1 -11 q 3 -2 1 3 q 2 -5 5 -3 q 5 4 -1 11 z'

const HEART_SPOTS = [
  { key: 'a', x: 100, y: 78, modifier: '' },
  { key: 'b', x: 82, y: 84, modifier: ' pet-avatar__heart--b' },
  { key: 'c', x: 118, y: 82, modifier: ' pet-avatar__heart--c' },
] as const

const ZZZ_SPOTS = [
  { key: 'a', x: 136, y: 74, size: 13, modifier: '' },
  { key: 'b', x: 146, y: 62, size: 10, modifier: ' pet-avatar__zzz--b' },
  { key: 'c', x: 154, y: 52, size: 8, modifier: ' pet-avatar__zzz--c' },
] as const

const SPARK_SPOTS = [
  { key: 'a', x: 44, y: 66, delay: '0s' },
  { key: 'b', x: 156, y: 60, delay: '0.4s' },
  { key: 'c', x: 60, y: 34, delay: '0.8s' },
  { key: 'd', x: 142, y: 30, delay: '1.2s' },
] as const

const CONFETTI_SPOTS = [
  { key: 'a', x: 52, y: 150, fill: 'var(--pet-accent-blue)', delay: '0s' },
  { key: 'b', x: 78, y: 138, fill: 'var(--pet-accent-green)', delay: '0.1s' },
  { key: 'c', x: 104, y: 132, fill: 'var(--pet-gold)', delay: '0.2s' },
  { key: 'd', x: 128, y: 140, fill: 'var(--pet-accent-pink)', delay: '0.3s' },
  { key: 'e', x: 150, y: 152, fill: 'var(--pet-accent-blue)', delay: '0.15s' },
] as const

export function Hearts() {
  return (
    <g>
      {HEART_SPOTS.map((spot) => (
        <path
          key={spot.key}
          className={`pet-avatar__heart${spot.modifier}`}
          d={HEART_PATH}
          transform={`translate(${spot.x} ${spot.y}) scale(1.6)`}
          fill="var(--pet-accent-pink)"
        />
      ))}
    </g>
  )
}

export function SleepZzz() {
  return (
    <g aria-hidden="true">
      {ZZZ_SPOTS.map((spot) => (
        <text
          key={spot.key}
          className={`pet-avatar__zzz${spot.modifier}`}
          x={spot.x}
          y={spot.y}
          fontSize={spot.size}
          fontWeight={700}
          fill="var(--pet-accent-blue)"
        >
          z
        </text>
      ))}
    </g>
  )
}

export function Sparks() {
  return (
    <g>
      {SPARK_SPOTS.map((spot) => (
        <path
          key={spot.key}
          className="pet-avatar__spark"
          d="M 0 -7 L 2 -2 L 7 0 L 2 2 L 0 7 L -2 2 L -7 0 L -2 -2 z"
          transform={`translate(${spot.x} ${spot.y})`}
          style={{ animationDelay: spot.delay }}
          fill="var(--pet-gold)"
        />
      ))}
    </g>
  )
}

export function Confetti() {
  return (
    <g>
      {CONFETTI_SPOTS.map((spot) => (
        <rect
          key={spot.key}
          className="pet-avatar__confetti"
          x={spot.x}
          y={spot.y}
          width={7}
          height={11}
          rx={2}
          fill={spot.fill}
          style={{ animationDelay: spot.delay }}
        />
      ))}
    </g>
  )
}

export function Aura() {
  return (
    <ellipse
      className="pet-avatar__aura"
      cx={100}
      cy={120}
      rx={78}
      ry={80}
      fill="var(--pet-aura)"
    />
  )
}

export function Flash() {
  return (
    <circle
      className="pet-avatar__flash"
      cx={100}
      cy={120}
      r={62}
      fill="var(--pet-gold)"
    />
  )
}

export function GroundShadow() {
  return (
    <ellipse
      cx={100}
      cy={196}
      rx={58}
      ry={9}
      fill="var(--pet-shadow)"
      opacity={0.16}
    />
  )
}
