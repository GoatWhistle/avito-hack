import { describe, expect, it } from 'vitest';

import { buildWsUrl, reconnectDelayMs } from './ws';

describe('reconnectDelayMs', () => {
  it('grows exponentially from one second', () => {
    expect(reconnectDelayMs(0)).toBe(1_000);
    expect(reconnectDelayMs(1)).toBe(2_000);
    expect(reconnectDelayMs(2)).toBe(4_000);
    expect(reconnectDelayMs(3)).toBe(8_000);
  });

  it('caps the delay at thirty seconds', () => {
    expect(reconnectDelayMs(10)).toBe(30_000);
    expect(reconnectDelayMs(100)).toBe(30_000);
  });

  it('treats negative attempts as the first one', () => {
    expect(reconnectDelayMs(-5)).toBe(1_000);
  });
});

describe('buildWsUrl', () => {
  it('switches http to the ws scheme and carries the token', () => {
    const url = new URL(buildWsUrl('jwt-token'));

    expect(url.protocol).toBe('ws:');
    expect(url.pathname.endsWith('/ws')).toBe(true);
    expect(url.searchParams.get('token')).toBe('jwt-token');
  });
});
