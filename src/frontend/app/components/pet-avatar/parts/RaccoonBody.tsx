import { stageGeometry, stageScale } from '../mood'
import type { EyeOffset, PetEmotion, PetMood, PetStage } from '../types'
import { Crown, Ears, Paws, Scarf, Tail } from './Body'
import { Eyes, Mask, Tongue } from './Face'

export interface RaccoonBodyProps {
  stage: PetStage
  mood: PetMood
  emotion: PetEmotion | null
  offset: EyeOffset
  eyesClosed: boolean
  smiling: boolean
}

export function RaccoonBody({
  stage,
  mood,
  emotion,
  offset,
  eyesClosed,
  smiling,
}: RaccoonBodyProps) {
  const geometry = stageGeometry(stage)
  const scale = stageScale[stage]
  const dressed = stage === 'adult' || stage === 'legend'
  const sadEyes = mood === 'sad' || mood === 'hungry'
  const curled = mood === 'sleeping'
  const bodyCy = curled ? geometry.bodyCy + 6 : geometry.bodyCy
  const bodyRx = curled ? geometry.bodyRx + 8 : geometry.bodyRx
  const bodyRy = curled ? geometry.bodyRy - 6 : geometry.bodyRy

  return (
    <g
      className="pet-avatar__body"
      transform={`translate(100 140) scale(${scale}) translate(-100 -140)`}
    >
      <Tail length={geometry.tailLength} />

      <ellipse
        cx={100}
        cy={bodyCy}
        rx={bodyRx}
        ry={bodyRy}
        fill="var(--pet-fur-base)"
      />
      <ellipse
        className="pet-avatar__belly"
        cx={100}
        cy={bodyCy + 4}
        rx={bodyRx - 16}
        ry={bodyRy - 8}
        fill="var(--pet-fur-belly)"
        opacity={0.85}
      />

      <Paws cy={bodyCy + 18} />

      {dressed && <Scarf />}

      <g className="pet-avatar__head">
        <Ears spread={geometry.earSpread} droopy={mood === 'sad'} />
        <circle
          cx={100}
          cy={92}
          r={geometry.headRadius}
          fill="var(--pet-fur-light)"
        />
        <circle
          cx={100}
          cy={92}
          r={geometry.headRadius}
          fill="var(--pet-fur-base)"
          opacity={0.35}
        />
        <Mask smiling={smiling} />
        {emotion === 'eating' && <Tongue />}
        <Eyes offset={offset} closed={eyesClosed} sad={sadEyes} />
        {stage === 'legend' && <Crown />}
      </g>
    </g>
  )
}
