import { describe, expect, it } from 'vitest'
import type { Pet } from '#/features/pet/types'
import {
  MAX_FREEZES,
  STREAK_BONUS_DAYS,
  STREAK_BONUS_FACTOR,
  daysToBonus,
  daysToNextMilestone,
  isBonusActive,
  nextMilestone,
  streakRules,
} from './streak-rules'

const makePet = (overrides: Partial<Pet> = {}): Pet =>
  ({
    id: 'pet-1',
    user_id: 'user-1',
    name: 'Енот',
    stage: 'baby',
    state: 'happy',
    level: 1,
    xp: 0,
    next_level_xp: 5,
    satiety: 70,
    happiness: 70,
    energy: 100,
    streak_days: 0,
    freezes: 0,
    is_hatched: true,
    ...overrides,
  }) as Pet

describe('streak rule constants match the backend domain', () => {
  it('uses the backend bonus threshold and factor', () => {
    expect(STREAK_BONUS_DAYS).toBe(7)
    expect(STREAK_BONUS_FACTOR).toBe(1.5)
    expect(MAX_FREEZES).toBe(3)
  })
})

describe('isBonusActive', () => {
  it('is inactive below the threshold', () => {
    expect(isBonusActive(6)).toBe(false)
  })

  it('activates exactly at the threshold', () => {
    expect(isBonusActive(7)).toBe(true)
    expect(isBonusActive(41)).toBe(true)
  })
})

describe('daysToBonus', () => {
  it('counts down to the threshold', () => {
    expect(daysToBonus(0)).toBe(7)
    expect(daysToBonus(5)).toBe(2)
  })

  it('never goes negative', () => {
    expect(daysToBonus(41)).toBe(0)
  })
})

describe('nextMilestone', () => {
  it('returns the next backend milestone', () => {
    expect(nextMilestone(0)).toBe(3)
    expect(nextMilestone(3)).toBe(7)
    expect(nextMilestone(20)).toBe(30)
  })

  it('returns null past the last milestone', () => {
    expect(nextMilestone(30)).toBeNull()
    expect(daysToNextMilestone(30)).toBeNull()
  })

  it('reports the gap to the next milestone', () => {
    expect(daysToNextMilestone(5)).toBe(2)
  })
})

describe('streakRules', () => {
  it('marks the bonus pending with days remaining for a fresh streak', () => {
    const [bonus] = streakRules(makePet({ streak_days: 2 }))

    expect(bonus.tone).toBe('pending')
    expect(bonus.badgeCount).toBe(5)
  })

  it('marks the bonus active once the streak reaches the threshold', () => {
    const [bonus] = streakRules(makePet({ streak_days: 41 }))

    expect(bonus.tone).toBe('active')
    expect(bonus.badgeCount).toBe(41)
  })

  it('reports freezes as unavailable when none are held', () => {
    const rules = streakRules(makePet({ freezes: 0 }))
    const freeze = rules.find((rule) => rule.key === 'freeze')

    expect(freeze?.tone).toBe('muted')
    expect(freeze?.badgeCount).toBe(0)
  })

  it('reports held freezes as active', () => {
    const rules = streakRules(makePet({ freezes: 2 }))
    const freeze = rules.find((rule) => rule.key === 'freeze')

    expect(freeze?.tone).toBe('active')
    expect(freeze?.badgeCount).toBe(2)
  })

  it('always returns the three documented rules', () => {
    expect(streakRules(makePet()).map((rule) => rule.key)).toEqual([
      'bonus',
      'reset',
      'freeze',
    ])
  })
})
