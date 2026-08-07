export interface EarsProps {
  spread: number
  droopy: boolean
}

export function Ears({ spread, droopy }: EarsProps) {
  const lift = droopy ? 12 : 0
  const rotate = droopy ? 26 : 0

  return (
    <g>
      <g transform={`rotate(${-rotate} ${100 - spread} 68)`}>
        <path
          d={`M ${100 - spread - 12} ${72 + lift} q -2 -22 14 -26 q 12 6 8 26 z`}
          fill="var(--pet-fur-base)"
        />
        <path
          d={`M ${100 - spread - 7} ${70 + lift} q -1 -14 9 -17 q 7 4 5 17 z`}
          fill="var(--pet-accent-pink)"
          opacity={0.55}
        />
      </g>
      <g transform={`rotate(${rotate} ${100 + spread} 68)`}>
        <path
          d={`M ${100 + spread + 12} ${72 + lift} q 2 -22 -14 -26 q -12 6 -8 26 z`}
          fill="var(--pet-fur-base)"
        />
        <path
          d={`M ${100 + spread + 7} ${70 + lift} q 1 -14 -9 -17 q -7 4 -5 17 z`}
          fill="var(--pet-accent-pink)"
          opacity={0.55}
        />
      </g>
    </g>
  )
}

export function Tail({ length }: { length: number }) {
  return (
    <g className="pet-avatar__tail">
      <path
        d={`M 138 150 q ${length} -6 ${length - 8} -${length}`}
        stroke="var(--pet-fur-base)"
        strokeWidth={20}
        strokeLinecap="round"
        fill="none"
      />
      <path
        d={`M 150 142 q ${length - 14} -4 ${length - 20} -${length - 12}`}
        stroke="var(--pet-mask-dark)"
        strokeWidth={7}
        strokeLinecap="round"
        fill="none"
        opacity={0.65}
      />
    </g>
  )
}

export function Scarf() {
  return (
    <g>
      <path
        d="M 70 128 q 30 14 60 0 q -4 12 -30 12 q -26 0 -30 -12 z"
        fill="var(--pet-accent-green)"
      />
      <path
        d="M 122 134 q 12 8 8 24 q -10 -2 -14 -18 z"
        fill="var(--pet-accent-blue)"
      />
    </g>
  )
}

export function Crown() {
  return (
    <g>
      <path
        d="M 84 58 l 6 -18 l 10 12 l 10 -12 l 6 18 z"
        fill="var(--pet-gold)"
        stroke="var(--pet-flame)"
        strokeWidth={1.6}
        strokeLinejoin="round"
      />
      <circle cx={100} cy={38} r={3} fill="var(--pet-accent-blue)" />
    </g>
  )
}

export function Paws({ cy }: { cy: number }) {
  return (
    <g>
      <ellipse cx={66} cy={cy} rx={13} ry={9} fill="var(--pet-fur-dark)" />
      <ellipse cx={134} cy={cy} rx={13} ry={9} fill="var(--pet-fur-dark)" />
    </g>
  )
}
