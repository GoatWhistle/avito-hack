import type { PetAvatarLabels, PetEmotion, PetMood, PetStage } from './types'

export const defaultLabels: PetAvatarLabels = {
  name: 'Ноти',
  stage: {
    baby: 'малыш',
    teen: 'подросток',
    adult: 'взрослый',
    legend: 'легенда',
  },
  mood: {
    idle: 'спокоен',
    happy: 'счастлив',
    hungry: 'голоден',
    sad: 'грустит',
    sleeping: 'спит',
  },
  emotion: {
    eating: 'ест',
    celebrate: 'празднует',
    levelup: 'получил новый уровень',
  },
  strokeHint: 'нажмите, чтобы погладить',
  loading: 'Загружаем питомца',
}

export function mergeLabels(
  overrides?: Partial<PetAvatarLabels>,
): PetAvatarLabels {
  if (overrides === undefined) {
    return defaultLabels
  }

  return {
    name: overrides.name ?? defaultLabels.name,
    stage: { ...defaultLabels.stage, ...overrides.stage },
    mood: { ...defaultLabels.mood, ...overrides.mood },
    emotion: { ...defaultLabels.emotion, ...overrides.emotion },
    strokeHint: overrides.strokeHint ?? defaultLabels.strokeHint,
    loading: overrides.loading ?? defaultLabels.loading,
  }
}

export interface AriaLabelInput {
  labels: PetAvatarLabels
  stage: PetStage
  mood: PetMood
  emotion: PetEmotion | null
}

export function buildAriaLabel({
  labels,
  stage,
  mood,
  emotion,
}: AriaLabelInput): string {
  const state = emotion === null ? labels.mood[mood] : labels.emotion[emotion]

  return `Енот ${labels.name}, ${labels.stage[stage]}, ${state}. ${labels.strokeHint}`
}
