import { screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { renderWithProviders } from '#/features/items/components/test-utils'
import type { Pet } from '#/features/pet/types'
import { StreakMechanicsCard } from './StreakMechanicsCard'

const makePet = (overrides: Partial<Pet> = {}): Pet =>
  ({
    id: 'pet-1',
    user_id: 'user-1',
    name: 'Енот',
    stage: 'baby',
    state: 'happy',
    level: 1,
    xp: 0,
    next_level_xp: 5,
    satiety: 70,
    happiness: 70,
    energy: 100,
    streak_days: 0,
    freezes: 0,
    is_hatched: true,
    ...overrides,
  }) as Pet

describe('StreakMechanicsCard', () => {
  it('renders the three streak rules', async () => {
    renderWithProviders(<StreakMechanicsCard pet={makePet()} />)

    expect(await screen.findByTestId('streak-rule-bonus')).toBeInTheDocument()
    expect(screen.getByTestId('streak-rule-reset')).toBeInTheDocument()
    expect(screen.getByTestId('streak-rule-freeze')).toBeInTheDocument()
  })

  it('states the real backend multiplier and threshold', async () => {
    renderWithProviders(<StreakMechanicsCard pet={makePet()} />)

    const bonus = await screen.findByTestId('streak-rule-bonus')

    expect(bonus).toHaveTextContent('x1.5')
    expect(bonus).toHaveTextContent('7')
  })

  it('shows days remaining before the multiplier turns on', async () => {
    renderWithProviders(
      <StreakMechanicsCard pet={makePet({ streak_days: 2 })} />,
    )

    expect(
      await screen.findByTestId('streak-rule-badge-bonus'),
    ).toHaveTextContent('5')
  })

  it('marks the multiplier active for a long streak', async () => {
    renderWithProviders(
      <StreakMechanicsCard pet={makePet({ streak_days: 41 })} />,
    )

    expect(
      await screen.findByTestId('streak-rule-badge-bonus'),
    ).toHaveTextContent('Активен')
  })

  it('reports freezes as unavailable when the pet holds none', async () => {
    renderWithProviders(<StreakMechanicsCard pet={makePet({ freezes: 0 })} />)

    expect(
      await screen.findByTestId('streak-rule-badge-freeze'),
    ).toHaveTextContent('Недоступно')
  })

  it('counts the freezes the pet actually holds', async () => {
    renderWithProviders(<StreakMechanicsCard pet={makePet({ freezes: 2 })} />)

    expect(
      await screen.findByTestId('streak-rule-badge-freeze'),
    ).toHaveTextContent('2')
  })

  it('explains that freezes come from milestones, capped at three', async () => {
    renderWithProviders(<StreakMechanicsCard pet={makePet()} />)

    const freeze = await screen.findByTestId('streak-rule-freeze')

    expect(freeze).toHaveTextContent('3, 7, 14, 30')
    expect(freeze).toHaveTextContent('3')
  })
})
