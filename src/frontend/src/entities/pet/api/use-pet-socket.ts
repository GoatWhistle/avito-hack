import { useUnit } from 'effector-react';
import { useCallback, useEffect, useMemo, useRef } from 'react';

import { emotionSettled, petStateCleared, $connectionStatus } from '../model/pet-store';
import type { ConnectionStatus } from '../model/pet-store';

import { createPetSocket, type PetSocket } from './ws-bridge';

const EMOTION_DURATION_MS = 2_600;

export interface PetSocketApi {
  status: ConnectionStatus;
  stroke: () => boolean;
  requestState: () => boolean;
}

export function usePetSocket(enabled: boolean, onReconnect?: () => void): PetSocketApi {
  const status = useUnit($connectionStatus);
  const socketRef = useRef<PetSocket | null>(null);
  const reconnectRef = useRef(onReconnect);

  useEffect(() => {
    reconnectRef.current = onReconnect;
  }, [onReconnect]);

  useEffect(() => {
    if (!enabled) {
      return;
    }

    const socket = createPetSocket(() => {
      reconnectRef.current?.();
      socket.requestState();
    });

    socketRef.current = socket;
    socket.client.connect();

    const initial = setTimeout(() => {
      socket.requestState();
    }, 0);

    return () => {
      clearTimeout(initial);
      socket.client.disconnect();
      socketRef.current = null;
      petStateCleared();
    };
  }, [enabled]);

  const stroke = useCallback(() => socketRef.current?.stroke() ?? false, []);
  const requestState = useCallback(() => socketRef.current?.requestState() ?? false, []);

  return useMemo(() => ({ status, stroke, requestState }), [status, stroke, requestState]);
}

export function useEmotionDecay(emotion: string | null): void {
  useEffect(() => {
    if (emotion === null) {
      return;
    }

    const timer = setTimeout(() => {
      emotionSettled();
    }, EMOTION_DURATION_MS);

    return () => {
      clearTimeout(timer);
    };
  }, [emotion]);
}
