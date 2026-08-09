import { describe, expect, it } from 'vitest'
import { makePet } from '../test-utils'
import { dashboardStats } from './dashboard-stats'

describe('dashboardStats', () => {
  it('lists energy, happiness and satiety in order', () => {
    const stats = dashboardStats(makePet())

    expect(stats.map((stat) => stat.key)).toEqual([
      'energy',
      'happiness',
      'satiety',
    ])
  })

  it('leaves streak and xp to the pet card', () => {
    const keys = dashboardStats(makePet()).map((stat) => stat.key)

    expect(keys).not.toContain('streak')
    expect(keys).not.toContain('xp')
  })

  it('uses raw values as percents for the 0..100 stats', () => {
    const stats = dashboardStats(
      makePet({ energy: 42, happiness: 13, satiety: 88 }),
    )

    expect(stats[0]).toMatchObject({ value: 42, percent: 42 })
    expect(stats[1]).toMatchObject({ value: 13, percent: 13 })
    expect(stats[2]).toMatchObject({ value: 88, percent: 88 })
  })

  it('clamps out-of-range values into the bar', () => {
    const stats = dashboardStats(makePet({ energy: 140, happiness: -20 }))

    expect(stats[0].percent).toBe(100)
    expect(stats[1].percent).toBe(0)
  })
})
