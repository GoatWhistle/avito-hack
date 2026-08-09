import { describe, expect, it } from 'vitest'
import { petSpeech } from './pet-speech'
import type { Pet } from '#/features/pet/types'

const basePet: Pet = {
  id: 'pet-1',
  user_id: 'user-1',
  name: 'Ноти',
  stage: 'baby',
  state: 'neutral',
  level: 3,
  xp: 15,
  next_level_xp: 22,
  satiety: 80,
  happiness: 80,
  energy: 80,
  streak_days: 2,
  freezes: 1,
  is_hatched: true,
  hatched_at: null,
  last_checkin_date: null,
}

const speechOf = (overrides: Partial<Pet>, canCheckIn = false) =>
  petSpeech({ pet: { ...basePet, ...overrides }, canCheckIn })

describe('petSpeech', () => {
  it('complains about hunger before anything else', () => {
    expect(speechOf({ satiety: 10, energy: 5, happiness: 5 }, true)).toBe(
      'hungry',
    )
  })

  it('asks for rest when energy is critical', () => {
    expect(speechOf({ energy: 10, happiness: 5 }, true)).toBe('tired')
  })

  it('asks for a stroke when the mood drops', () => {
    expect(speechOf({ happiness: 40 }, true)).toBe('lonely')
  })

  it('nudges the check in when stats are fine', () => {
    expect(speechOf({}, true)).toBe('checkIn')
  })

  it('celebrates a long streak once the check in is done', () => {
    expect(speechOf({ streak_days: 9 })).toBe('streak')
  })

  it('celebrates the max level when there is no long streak', () => {
    expect(speechOf({ level: 15, streak_days: 1 })).toBe('legend')
  })

  it('falls back to a happy line', () => {
    expect(speechOf({})).toBe('happy')
  })
})
