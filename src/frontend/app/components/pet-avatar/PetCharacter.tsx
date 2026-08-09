import { buildAriaLabel, mergeLabels } from './labels'
import { isLegendStage, lottieSourceFor } from './lottie-sources'
import { deriveMood } from './mood'
import {
  EMOTION_RESET_MS,
  useTransientEmotion,
} from './hooks/useTransientEmotion'
import { PetLottie } from './PetLottie'
import type { PetAvatarProps } from './types'

export type PetCharacterProps = PetAvatarProps

export function PetCharacter({
  stage,
  emotion = null,
  satiety,
  happiness,
  energy,
  size = 'lg',
  labels,
  emotionResetMs = EMOTION_RESET_MS,
  onStroke,
  onEmotionEnd,
  className,
}: PetCharacterProps) {
  const activeEmotion = useTransientEmotion(
    emotion,
    emotionResetMs,
    onEmotionEnd,
  )
  const mood = deriveMood({ satiety, happiness, energy })
  const merged = mergeLabels(labels)

  return (
    <PetLottie
      src={lottieSourceFor(stage)}
      ariaLabel={buildAriaLabel({
        labels: merged,
        stage,
        mood,
        emotion: activeEmotion,
      })}
      placeholderLabel={merged.loading}
      sizeClass={`pet-lottie--${size}`}
      stage={stage}
      mood={mood}
      emotion={activeEmotion ?? 'none'}
      legend={isLegendStage(stage)}
      className={className}
      onStroke={onStroke}
    />
  )
}
