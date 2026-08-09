import { describe, expect, it } from 'vitest'
import { deriveMood } from './mood'
import {
  isLegendStage,
  lottieSourceFor,
  RACCOON_ADULT_SRC,
  RACCOON_TEEN_SRC,
} from './lottie-sources'
import type { PetStage } from './types'

const FULL = { satiety: 90, happiness: 90, energy: 90 }

describe('deriveMood', () => {
  it('returns happy for healthy stats', () => {
    expect(deriveMood(FULL)).toBe('happy')
  })

  it('returns idle for mid stats', () => {
    expect(deriveMood({ satiety: 60, happiness: 50, energy: 60 })).toBe('idle')
  })

  it('prioritises sleeping when energy is low', () => {
    expect(deriveMood({ satiety: 5, happiness: 5, energy: 10 })).toBe(
      'sleeping',
    )
  })

  it('returns sleeping when forced regardless of energy', () => {
    expect(deriveMood({ ...FULL, forcedSleep: true })).toBe('sleeping')
  })

  it('returns hungry below the satiety threshold', () => {
    expect(deriveMood({ satiety: 20, happiness: 90, energy: 90 })).toBe(
      'hungry',
    )
  })

  it('returns sad below the happiness threshold', () => {
    expect(deriveMood({ satiety: 90, happiness: 20, energy: 90 })).toBe('sad')
  })

  it('prefers hungry over sad when both are low', () => {
    expect(deriveMood({ satiety: 10, happiness: 10, energy: 90 })).toBe(
      'hungry',
    )
  })

  it.each([
    [30, 'idle'],
    [29, 'hungry'],
  ] as const)('handles satiety boundary %i', (satiety, expected) => {
    expect(deriveMood({ satiety, happiness: 50, energy: 90 })).toBe(expected)
  })
})

describe('lottie source mapping', () => {
  const stages: PetStage[] = ['baby', 'teen', 'adult', 'legend']

  it('maps every stage to a designer animation', () => {
    for (const stage of stages) {
      expect(lottieSourceFor(stage)).toMatch(/^\/lottie\/raccoon-.+\.json$/)
    }
  })

  it('uses the young raccoon for the early stages', () => {
    for (const stage of ['baby', 'teen'] as PetStage[]) {
      expect(lottieSourceFor(stage)).toBe(RACCOON_TEEN_SRC)
    }
  })

  it('uses the grown raccoon for the late stages', () => {
    for (const stage of ['adult', 'legend'] as PetStage[]) {
      expect(lottieSourceFor(stage)).toBe(RACCOON_ADULT_SRC)
    }
  })

  it('marks only the legend stage as regal', () => {
    expect(isLegendStage('legend')).toBe(true)
    expect(isLegendStage('adult')).toBe(false)
  })
})
