import { useQueryClient } from '@tanstack/react-query';
import { useUnit } from 'effector-react';
import { useCallback, useMemo } from 'react';

import {
  $pet,
  $reactiveEmotion,
  petKeys,
  resolveBackgroundEmotion,
  useEmotionDecay,
  useIdleSleep,
  usePetSocket,
  xpProgressPercent,
  xpRemaining,
  isMaxLevel,
  type ConnectionStatus,
  type PetEmotion,
  type PetState,
} from '@/entities/pet';

export interface PetLive {
  pet: PetState | null;
  emotion: PetEmotion;
  status: ConnectionStatus;
  xpPercent: number;
  xpLeft: number;
  maxLevel: boolean;
  stroke: () => void;
}

export function usePetLive(enabled: boolean): PetLive {
  const queryClient = useQueryClient();
  const [pet, reactiveEmotion] = useUnit([$pet, $reactiveEmotion]);
  const isAsleep = useIdleSleep();

  const onReconnect = useCallback(() => {
    void queryClient.invalidateQueries({ queryKey: petKeys.all });
  }, [queryClient]);

  const socket = usePetSocket(enabled, onReconnect);
  useEmotionDecay(reactiveEmotion);

  const background = resolveBackgroundEmotion(pet, isAsleep);
  const emotion = reactiveEmotion ?? background;

  const stroke = useCallback(() => {
    socket.stroke();
  }, [socket]);

  return useMemo(
    () => ({
      pet,
      emotion,
      status: socket.status,
      xpPercent: xpProgressPercent(pet),
      xpLeft: xpRemaining(pet),
      maxLevel: isMaxLevel(pet),
      stroke,
    }),
    [pet, emotion, socket.status, stroke],
  );
}
