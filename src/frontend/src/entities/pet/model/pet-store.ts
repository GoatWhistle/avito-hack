import { createEvent, createStore, sample } from 'effector';

import { isReactiveEmotion, type PetEmotion, type PetState } from './types';

export type ConnectionStatus = 'idle' | 'connecting' | 'open' | 'closed';

export const petStateReceived = createEvent<PetState>();
export const petStateCleared = createEvent();
export const emotionTriggered = createEvent<PetEmotion>();
export const emotionSettled = createEvent();
export const connectionStatusChanged = createEvent<ConnectionStatus>();
export const asleepChanged = createEvent<boolean>();

export const $pet = createStore<PetState | null>(null)
  .on(petStateReceived, (_, pet) => pet)
  .reset(petStateCleared);

export const $reactiveEmotion = createStore<PetEmotion | null>(null)
  .on(emotionTriggered, (_, emotion) => (isReactiveEmotion(emotion) ? emotion : null))
  .reset(emotionSettled, petStateCleared);

export const $connectionStatus = createStore<ConnectionStatus>('idle').on(
  connectionStatusChanged,
  (_, status) => status,
);

export const $isAsleep = createStore(false)
  .on(asleepChanged, (_, value) => value)
  .reset(petStateCleared);

export const $isConnected = $connectionStatus.map((status) => status === 'open');

sample({
  clock: petStateCleared,
  fn: (): ConnectionStatus => 'idle',
  target: connectionStatusChanged,
});
