import { API_BASE_URL } from '@/shared/config/env';

import { tokenStorage } from './token-storage';

export interface WsMessage {
  type: string;
  request_id?: string;
  payload?: unknown;
}

export type WsStatus = 'idle' | 'connecting' | 'open' | 'closed';

export interface WsClientHandlers {
  onMessage: (message: WsMessage) => void;
  onStatus: (status: WsStatus) => void;
  onReconnect: () => void;
}

export interface WsClient {
  connect: () => void;
  disconnect: () => void;
  send: (message: WsMessage) => boolean;
}

const HEARTBEAT_INTERVAL_MS = 25_000;
const RECONNECT_BASE_MS = 1_000;
const RECONNECT_MAX_MS = 30_000;
const NORMAL_CLOSURE = 1000;

export function buildWsUrl(token: string): string {
  const httpUrl = new URL(`${API_BASE_URL}/ws`);
  httpUrl.protocol = httpUrl.protocol === 'https:' ? 'wss:' : 'ws:';
  httpUrl.searchParams.set('token', token);

  return httpUrl.toString();
}

export function reconnectDelayMs(attempt: number): number {
  const delay = RECONNECT_BASE_MS * Math.pow(2, Math.max(0, attempt));

  return Math.min(delay, RECONNECT_MAX_MS);
}

function parseMessage(raw: unknown): WsMessage | null {
  if (typeof raw !== 'string') {
    return null;
  }

  try {
    const parsed: unknown = JSON.parse(raw);

    if (typeof parsed !== 'object' || parsed === null || !('type' in parsed)) {
      return null;
    }

    const { type } = parsed;

    return typeof type === 'string' ? (parsed as WsMessage) : null;
  } catch {
    return null;
  }
}

export function createWsClient(handlers: WsClientHandlers): WsClient {
  let socket: WebSocket | null = null;
  let heartbeat: ReturnType<typeof setInterval> | null = null;
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  let attempt = 0;
  let stopped = true;
  let hasConnectedOnce = false;

  const clearTimers = (): void => {
    if (heartbeat !== null) {
      clearInterval(heartbeat);
      heartbeat = null;
    }
    if (reconnectTimer !== null) {
      clearTimeout(reconnectTimer);
      reconnectTimer = null;
    }
  };

  const scheduleReconnect = (): void => {
    if (stopped || reconnectTimer !== null) {
      return;
    }

    const delay = reconnectDelayMs(attempt);
    attempt += 1;
    reconnectTimer = setTimeout(() => {
      reconnectTimer = null;
      open();
    }, delay);
  };

  const startHeartbeat = (): void => {
    heartbeat = setInterval(() => {
      if (socket?.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({ type: 'ping' }));
      }
    }, HEARTBEAT_INTERVAL_MS);
  };

  function open(): void {
    const token = tokenStorage.get();

    if (stopped || token === null) {
      return;
    }

    handlers.onStatus('connecting');

    try {
      socket = new WebSocket(buildWsUrl(token));
    } catch {
      handlers.onStatus('closed');
      scheduleReconnect();

      return;
    }

    socket.onopen = () => {
      attempt = 0;
      handlers.onStatus('open');
      startHeartbeat();

      if (hasConnectedOnce) {
        handlers.onReconnect();
      }
      hasConnectedOnce = true;
    };

    socket.onmessage = (event: MessageEvent<unknown>) => {
      const message = parseMessage(event.data);

      if (message !== null) {
        handlers.onMessage(message);
      }
    };

    socket.onerror = () => {
      handlers.onStatus('closed');
    };

    socket.onclose = () => {
      clearTimers();
      socket = null;
      handlers.onStatus('closed');
      scheduleReconnect();
    };
  }

  return {
    connect: () => {
      if (!stopped) {
        return;
      }
      stopped = false;
      attempt = 0;
      open();
    },

    disconnect: () => {
      stopped = true;
      clearTimers();
      hasConnectedOnce = false;

      if (socket !== null) {
        socket.onclose = null;
        socket.close(NORMAL_CLOSURE, 'client disconnect');
        socket = null;
      }

      handlers.onStatus('idle');
    },

    send: (message: WsMessage): boolean => {
      if (socket?.readyState !== WebSocket.OPEN) {
        return false;
      }
      socket.send(JSON.stringify(message));

      return true;
    },
  };
}
