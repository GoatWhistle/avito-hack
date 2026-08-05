import { beforeEach, describe, expect, it } from 'vitest';

import { $pet, $reactiveEmotion, petStateCleared } from '../model/pet-store';

import { handlePetMessage, PET_MESSAGE_TYPES } from './ws-bridge';

const payload = {
  id: 'pet-1',
  user_id: 'user-1',
  name: 'Noti',
  stage: 'teen',
  level: 4,
  xp: 120,
  next_level_xp: 300,
  satiety: 55,
  happiness: 71,
  streak_days: 6,
  last_decay_time: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-01T00:00:00Z',
};

describe('handlePetMessage', () => {
  beforeEach(() => {
    petStateCleared();
  });

  it('stores the state coming from pet.state', () => {
    handlePetMessage({ type: PET_MESSAGE_TYPES.state, payload });

    expect($pet.getState()?.level).toBe(4);
    expect($pet.getState()?.stage).toBe('teen');
    expect($reactiveEmotion.getState()).toBeNull();
  });

  it('accepts the server push pet.updated', () => {
    handlePetMessage({ type: PET_MESSAGE_TYPES.updated, payload });

    expect($pet.getState()?.happiness).toBe(71);
  });

  it('triggers the level up emotion', () => {
    handlePetMessage({ type: PET_MESSAGE_TYPES.levelUp, payload });

    expect($reactiveEmotion.getState()).toBe('levelup');
  });

  it('triggers the eating emotion on gained xp', () => {
    handlePetMessage({ type: PET_MESSAGE_TYPES.xpGained, payload });

    expect($reactiveEmotion.getState()).toBe('eating');
  });

  it('triggers hatching', () => {
    handlePetMessage({ type: PET_MESSAGE_TYPES.hatched, payload });

    expect($reactiveEmotion.getState()).toBe('hatching');
  });

  it('ignores malformed payloads without throwing', () => {
    handlePetMessage({ type: PET_MESSAGE_TYPES.state, payload: { nope: true } });

    expect($pet.getState()).toBeNull();
  });

  it('ignores unknown message types', () => {
    handlePetMessage({ type: PET_MESSAGE_TYPES.pong });

    expect($pet.getState()).toBeNull();
    expect($reactiveEmotion.getState()).toBeNull();
  });
});
