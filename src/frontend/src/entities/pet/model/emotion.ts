import {
  HAPPINESS_LOW_THRESHOLD,
  SATIETY_LOW_THRESHOLD,
  type PetEmotion,
  type PetState,
} from './types';

export function resolveBackgroundEmotion(pet: PetState | null, isAsleep: boolean): PetEmotion {
  if (pet === null) {
    return 'idle';
  }

  if (isAsleep) {
    return 'sleeping';
  }

  if (pet.satiety < SATIETY_LOW_THRESHOLD) {
    return 'hungry';
  }

  if (pet.happiness < HAPPINESS_LOW_THRESHOLD) {
    return 'sad';
  }

  if (pet.happiness >= HAPPINESS_LOW_THRESHOLD * 2 && pet.satiety >= SATIETY_LOW_THRESHOLD * 2) {
    return 'happy';
  }

  return 'idle';
}

export function xpProgressPercent(pet: PetState | null): number {
  if (pet === null || pet.nextLevelXp <= 0) {
    return 100;
  }

  const ratio = (pet.xp / pet.nextLevelXp) * 100;

  return Math.max(0, Math.min(100, Math.round(ratio)));
}

export function xpRemaining(pet: PetState | null): number {
  if (pet === null || pet.nextLevelXp <= 0) {
    return 0;
  }

  return Math.max(0, pet.nextLevelXp - pet.xp);
}

export function isMaxLevel(pet: PetState | null): boolean {
  return pet !== null && pet.nextLevelXp <= 0;
}
