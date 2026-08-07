import type { PetAvatarLabels, PetEmotion, PetMood, PetStage } from './types'

export const defaultLabels: PetAvatarLabels = {
  name: 'Ноти',
  stage: {
    egg: 'яйцо',
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
    hatching: 'вылупляется',
  },
  strokeHint: 'нажмите, чтобы погладить',
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
