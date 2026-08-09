import type { PetMood } from './types'

export const SATIETY_HUNGRY_BELOW = 30
export const HAPPINESS_SAD_BELOW = 30
export const ENERGY_SLEEPING_BELOW = 20
export const HAPPINESS_HAPPY_FROM = 70

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
