import { describe, expect, it } from 'vitest';

import { formatPrice, fromKopeks, toKopeks } from './format';

describe('format', () => {
  it('converts units to kopeks and back', () => {
    expect(toKopeks(199.99)).toBe(19999);
    expect(fromKopeks(19999)).toBeCloseTo(199.99);
  });

  it('rounds fractional kopeks', () => {
    expect(toKopeks(0.005)).toBe(1);
    expect(toKopeks(0)).toBe(0);
  });

  it('formats price for the given locale', () => {
    expect(formatPrice(150000, 'ru')).toContain('1');
    expect(formatPrice(150000, 'en')).toContain('$');
  });
});
