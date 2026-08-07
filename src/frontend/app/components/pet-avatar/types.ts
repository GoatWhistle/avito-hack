export type PetStage = 'egg' | 'baby' | 'teen' | 'adult' | 'legend'

export type PetEmotion = 'eating' | 'celebrate' | 'levelup' | 'hatching'

export type PetMood = 'idle' | 'happy' | 'hungry' | 'sad' | 'sleeping'

export type PetAvatarSize = 'sm' | 'md' | 'lg' | 'xl'

export interface EyeOffset {
  x: number
  y: number
}

export interface PetAvatarLabels {
  name: string
  stage: Record<PetStage, string>
  mood: Record<PetMood, string>
  emotion: Record<PetEmotion, string>
  strokeHint: string
}

export interface PetAvatarProps {
  stage: PetStage
  emotion?: PetEmotion | null
  satiety: number
  happiness: number
  energy: number
  size?: PetAvatarSize
  labels?: Partial<PetAvatarLabels>
  idleSleepEnabled?: boolean
  emotionResetMs?: number
  onStroke?: () => void
  onEmotionEnd?: () => void
  className?: string
}
