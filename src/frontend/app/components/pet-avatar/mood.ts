import type { PetMood, PetStage } from './types'

export const SATIETY_HUNGRY_BELOW = 30
export const HAPPINESS_SAD_BELOW = 30
export const ENERGY_SLEEPING_BELOW = 20
export const HAPPINESS_HAPPY_FROM = 70
export const SMILE_HAPPINESS_FROM = 50
export const SMILE_SATIETY_FROM = 30

export interface MoodInput {
  satiety: number
  happiness: number
  energy: number
  forcedSleep?: boolean
}

export function deriveMood({
  satiety,
  happiness,
  energy,
  forcedSleep = false,
}: MoodInput): PetMood {
  if (forcedSleep || energy < ENERGY_SLEEPING_BELOW) {
    return 'sleeping'
  }

  if (satiety < SATIETY_HUNGRY_BELOW) {
    return 'hungry'
  }

  if (happiness < HAPPINESS_SAD_BELOW) {
    return 'sad'
  }

  if (happiness >= HAPPINESS_HAPPY_FROM && satiety >= SATIETY_HUNGRY_BELOW) {
    return 'happy'
  }

  return 'idle'
}

export function isSmiling(satiety: number, happiness: number): boolean {
  return happiness >= SMILE_HAPPINESS_FROM && satiety >= SMILE_SATIETY_FROM
}

export const stageScale: Record<PetStage, number> = {
  egg: 1,
  baby: 1,
  teen: 0.94,
  adult: 0.9,
  legend: 0.9,
}

export interface StageGeometry {
  headRadius: number
  bodyRx: number
  bodyRy: number
  bodyCy: number
  earSpread: number
  tailLength: number
}

const GEOMETRY: Record<PetStage, StageGeometry> = {
  egg: {
    headRadius: 48,
    bodyRx: 44,
    bodyRy: 34,
    bodyCy: 156,
    earSpread: 30,
    tailLength: 30,
  },
  baby: {
    headRadius: 48,
    bodyRx: 40,
    bodyRy: 32,
    bodyCy: 158,
    earSpread: 30,
    tailLength: 32,
  },
  teen: {
    headRadius: 42,
    bodyRx: 44,
    bodyRy: 44,
    bodyCy: 154,
    earSpread: 27,
    tailLength: 42,
  },
  adult: {
    headRadius: 38,
    bodyRx: 47,
    bodyRy: 52,
    bodyCy: 150,
    earSpread: 25,
    tailLength: 48,
  },
  legend: {
    headRadius: 38,
    bodyRx: 47,
    bodyRy: 52,
    bodyCy: 150,
    earSpread: 25,
    tailLength: 48,
  },
}

export function stageGeometry(stage: PetStage): StageGeometry {
  return GEOMETRY[stage]
}
