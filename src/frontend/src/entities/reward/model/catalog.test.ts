import { describe, expect, it } from 'vitest';

import { buildRewardProgress, computeRewardProgress, REWARD_CATALOG } from './catalog';

const levelReward = {
  id: 'free_delivery',
  kind: 'promocode',
  condition: 'level',
  target: 3,
} as const;
const streakReward = {
  id: 'avito_scarf',
  kind: 'cosmetic',
  condition: 'streak',
  target: 7,
} as const;

describe('computeRewardProgress', () => {
  it('locks a reward whose condition is not met yet', () => {
    const progress = computeRewardProgress(levelReward, 1, 0);

    expect(progress.unlocked).toBe(false);
    expect(progress.remaining).toBe(2);
    expect(progress.percent).toBe(33);
  });

  it('unlocks the reward exactly on the target', () => {
    const progress = computeRewardProgress(levelReward, 3, 0);

    expect(progress.unlocked).toBe(true);
    expect(progress.remaining).toBe(0);
    expect(progress.percent).toBe(100);
  });

  it('never reports more than a hundred percent', () => {
    expect(computeRewardProgress(levelReward, 99, 0).percent).toBe(100);
  });

  it('reads the streak for streak conditions', () => {
    const progress = computeRewardProgress(streakReward, 15, 5);

    expect(progress.current).toBe(5);
    expect(progress.unlocked).toBe(false);
    expect(progress.remaining).toBe(2);
  });
});

describe('buildRewardProgress', () => {
  it('covers the whole catalog', () => {
    expect(buildRewardProgress(1, 0)).toHaveLength(REWARD_CATALOG.length);
  });

  it('always shows the user what is left to do', () => {
    for (const item of buildRewardProgress(2, 1)) {
      expect(item.unlocked || item.remaining > 0).toBe(true);
    }
  });
});
