export const PET_STAGES = ['egg', 'baby', 'teen', 'adult', 'legend'] as const;

export type PetStage = (typeof PET_STAGES)[number];

export const PET_EMOTIONS = [
  'idle',
  'happy',
  'hungry',
  'sad',
  'sleeping',
  'eating',
  'celebrate',
  'levelup',
  'hatching',
] as const;

export type PetEmotion = (typeof PET_EMOTIONS)[number];

export const REACTIVE_EMOTIONS = ['eating', 'celebrate', 'levelup', 'hatching'] as const;

export type ReactiveEmotion = (typeof REACTIVE_EMOTIONS)[number];

export interface PetState {
  id: string;
  userId: string;
  name: string;
  stage: PetStage;
  level: number;
  xp: number;
  nextLevelXp: number;
  satiety: number;
  happiness: number;
  streakDays: number;
  lastCheckInDate: string | null;
  lastDecayTime: string;
  updatedAt: string;
}

export interface Badge {
  id: string;
  name: string;
  description: string;
  iconUrl: string;
  earnedAt: string | null;
}

export interface RaccoonProfile {
  id: string;
  userId: string;
  name: string;
  level: number;
  xp: number;
  xpToNextLevel: number;
  currentStreak: number;
  badges: Badge[];
}

export interface ClaimedReward {
  rewardId: string;
  promocode: string;
}

export interface XpGain {
  amount: number;
  reason: string;
  total: number;
}

export interface LevelUp {
  level: number;
  unlockedRewards: string[];
}

export interface RewardGrant {
  rewardId: string;
  title: string;
}

export interface StreakUpdate {
  days: number;
  milestone: boolean;
}

export const SATIETY_LOW_THRESHOLD = 30;
export const HAPPINESS_LOW_THRESHOLD = 30;
export const HAPPINESS_HIGH_THRESHOLD = 70;
export const PARAM_MAX = 100;

export function isPetStage(value: string): value is PetStage {
  return (PET_STAGES as readonly string[]).includes(value);
}

export function isReactiveEmotion(value: PetEmotion): value is ReactiveEmotion {
  return (REACTIVE_EMOTIONS as readonly PetEmotion[]).includes(value);
}
