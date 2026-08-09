import { screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { renderWithProviders } from '#/features/items/components/test-utils'
import type { BadgeItem } from '#/features/rewards/types'
import { AchievementCard } from './AchievementCard'

const makeBadge = (overrides: Partial<BadgeItem> = {}): BadgeItem => ({
  id: 'explorer',
  name: 'Исследователь',
  description: 'Добавить 5 объявлений в избранное',
  icon_url: '',
  earned_at: '',
  progress_current: 0,
  progress_target: 5,
  ...overrides,
})

describe('AchievementCard', () => {
  it('shows what is left to do for a locked badge', async () => {
    renderWithProviders(
      <ul>
        <AchievementCard badge={makeBadge({ progress_current: 2 })} />
      </ul>,
    )

    expect(await screen.findByTestId('badge-progress')).toHaveTextContent(
      '2 из 5',
    )
    expect(screen.getByText(/Осталось 3/)).toBeInTheDocument()
  })

  it('states the unlock condition on the card', async () => {
    renderWithProviders(
      <ul>
        <AchievementCard badge={makeBadge()} />
      </ul>,
    )

    expect(
      await screen.findByText('Добавить 5 объявлений в избранное'),
    ).toBeInTheDocument()
  })

  it('replaces progress with the earn date once earned', async () => {
    renderWithProviders(
      <ul>
        <AchievementCard
          badge={makeBadge({
            progress_current: 5,
            earned_at: '2026-08-01T10:00:00Z',
          })}
        />
      </ul>,
    )

    const tile = await screen.findByTestId('badge-tile')
    expect(tile).toHaveAttribute('data-earned', 'true')
    expect(screen.queryByTestId('badge-progress')).not.toBeInTheDocument()
  })

  it('falls back to the plain locked label without progress data', async () => {
    renderWithProviders(
      <ul>
        <AchievementCard badge={makeBadge({ progress_target: 0 })} />
      </ul>,
    )

    expect(await screen.findByText('Ещё не получен')).toBeInTheDocument()
    expect(screen.queryByTestId('badge-progress')).not.toBeInTheDocument()
  })
})
