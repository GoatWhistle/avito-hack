import { useCallback, type KeyboardEvent } from 'react'
import { buildAriaLabel, mergeLabels } from './labels'
import { deriveMood, isSmiling } from './mood'
import { useBlink } from './hooks/useBlink'
import { useEyeTracking } from './hooks/useEyeTracking'
import { useIdleSleep } from './hooks/useIdleSleep'
import { useReducedMotion } from './hooks/useReducedMotion'
import { useStroke } from './hooks/useStroke'
import {
  EMOTION_RESET_MS,
  useTransientEmotion,
} from './hooks/useTransientEmotion'
import {
  Aura,
  Confetti,
  Flash,
  GroundShadow,
  Hearts,
  SleepZzz,
  Sparks,
} from './parts/Effects'
import { Egg, HatchingEgg } from './parts/Egg'
import { RaccoonBody } from './parts/RaccoonBody'
import type { PetAvatarProps } from './types'

import './pet-avatar.css'
import './pet-effect-keyframes.css'
import './pet-motion-keyframes.css'
import './pet-states.css'

export function PetAvatar({
  stage,
  emotion = null,
  satiety,
  happiness,
  energy,
  size = 'lg',
  labels,
  idleSleepEnabled = true,
  emotionResetMs = EMOTION_RESET_MS,
  onStroke,
  onEmotionEnd,
  className,
}: PetAvatarProps) {
  const reducedMotion = useReducedMotion()
  const activeEmotion = useTransientEmotion(
    emotion,
    emotionResetMs,
    onEmotionEnd,
  )
  const { isAsleep, isWaving } = useIdleSleep(
    idleSleepEnabled && !reducedMotion,
  )
  const { isPetted, strokeKey, trigger } = useStroke(onStroke)

  const mood = deriveMood({
    satiety,
    happiness,
    energy,
    forcedSleep: isAsleep,
  })
  const sleeping = mood === 'sleeping'
  const hatching = activeEmotion === 'hatching'

  const blinking = useBlink(
    !reducedMotion && !sleeping && activeEmotion === null,
  )
  const { ref, offset } = useEyeTracking<SVGSVGElement>(
    !reducedMotion && !sleeping && !isPetted,
  )

  const handleKeyDown = useCallback(
    (event: KeyboardEvent<SVGSVGElement>) => {
      if (event.key === 'Enter' || event.key === ' ') {
        event.preventDefault()
        trigger()
      }
    },
    [trigger],
  )

  const isLegend = stage === 'legend'
  const showSparks =
    activeEmotion === 'levelup' || activeEmotion === 'celebrate' || isLegend
  const showEgg = stage === 'egg'

  const classes = [
    'pet-avatar',
    `pet-avatar--${size}`,
    `pet-avatar--mood-${mood}`,
    activeEmotion !== null ? `pet-avatar--emotion-${activeEmotion}` : '',
    isPetted ? 'pet-avatar--petted' : '',
    isWaving ? 'pet-avatar--waving' : '',
    reducedMotion ? 'pet-avatar--reduced' : '',
    className ?? '',
  ]
    .filter((value) => value !== '')
    .join(' ')

  return (
    <svg
      ref={ref}
      className={classes}
      viewBox="0 0 200 210"
      role="img"
      tabIndex={0}
      aria-label={buildAriaLabel({
        labels: mergeLabels(labels),
        stage,
        mood,
        emotion: activeEmotion,
      })}
      data-stage={stage}
      data-mood={mood}
      data-emotion={activeEmotion ?? 'none'}
      onClick={trigger}
      onKeyDown={handleKeyDown}
    >
      <GroundShadow />

      {isLegend && <Aura />}
      {showSparks && <Sparks />}

      {showEgg ? (
        hatching ? (
          <HatchingEgg offset={offset} />
        ) : (
          <Egg />
        )
      ) : (
        <RaccoonBody
          stage={stage}
          mood={mood}
          emotion={activeEmotion}
          offset={offset}
          eyesClosed={sleeping || isPetted || blinking}
          smiling={isSmiling(satiety, happiness)}
        />
      )}

      {(activeEmotion === 'levelup' || hatching) && <Flash />}
      {activeEmotion === 'celebrate' && <Confetti />}
      {sleeping && <SleepZzz />}
      {isPetted && <Hearts key={strokeKey} />}
    </svg>
  )
}
