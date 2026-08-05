import { petPalette } from '@/shared/design';

const HEART_PATH = 'M 0 4 q -6 -7 -1 -11 q 3 -2 1 3 q 2 -5 5 -3 q 5 4 -1 11 z';
const HEART_SPOTS = [
  { key: 'a', x: 100, y: 78, className: '' },
  { key: 'b', x: 82, y: 84, className: 'pet-avatar__heart--b' },
  { key: 'c', x: 118, y: 82, className: 'pet-avatar__heart--c' },
] as const;

const ZZZ_GLYPH = 'z';

const ZZZ_SPOTS = [
  { key: 'a', x: 136, y: 74, size: 13, className: '' },
  { key: 'b', x: 146, y: 62, size: 10, className: 'pet-avatar__zzz--b' },
  { key: 'c', x: 154, y: 52, size: 8, className: 'pet-avatar__zzz--c' },
] as const;

const SPARK_SPOTS = [
  { key: 'a', x: 44, y: 66 },
  { key: 'b', x: 156, y: 60 },
  { key: 'c', x: 60, y: 34 },
  { key: 'd', x: 142, y: 30 },
] as const;

export function Hearts() {
  return (
    <g>
      {HEART_SPOTS.map((spot) => (
        <path
          key={spot.key}
          className={`pet-avatar__heart ${spot.className}`}
          d={HEART_PATH}
          transform={`translate(${spot.x} ${spot.y}) scale(1.6)`}
          fill={petPalette.accentPink}
        />
      ))}
    </g>
  );
}

export function SleepZzz() {
  return (
    <g>
      {ZZZ_SPOTS.map((spot) => (
        <text
          key={spot.key}
          className={`pet-avatar__zzz ${spot.className}`}
          x={spot.x}
          y={spot.y}
          fontSize={spot.size}
          fontWeight={700}
          fill={petPalette.accentBlue}
        >
          {ZZZ_GLYPH}
        </text>
      ))}
    </g>
  );
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
          fill={petPalette.streakCore}
        />
      ))}
    </g>
  );
}

export function Aura() {
  return (
    <ellipse
      className="pet-avatar__aura"
      cx={100}
      cy={120}
      rx={78}
      ry={80}
      fill={petPalette.auraGlow}
    />
  );
}

export function Egg({ hatching }: { hatching: boolean }) {
  return (
    <g className="pet-avatar__body">
      <ellipse cx={100} cy={120} rx={54} ry={68} fill={petPalette.eggShell} />
      <ellipse cx={100} cy={120} rx={54} ry={68} fill={petPalette.furBase} opacity={0.08} />
      <circle cx={78} cy={102} r={9} fill={petPalette.eggSpot} opacity={0.35} />
      <circle cx={118} cy={132} r={12} fill={petPalette.accentGreen} opacity={0.3} />
      <circle cx={110} cy={92} r={6} fill={petPalette.accentGreen} opacity={0.25} />
      {hatching && (
        <path
          d="M 52 116 l 16 -12 l 12 12 l 14 -14 l 12 14 l 14 -12 l 16 12"
          stroke={petPalette.eggCrack}
          strokeWidth={3.4}
          strokeLinejoin="round"
          fill="none"
        />
      )}
    </g>
  );
}
