import { describe, expect, it } from 'vitest'
import { routeTransitionKey } from './route-transition-key'

describe('routeTransitionKey', () => {
  it('returns the top-level segment for a nested path', () => {
    expect(routeTransitionKey('/pet/games')).toBe('pet')
    expect(routeTransitionKey('/items/123/edit')).toBe('items')
  })

  it('keeps the same key across tabs within a section', () => {
    expect(routeTransitionKey('/pet')).toBe(routeTransitionKey('/pet/progress'))
    expect(routeTransitionKey('/pet/games')).toBe(
      routeTransitionKey('/pet/rewards'),
    )
  })

  it('changes the key when moving between sections', () => {
    expect(routeTransitionKey('/items')).not.toBe(routeTransitionKey('/pet'))
    expect(routeTransitionKey('/pet')).not.toBe(routeTransitionKey('/sign-in'))
  })

  it('returns a stable key for the root path', () => {
    expect(routeTransitionKey('/')).toBe('/')
    expect(routeTransitionKey('')).toBe('/')
  })

  it('ignores a trailing slash', () => {
    expect(routeTransitionKey('/pet/')).toBe('pet')
  })

  it('ignores a query string if present', () => {
    expect(routeTransitionKey('/items?page=2')).toBe('items')
  })
})
