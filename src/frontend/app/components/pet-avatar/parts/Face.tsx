import type { EyeOffset } from '../types'

const EYE_LEFT_X = 82
const EYE_RIGHT_X = 118
const EYE_Y = 92

export interface EyesProps {
  offset: EyeOffset
  closed: boolean
  sad: boolean
}

export function Eyes({ offset, closed, sad }: EyesProps) {
  const radius = sad ? 7 : 8.5

  return (
    <g>
      {[EYE_LEFT_X, EYE_RIGHT_X].map((cx) => (
        <g key={cx}>
          <ellipse
            cx={cx}
            cy={EYE_Y}
            rx={radius}
            ry={sad ? 6 : radius}
            fill="var(--pet-eye-white)"
          />
          <circle
            className="pet-avatar__pupil"
            cx={cx}
            cy={EYE_Y}
            r={4.4}
            fill="var(--pet-eye-pupil)"
            style={{ transform: `translate(${offset.x}px, ${offset.y}px)` }}
          />
          <circle
            cx={cx + 1.8}
            cy={EYE_Y - 2.4}
            r={1.5}
            fill="var(--pet-eye-white)"
            opacity={0.9}
          />
          <ellipse
            className={`pet-avatar__lid${closed ? ' pet-avatar__lid--closed' : ''}`}
            cx={cx}
            cy={EYE_Y}
            rx={radius + 0.8}
            ry={radius + 0.8}
            fill="var(--pet-mask-dark)"
          />
          {closed && (
            <path
              d={`M ${cx - radius} ${EYE_Y} q ${radius} ${sad ? -5 : 6} ${radius * 2} 0`}
              stroke="var(--pet-nose)"
              strokeWidth={2}
              strokeLinecap="round"
              fill="none"
              opacity={0.7}
            />
          )}
        </g>
      ))}
    </g>
  )
}

export function Mask({ smiling }: { smiling: boolean }) {
  return (
    <g>
      <path
        d="M 66 86 q 8 -16 34 -16 q 26 0 34 16 q -6 20 -34 20 q -28 0 -34 -20 z"
        fill="var(--pet-mask-dark)"
        opacity={0.92}
      />
      <path
        d="M 74 78 q 26 -12 52 0 q -26 -6 -52 0 z"
        fill="var(--pet-mask-light)"
        opacity={0.5}
      />
      <ellipse cx={100} cy={108} rx={14} ry={11} fill="var(--pet-mask-light)" />
      <ellipse cx={100} cy={104} rx={5.6} ry={4.4} fill="var(--pet-nose)" />
      <path
        d={smiling ? 'M 92 112 q 8 8 16 0' : 'M 92 114 q 8 -5 16 0'}
        stroke="var(--pet-nose)"
        strokeWidth={2.2}
        strokeLinecap="round"
        fill="none"
      />
    </g>
  )
}

export function Tongue() {
  return (
    <ellipse
      cx={100}
      cy={116}
      rx={6}
      ry={4.5}
      fill="var(--pet-accent-pink)"
      opacity={0.9}
    />
  )
}
