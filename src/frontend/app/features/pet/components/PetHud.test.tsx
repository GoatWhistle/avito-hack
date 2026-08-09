import { screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { levelProgress } from '#/features/pet/lib'
import { PetHud } from './PetHud'
import { makePet, renderWithProviders } from './test-utils'

const renderHud = (
  overrides: Parameters<typeof makePet>[0] = {},
  streakAtRisk = false,
) => {
  const pet = makePet(overrides)

  return renderWithProviders(
    <PetHud
      pet={pet}
      progress={levelProgress(pet)}
      streakAtRisk={streakAtRisk}
    />,
  )
}

describe('PetHud', () => {
  it('shows the pet name, level and xp progress', () => {
    renderHud()

    expect(
      screen.getByRole('heading', { name: 'Ноти · уровень 3' }),
    ).toBeInTheDocument()
    expect(screen.getByText('15 / 22 XP')).toBeInTheDocument()
    expect(screen.getByRole('progressbar')).toHaveAttribute(
      'aria-valuenow',
      '68',
    )
  })

  it('shows compact streak, stage and state chips', () => {
    renderHud()

    expect(screen.getByText('серия 4 дн.')).toBeInTheDocument()
    expect(screen.getByText('Малыш')).toBeInTheDocument()
    expect(screen.getByText('Спокоен')).toBeInTheDocument()
  })

  it('keeps the streak chip readable when there is no streak', () => {
    renderHud({ streak_days: 0 })

    expect(screen.getByText('серия 0 дн.')).toBeInTheDocument()
  })

  it('marks the max level instead of an xp ratio', () => {
    renderHud({ level: 15, xp: 500, next_level_xp: 0 })

    expect(screen.getByText('МАКС')).toBeInTheDocument()
    expect(screen.getByRole('progressbar')).toHaveAttribute(
      'aria-label',
      'Достигнут максимальный уровень',
    )
  })
})
