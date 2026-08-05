import { createWsClient, type WsClient, type WsMessage, type WsStatus } from '@/shared/api/ws';

import { toPetState } from '../model/mappers';
import {
  connectionStatusChanged,
  emotionTriggered,
  petStateReceived,
  type ConnectionStatus,
} from '../model/pet-store';
import { petStateDtoSchema } from '../model/schemas';
import type { PetEmotion } from '../model/types';

export const PET_MESSAGE_TYPES = {
  ping: 'ping',
  pong: 'pong',
  get: 'pet.get',
  pet: 'pet.pet',
  state: 'pet.state',
  updated: 'pet.updated',
  hatched: 'pet.hatched',
  xpGained: 'xp.gained',
  levelUp: 'level.up',
  rewardGranted: 'reward.granted',
  streakUpdated: 'streak.updated',
  error: 'error',
} as const;

const EMOTION_BY_TYPE: Record<string, PetEmotion> = {
  [PET_MESSAGE_TYPES.hatched]: 'hatching',
  [PET_MESSAGE_TYPES.xpGained]: 'eating',
  [PET_MESSAGE_TYPES.levelUp]: 'levelup',
  [PET_MESSAGE_TYPES.rewardGranted]: 'celebrate',
  [PET_MESSAGE_TYPES.streakUpdated]: 'celebrate',
};

function applyStatePayload(payload: unknown): boolean {
  const parsed = petStateDtoSchema.safeParse(payload);

  if (!parsed.success) {
    return false;
  }

  petStateReceived(toPetState(parsed.data));

  return true;
}

export function handlePetMessage(message: WsMessage): void {
  if (message.type === PET_MESSAGE_TYPES.state || message.type === PET_MESSAGE_TYPES.updated) {
    applyStatePayload(message.payload);

    return;
  }

  const emotion = EMOTION_BY_TYPE[message.type];

  if (emotion !== undefined) {
    applyStatePayload(message.payload);
    emotionTriggered(emotion);
  }
}

function toConnectionStatus(status: WsStatus): ConnectionStatus {
  return status;
}

export interface PetSocket {
  client: WsClient;
  requestState: () => boolean;
  stroke: () => boolean;
}

export function createPetSocket(onReconnect: () => void): PetSocket {
  const client = createWsClient({
    onMessage: handlePetMessage,
    onStatus: (status) => {
      connectionStatusChanged(toConnectionStatus(status));
    },
    onReconnect,
  });

  return {
    client,
    requestState: () => client.send({ type: PET_MESSAGE_TYPES.get }),
    stroke: () => client.send({ type: PET_MESSAGE_TYPES.pet }),
  };
}
