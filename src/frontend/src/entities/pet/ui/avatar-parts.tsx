import { petPalette } from '@/shared/design';

import type { PetStage } from '../model/types';
import type { EyeOffset } from '../model/use-eye-tracking';

interface EyesProps {
  offset: EyeOffset;
  closed: boolean;
  sad: boolean;
}

const EYE_LEFT_X = 82;
const EYE_RIGHT_X = 118;
const EYE_Y = 92;

export function Eyes({ offset, closed, sad }: EyesProps) {
  const radius = sad ? 7 : 8.5;

  return (
    <g>
      {[EYE_LEFT_X, EYE_RIGHT_X].map((cx) => (
        <g key={cx}>
          <ellipse
            cx={cx}
            cy={EYE_Y}
            rx={radius}
            ry={sad ? 6 : radius}
            fill={petPalette.eyeWhite}
          />
          <circle
            className="pet-avatar__pupil"
            cx={cx}
            cy={EYE_Y}
            r={4.4}
            fill={petPalette.eyePupil}
            style={{ transform: `translate(${offset.x}px, ${offset.y}px)` }}
          />
          <circle cx={cx + 1.8} cy={EYE_Y - 2.4} r={1.5} fill={petPalette.eyeShine} opacity={0.9} />
          {closed ? (
            <path
              d={`M ${cx - radius} ${EYE_Y} q ${radius} ${sad ? -5 : 6} ${radius * 2} 0`}
              stroke={petPalette.maskDark}
              strokeWidth={2.4}
              strokeLinecap="round"
              fill="none"
            />
          ) : (
            <ellipse
              className="pet-avatar__lid"
              cx={cx}
              cy={EYE_Y}
              rx={radius + 0.6}
              ry={radius + 0.6}
              fill={petPalette.maskDark}
            />
          )}
        </g>
      ))}
    </g>
  );
}

interface EarsProps {
  stage: PetStage;
  droopy: boolean;
}

export function Ears({ stage, droopy }: EarsProps) {
  const spread = stage === 'baby' ? 30 : 26;
  const lift = droopy ? 12 : 0;
  const rotate = droopy ? 26 : 0;

  return (
    <g>
      <g transform={`rotate(${-rotate} ${100 - spread} 68)`}>
        <path
          d={`M ${100 - spread - 12} ${72 + lift} q -2 -22 14 -26 q 12 6 8 26 z`}
          fill={petPalette.furBase}
        />
        <path
          d={`M ${100 - spread - 7} ${70 + lift} q -1 -14 9 -17 q 7 4 5 17 z`}
          fill={petPalette.accentPink}
          opacity={0.55}
        />
      </g>
      <g transform={`rotate(${rotate} ${100 + spread} 68)`}>
        <path
          d={`M ${100 + spread + 12} ${72 + lift} q 2 -22 -14 -26 q -12 6 -8 26 z`}
          fill={petPalette.furBase}
        />
        <path
          d={`M ${100 + spread + 7} ${70 + lift} q 1 -14 -9 -17 q -7 4 -5 17 z`}
          fill={petPalette.accentPink}
          opacity={0.55}
        />
      </g>
    </g>
  );
}

export function Mask({ smiling }: { smiling: boolean }) {
  return (
    <g>
      <path
        d="M 66 86 q 8 -16 34 -16 q 26 0 34 16 q -6 20 -34 20 q -28 0 -34 -20 z"
        fill={petPalette.maskDark}
        opacity={0.92}
      />
      <path d="M 74 78 q 26 -12 52 0 q -26 -6 -52 0 z" fill={petPalette.maskLight} opacity={0.5} />
      <ellipse cx={100} cy={108} rx={14} ry={11} fill={petPalette.maskLight} />
      <ellipse cx={100} cy={104} rx={5.6} ry={4.4} fill={petPalette.noseDark} />
      <path
        d={smiling ? 'M 92 112 q 8 8 16 0' : 'M 92 114 q 8 -5 16 0'}
        stroke={petPalette.noseDark}
        strokeWidth={2.2}
        strokeLinecap="round"
        fill="none"
      />
    </g>
  );
}

export function Tail({ stage }: { stage: PetStage }) {
  const length = stage === 'baby' ? 30 : 46;

  return (
    <g className="pet-avatar__tail">
      <path
        d={`M 138 150 q ${length} -6 ${length - 8} -${length}`}
        stroke={petPalette.furBase}
        strokeWidth={20}
        strokeLinecap="round"
        fill="none"
      />
      <path
        d={`M 150 142 q ${length - 14} -4 ${length - 20} -${length - 12}`}
        stroke={petPalette.maskDark}
        strokeWidth={7}
        strokeLinecap="round"
        fill="none"
        opacity={0.65}
      />
    </g>
  );
}

export function Scarf() {
  return (
    <g>
      <path
        d="M 70 128 q 30 14 60 0 q -4 12 -30 12 q -26 0 -30 -12 z"
        fill={petPalette.accentGreen}
      />
      <path d="M 122 134 q 12 8 8 24 q -10 -2 -14 -18 z" fill={petPalette.accentGreen} />
    </g>
  );
}

export function Crown() {
  return (
    <g>
      <path
        d="M 84 58 l 6 -18 l 10 12 l 10 -12 l 6 18 z"
        fill={petPalette.streakCore}
        stroke={petPalette.streakFlame}
        strokeWidth={1.6}
      />
      <circle cx={100} cy={38} r={3} fill={petPalette.accentBlue} />
    </g>
  );
}
