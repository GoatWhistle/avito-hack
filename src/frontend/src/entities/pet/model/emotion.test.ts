import { describe, expect, it } from 'vitest';

import { isMaxLevel, resolveBackgroundEmotion, xpProgressPercent, xpRemaining } from './emotion';
import type { PetState } from './types';

function makePet(overrides: Partial<PetState> = {}): PetState {
  return {
    id: 'pet-1',
    userId: 'user-1',
    name: 'Noti',
    stage: 'baby',
    level: 2,
    xp: 40,
    nextLevelXp: 100,
    satiety: 80,
    happiness: 80,
    streakDays: 3,
    lastCheckInDate: null,
    lastDecayTime: '2026-01-01T00:00:00Z',
    updatedAt: '2026-01-01T00:00:00Z',
    ...overrides,
  };
}

describe('resolveBackgroundEmotion', () => {
  it('falls back to idle without a pet', () => {
    expect(resolveBackgroundEmotion(null, false)).toBe('idle');
  });

  it('sleeps whenever the user is idle, whatever the stats are', () => {
    expect(resolveBackgroundEmotion(makePet({ satiety: 5 }), true)).toBe('sleeping');
  });

  it('prefers hunger over sadness', () => {
    expect(resolveBackgroundEmotion(makePet({ satiety: 10, happiness: 10 }), false)).toBe('hungry');
  });

  it('reports sadness when only happiness is low', () => {
    expect(resolveBackgroundEmotion(makePet({ satiety: 90, happiness: 10 }), false)).toBe('sad');
  });

  it('is happy when both stats are comfortable', () => {
    expect(resolveBackgroundEmotion(makePet({ satiety: 90, happiness: 90 }), false)).toBe('happy');
  });

  it('stays idle in the middle band', () => {
    expect(resolveBackgroundEmotion(makePet({ satiety: 40, happiness: 40 }), false)).toBe('idle');
  });
});

describe('xp helpers', () => {
  it('computes the percentage towards the next level', () => {
    expect(xpProgressPercent(makePet({ xp: 40, nextLevelXp: 100 }))).toBe(40);
  });

  it('clamps the percentage to a hundred', () => {
    expect(xpProgressPercent(makePet({ xp: 500, nextLevelXp: 100 }))).toBe(100);
  });

  it('treats a zero threshold as the max level', () => {
    const maxed = makePet({ nextLevelXp: 0 });

    expect(isMaxLevel(maxed)).toBe(true);
    expect(xpProgressPercent(maxed)).toBe(100);
    expect(xpRemaining(maxed)).toBe(0);
  });

  it('reports the remaining xp', () => {
    expect(xpRemaining(makePet({ xp: 40, nextLevelXp: 100 }))).toBe(60);
  });
});
