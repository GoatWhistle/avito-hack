export { buildSocketUrl, connectPetSocket } from './pet-socket'
export type {
  PetSocketHandlers,
  PetSocketOptions,
  PetSocketStatus,
} from './pet-socket'
export {
  canCheckInToday,
  clampPercent,
  isSameDay,
  levelProgress,
  MAX_LEVEL,
  nextStep,
  STAT_CRITICAL_BELOW,
  STAT_WARNING_BELOW,
  statTone,
  statViews,
} from './pet-stats'
export type {
  LevelProgress,
  NextStepKey,
  StatKey,
  StatTone,
  StatView,
} from './pet-stats'
export {
  MAX_FREEZES,
  MILESTONE_BONUS_XP,
  STREAK_BONUS_DAYS,
  STREAK_BONUS_FACTOR,
  STREAK_MILESTONES,
  STREAK_RESET_HOURS,
  daysToBonus,
  daysToNextMilestone,
  isBonusActive,
  nextMilestone,
  streakRules,
} from './streak-rules'
export type {
  StreakRuleKey,
  StreakRuleTone,
  StreakRuleView,
} from './streak-rules'
export { buildAvatarLabels } from './pet-avatar-labels'
export { adviceText, summaryMessage } from './summary-text'
export { petSpeech } from './pet-speech'
export type { PetSpeechInput, PetSpeechKey } from './pet-speech'
