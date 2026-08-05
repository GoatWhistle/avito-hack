import { useCallback, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { petPalette } from '@/shared/design';

import type { PetEmotion, PetStage } from '../model/types';
import { useEyeTracking, useReducedMotion } from '../model/use-eye-tracking';

import { Aura, Egg, Hearts, SleepZzz, Sparks } from './avatar-effects';
import { AvatarBody } from './AvatarBody';
import './pet-avatar.css';
import './pet-keyframes.css';

const PET_FEEDBACK_MS = 1_200;
const SMILE_HAPPINESS = 50;
const SMILE_SATIETY = 30;

export interface PetAvatarProps {
  stage: PetStage;
  emotion: PetEmotion;
  satiety: number;
  happiness: number;
  onStroke?: () => void;
}

function usePettedFlag(): [boolean, () => void] {
  const [isPetted, setIsPetted] = useState(false);

  useEffect(() => {
    if (!isPetted) {
      return;
    }

    const timer = setTimeout(() => {
      setIsPetted(false);
    }, PET_FEEDBACK_MS);

    return () => {
      clearTimeout(timer);
    };
  }, [isPetted]);

  const trigger = useCallback(() => {
    setIsPetted(true);
  }, []);

  return [isPetted, trigger];
}

export function PetAvatar({ stage, emotion, satiety, happiness, onStroke }: PetAvatarProps) {
  const { t } = useTranslation('pet');
  const reducedMotion = useReducedMotion();
  const [isPetted, markPetted] = usePettedFlag();
  const isSleeping = emotion === 'sleeping';
  const { ref, offset } = useEyeTracking<SVGSVGElement>(!reducedMotion && !isSleeping && !isPetted);

  const handleStroke = useCallback(() => {
    markPetted();
    onStroke?.();
  }, [markPetted, onStroke]);

  const handleKeyDown = useCallback(
    (event: React.KeyboardEvent<SVGSVGElement>) => {
      if (event.key === 'Enter' || event.key === ' ') {
        event.preventDefault();
        handleStroke();
      }
    },
    [handleStroke],
  );

  const isLegend = stage === 'legend';
  const showSparks = emotion === 'levelup' || emotion === 'celebrate' || isLegend;
  const classes = ['pet-avatar', `pet-avatar--${emotion}`, isPetted ? 'pet-avatar--petted' : '']
    .filter((value) => value !== '')
    .join(' ');

  return (
    <svg
      ref={ref}
      className={classes}
      viewBox="0 0 200 210"
      role="img"
      tabIndex={0}
      aria-label={t('avatar.label', {
        stage: t(`stage.${stage}`),
        emotion: t(`emotion.${emotion}`),
      })}
      onClick={handleStroke}
      onKeyDown={handleKeyDown}
    >
      <ellipse cx={100} cy={196} rx={58} ry={9} fill={petPalette.furDark} opacity={0.16} />

      {isLegend && <Aura />}
      {showSparks && <Sparks />}

      {stage === 'egg' ? (
        <Egg hatching={emotion === 'hatching'} />
      ) : (
        <AvatarBody
          stage={stage}
          emotion={emotion}
          offset={offset}
          eyesClosed={isSleeping || isPetted}
          smiling={happiness >= SMILE_HAPPINESS && satiety >= SMILE_SATIETY}
        />
      )}

      {isSleeping && <SleepZzz />}
      {isPetted && <Hearts />}
    </svg>
  );
}
