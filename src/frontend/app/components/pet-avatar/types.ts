export type PetStage = 'baby' | 'teen' | 'adult' | 'legend'

export type PetEmotion = 'eating' | 'celebrate' | 'levelup'

export type PetMood = 'idle' | 'happy' | 'hungry' | 'sad' | 'sleeping'

export type PetAvatarSize = 'sm' | 'md' | 'lg' | 'xl'

export interface PetAvatarLabels {
  name: string
  stage: Record<PetStage, string>
  mood: Record<PetMood, string>
  emotion: Record<PetEmotion, string>
  strokeHint: string
  loading: string
}

export interface PetAvatarProps {
  stage: PetStage
  emotion?: PetEmotion | null
  satiety: number
  happiness: number
  energy: number
  size?: PetAvatarSize
  labels?: Partial<PetAvatarLabels>
  emotionResetMs?: number
  onStroke?: () => void
  onEmotionEnd?: () => void
  className?: string
}
