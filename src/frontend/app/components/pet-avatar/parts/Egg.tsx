import type { EyeOffset } from '../types'
import { Ears, Tail } from './Body'
import { Eyes, Mask } from './Face'

const CRACK_PATH = 'M 52 116 l 16 -12 l 12 12 l 14 -14 l 12 14 l 14 -12 l 16 12'

function Shell() {
  return (
    <>
      <ellipse cx={100} cy={120} rx={54} ry={68} fill="var(--pet-egg-shell)" />
      <ellipse
        cx={100}
        cy={120}
        rx={54}
        ry={68}
        fill="var(--pet-fur-base)"
        opacity={0.08}
      />
      <circle
        cx={78}
        cy={102}
        r={9}
        fill="var(--pet-accent-blue)"
        opacity={0.35}
      />
      <circle
        cx={118}
        cy={132}
        r={12}
        fill="var(--pet-accent-green)"
        opacity={0.3}
      />
      <circle
        cx={110}
        cy={92}
        r={6}
        fill="var(--pet-accent-green)"
        opacity={0.25}
      />
    </>
  )
}

export function Egg() {
  return (
    <g className="pet-avatar__body">
      <g className="pet-avatar__shell">
        <Shell />
      </g>
    </g>
  )
}

export interface HatchingEggProps {
  offset: EyeOffset
}

export function HatchingEgg({ offset }: HatchingEggProps) {
  return (
    <g className="pet-avatar__body">
      <g className="pet-avatar__hatchling">
        <Tail length={30} />
        <ellipse cx={100} cy={158} rx={40} ry={32} fill="var(--pet-fur-base)" />
        <ellipse
          cx={100}
          cy={162}
          rx={26}
          ry={24}
          fill="var(--pet-fur-belly)"
          opacity={0.85}
        />
        <g className="pet-avatar__head">
          <Ears spread={30} droopy={false} />
          <circle cx={100} cy={92} r={48} fill="var(--pet-fur-light)" />
          <circle
            cx={100}
            cy={92}
            r={48}
            fill="var(--pet-fur-base)"
            opacity={0.35}
          />
          <Mask smiling />
          <Eyes offset={offset} closed={false} sad={false} />
        </g>
      </g>

      <g className="pet-avatar__shell">
        <g className="pet-avatar__shell-top">
          <clipPath id="pet-shell-top-clip">
            <rect x={40} y={40} width={120} height={78} />
          </clipPath>
          <g clipPath="url(#pet-shell-top-clip)">
            <Shell />
          </g>
          <path
            d={CRACK_PATH}
            stroke="var(--pet-egg-crack)"
            strokeWidth={3.4}
            strokeLinejoin="round"
            fill="none"
          />
        </g>

        <g className="pet-avatar__shell-bottom">
          <clipPath id="pet-shell-bottom-clip">
            <rect x={40} y={118} width={120} height={80} />
          </clipPath>
          <g clipPath="url(#pet-shell-bottom-clip)">
            <Shell />
          </g>
          <path
            d={CRACK_PATH}
            stroke="var(--pet-egg-crack)"
            strokeWidth={3.4}
            strokeLinejoin="round"
            fill="none"
          />
        </g>
      </g>
    </g>
  )
}
