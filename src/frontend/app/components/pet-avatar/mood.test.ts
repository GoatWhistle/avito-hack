import { describe, expect, it } from 'vitest'
import { deriveMood, isSmiling, stageGeometry, stageScale } from './mood'
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

describe('isSmiling', () => {
  it('smiles when both stats are comfortable', () => {
    expect(isSmiling(60, 60)).toBe(true)
  })

  it('does not smile when happiness is low', () => {
    expect(isSmiling(90, 20)).toBe(false)
  })

  it('does not smile when satiety is low', () => {
    expect(isSmiling(10, 90)).toBe(false)
  })
})

describe('stage geometry', () => {
  const stages: PetStage[] = ['egg', 'baby', 'teen', 'adult', 'legend']

  it('exposes geometry for every stage', () => {
    for (const stage of stages) {
      expect(stageGeometry(stage).headRadius).toBeGreaterThan(0)
      expect(stageScale[stage]).toBeGreaterThan(0)
    }
  })

  it('shrinks the head as the pet grows up', () => {
    expect(stageGeometry('baby').headRadius).toBeGreaterThan(
      stageGeometry('adult').headRadius,
    )
  })

  it('grows the body as the pet grows up', () => {
    expect(stageGeometry('adult').bodyRy).toBeGreaterThan(
      stageGeometry('baby').bodyRy,
    )
  })
})
