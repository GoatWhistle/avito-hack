import type { ReactNode } from 'react'
import type { PetEmotion } from '#/components/pet-avatar'
import type { PetSpeechKey } from '#/features/pet/lib'
import type { XpToast } from '#/features/pet/hooks'
import type { Pet } from '#/features/pet/types'
import { PetHero } from './PetHero'
import { PetSpeech } from './PetSpeech'

export interface PetStageProps {
  pet: Pet
  speech: PetSpeechKey
  emotion: PetEmotion | null
  xpToasts: XpToast[]
  onStroke: () => void
  onEmotionEnd: () => void
  children: ReactNode
}

export function PetStage({
  pet,
  speech,
  emotion,
  xpToasts,
  onStroke,
  onEmotionEnd,
  children,
}: PetStageProps) {
  return (
    <div className="flex flex-col items-center gap-5">
      <PetHero
        pet={pet}
        emotion={emotion}
        xpToasts={xpToasts}
        onStroke={onStroke}
        onEmotionEnd={onEmotionEnd}
      />

      <PetSpeech speech={speech} petName={pet.name} />

      <div className="w-full max-w-sm">{children}</div>
    </div>
  )
}
