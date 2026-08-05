import type { BadgeDto, ClaimRewardDto, PetStateDto, RaccoonProfileDto } from './schemas';
import type { Badge, ClaimedReward, PetState, RaccoonProfile } from './types';

export function toPetState(dto: PetStateDto): PetState {
  return {
    id: dto.id,
    userId: dto.user_id,
    name: dto.name,
    stage: dto.stage,
    level: dto.level,
    xp: dto.xp,
    nextLevelXp: dto.next_level_xp,
    satiety: dto.satiety,
    happiness: dto.happiness,
    streakDays: dto.streak_days,
    lastCheckInDate: dto.last_checkin_date ?? null,
    lastDecayTime: dto.last_decay_time,
    updatedAt: dto.updated_at,
  };
}

export function toBadge(dto: BadgeDto): Badge {
  return {
    id: dto.id,
    name: dto.name,
    description: dto.description,
    iconUrl: dto.icon_url,
    earnedAt: dto.earned_at ?? null,
  };
}

export function toRaccoonProfile(dto: RaccoonProfileDto): RaccoonProfile {
  return {
    id: dto.id,
    userId: dto.user_id,
    name: dto.name,
    level: dto.level,
    xp: dto.xp,
    xpToNextLevel: dto.xp_to_next_level,
    currentStreak: dto.current_streak,
    badges: dto.badges.map(toBadge),
  };
}

export function toClaimedReward(dto: ClaimRewardDto): ClaimedReward {
  return { rewardId: dto.reward_id, promocode: dto.promocode };
}
