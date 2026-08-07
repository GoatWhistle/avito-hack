import { describe, expect, it } from 'vitest'
import {
  formatBytes,
  formatDate,
  formatPrice,
  kopeksToRubles,
  rublesToKopeks,
} from './format'
import { availableActions, isFavoritable, statusTone } from './status'
import type { ItemStatus } from '#/features/items/types'

describe('kopeksToRubles', () => {
  it('converts whole and fractional amounts', () => {
    expect(kopeksToRubles(150_000)).toBe(1500)
    expect(kopeksToRubles(1_050)).toBe(10.5)
  })

  it('rounds sub kopek input', () => {
    expect(kopeksToRubles(1_00.4)).toBe(1)
  })
})

describe('rublesToKopeks', () => {
  it('converts back without floating point drift', () => {
    expect(rublesToKopeks(10.5)).toBe(1_050)
    expect(rublesToKopeks(0.1)).toBe(10)
  })
})

describe('formatPrice', () => {
  it('omits decimals for whole rubles', () => {
    expect(formatPrice(150_000, 'ru-RU').replace(/\s/g, '')).toBe('1500')
  })

  it('keeps two decimals for fractional rubles', () => {
    expect(formatPrice(1_050, 'ru-RU')).toBe('10,50')
  })

  it('respects the requested locale', () => {
    expect(formatPrice(150_000, 'en-US')).toBe('1,500')
  })
})

describe('formatBytes', () => {
  it('prints whole megabytes without decimals', () => {
    expect(formatBytes(2 * 1024 * 1024)).toBe('2 MB')
  })

  it('rounds partial megabytes to one decimal', () => {
    expect(formatBytes(1.5 * 1024 * 1024)).toBe('1.5 MB')
  })

  it('handles zero', () => {
    expect(formatBytes(0)).toBe('0 MB')
  })
})

describe('formatDate', () => {
  it('formats a valid iso date', () => {
    expect(formatDate('2026-01-15T09:00:00Z', 'ru-RU')).toContain('2026')
  })

  it('returns an empty string for an invalid date', () => {
    expect(formatDate('not-a-date', 'ru-RU')).toBe('')
  })
})

describe('availableActions', () => {
  it('lists the transitions of every status', () => {
    expect(availableActions('draft')).toEqual(['publish', 'submit', 'archive'])
    expect(availableActions('moderation')).toEqual(['publish', 'archive'])
    expect(availableActions('published')).toEqual(['sell', 'archive'])
    expect(availableActions('sold')).toEqual(['archive'])
    expect(availableActions('archived')).toEqual(['restore'])
  })

  it('returns nothing for an unknown status', () => {
    expect(availableActions('ghost' as ItemStatus)).toEqual([])
  })
})

describe('isFavoritable', () => {
  it('allows published and moderated items', () => {
    expect(isFavoritable('published')).toBe(true)
    expect(isFavoritable('moderation')).toBe(true)
  })

  it('rejects the remaining statuses', () => {
    expect(isFavoritable('draft')).toBe(false)
    expect(isFavoritable('sold')).toBe(false)
    expect(isFavoritable('archived')).toBe(false)
  })
})

describe('statusTone', () => {
  it('maps every status to a tone', () => {
    expect(statusTone('draft')).toBe('neutral')
    expect(statusTone('moderation')).toBe('info')
    expect(statusTone('published')).toBe('success')
    expect(statusTone('sold')).toBe('info')
    expect(statusTone('archived')).toBe('muted')
  })

  it('falls back to neutral for an unknown status', () => {
    expect(statusTone('ghost' as ItemStatus)).toBe('neutral')
  })
})
