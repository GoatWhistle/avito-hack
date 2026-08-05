import { petPalette, petStageScale } from '@/shared/design';

import type { PetEmotion, PetStage } from '../model/types';
import type { EyeOffset } from '../model/use-eye-tracking';

import { bodyGeometry } from './avatar-geometry';
import { Crown, Ears, Eyes, Mask, Scarf, Tail } from './avatar-parts';

export interface AvatarBodyProps {
  stage: PetStage;
  emotion: PetEmotion;
  offset: EyeOffset;
  eyesClosed: boolean;
  smiling: boolean;
}

export function AvatarBody({ stage, emotion, offset, eyesClosed, smiling }: AvatarBodyProps) {
  const geometry = bodyGeometry(stage);
  const scale = petStageScale[stage];
  const dressed = stage === 'adult' || stage === 'legend';
  const sadEyes = emotion === 'sad' || emotion === 'hungry';

  return (
    <g
      className="pet-avatar__body"
      transform={`translate(100 140) scale(${scale}) translate(-100 -140)`}
    >
      <Tail stage={stage} />

      <ellipse
        cx={100}
        cy={geometry.bodyCy}
        rx={44}
        ry={geometry.bodyRy}
        fill={petPalette.furBase}
      />
      <ellipse
        cx={100}
        cy={geometry.bodyCy + 4}
        rx={28}
        ry={geometry.bodyRy - 8}
        fill={petPalette.furBelly}
        opacity={0.85}
      />
      <ellipse cx={66} cy={geometry.bodyCy + 18} rx={13} ry={9} fill={petPalette.furDark} />
      <ellipse cx={134} cy={geometry.bodyCy + 18} rx={13} ry={9} fill={petPalette.furDark} />

      {dressed && <Scarf />}

      <g className="pet-avatar__head">
        <Ears stage={stage} droopy={emotion === 'sad'} />
        <circle cx={100} cy={92} r={geometry.headRadius} fill={petPalette.furLight} />
        <circle cx={100} cy={92} r={geometry.headRadius} fill={petPalette.furBase} opacity={0.35} />
        <Mask smiling={smiling} />
        <Eyes offset={offset} closed={eyesClosed} sad={sadEyes} />
        {stage === 'legend' && <Crown />}
      </g>
    </g>
  );
}
