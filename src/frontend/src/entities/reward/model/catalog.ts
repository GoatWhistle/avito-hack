export const REWARD_CONDITIONS = ['level', 'streak'] as const;

export type RewardCondition = (typeof REWARD_CONDITIONS)[number];

export const REWARD_KINDS = ['promocode', 'utility', 'cosmetic'] as const;

export type RewardKind = (typeof REWARD_KINDS)[number];

export const REWARD_IDS = [
  'hatched_badge',
  'free_delivery',
  'streak_freeze',
  'promo_discount',
  'teen_skin',
  'avito_scarf',
  'xp_booster',
  'autoteka_discount',
  'legend_skin',
] as const;

export type RewardId = (typeof REWARD_IDS)[number];

export interface RewardDefinition {
  id: RewardId;
  kind: RewardKind;
  condition: RewardCondition;
  target: number;
}

export interface RewardProgress {
  definition: RewardDefinition;
  current: number;
  percent: number;
  unlocked: boolean;
  remaining: number;
}

export const REWARD_CATALOG: readonly RewardDefinition[] = [
  { id: 'hatched_badge', kind: 'cosmetic', condition: 'level', target: 1 },
  { id: 'free_delivery', kind: 'promocode', condition: 'level', target: 3 },
  { id: 'streak_freeze', kind: 'utility', condition: 'streak', target: 3 },
  { id: 'promo_discount', kind: 'promocode', condition: 'level', target: 5 },
  { id: 'teen_skin', kind: 'cosmetic', condition: 'level', target: 5 },
  { id: 'avito_scarf', kind: 'cosmetic', condition: 'streak', target: 7 },
  { id: 'xp_booster', kind: 'utility', condition: 'streak', target: 14 },
  { id: 'autoteka_discount', kind: 'promocode', condition: 'level', target: 8 },
  { id: 'legend_skin', kind: 'cosmetic', condition: 'level', target: 15 },
] as const;

export function computeRewardProgress(
  definition: RewardDefinition,
  level: number,
  streak: number,
): RewardProgress {
  const current = definition.condition === 'level' ? level : streak;
  const percent = Math.max(0, Math.min(100, Math.round((current / definition.target) * 100)));

  return {
    definition,
    current,
    percent,
    unlocked: current >= definition.target,
    remaining: Math.max(0, definition.target - current),
  };
}

export function buildRewardProgress(level: number, streak: number): RewardProgress[] {
  return REWARD_CATALOG.map((definition) => computeRewardProgress(definition, level, streak));
}
