export { petApi } from './api/pet-api';
export { petKeys } from './api/query-keys';
export { useBadges, useClaimReward, usePetProfile } from './api/use-pet';
export { useEmotionDecay, usePetSocket } from './api/use-pet-socket';
export type { PetSocketApi } from './api/use-pet-socket';
export { createPetSocket, handlePetMessage, PET_MESSAGE_TYPES } from './api/ws-bridge';
export {
  isMaxLevel,
  resolveBackgroundEmotion,
  xpProgressPercent,
  xpRemaining,
} from './model/emotion';
export { toBadge, toClaimedReward, toPetState, toRaccoonProfile } from './model/mappers';
export {
  $connectionStatus,
  $isAsleep,
  $isConnected,
  $pet,
  $reactiveEmotion,
  asleepChanged,
  connectionStatusChanged,
  emotionSettled,
  emotionTriggered,
  petStateCleared,
  petStateReceived,
} from './model/pet-store';
export type { ConnectionStatus } from './model/pet-store';
export {
  badgeDtoSchema,
  claimRewardDtoSchema,
  petStateDtoSchema,
  raccoonProfileDtoSchema,
} from './model/schemas';
export type { BadgeDto, ClaimRewardDto, PetStateDto, RaccoonProfileDto } from './model/schemas';
export { isPetStage, isReactiveEmotion, PET_EMOTIONS, PET_STAGES } from './model/types';
export type {
  Badge,
  ClaimedReward,
  PetEmotion,
  PetStage,
  PetState,
  RaccoonProfile,
} from './model/types';
export { useEyeTracking, useReducedMotion } from './model/use-eye-tracking';
export { useIdleSleep } from './model/use-idle-sleep';
export { PetAvatar } from './ui/PetAvatar';
export type { PetAvatarProps } from './ui/PetAvatar';
export { PetStatBar } from './ui/PetStatBar';
export { StreakFlame } from './ui/StreakFlame';
export { XpBar } from './ui/XpBar';
