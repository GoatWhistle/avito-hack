import { screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { GamePlayScreen } from './GamePlayScreen'
import { makeGameState, renderWithProviders } from './test-utils'

const state = vi.fn()
const startRound = vi.fn()
const guess = vi.fn()
const claimReward = vi.fn()

vi.mock('#/features/games/repository', () => ({
  gameRepository: {
    state: (slug: string) => state(slug),
    startRound: (slug: string) => startRound(slug),
    guess: (slug: string, roundId: string, move: unknown) =>
      guess(slug, roundId, move),
    claimReward: () => claimReward(),
  },
}))

const renderGame = (slug: string) =>
  renderWithProviders(<GamePlayScreen />, {
    route: `/play/${slug}`,
    path: '/play/:gameSlug',
  })

beforeEach(() => {
  vi.clearAllMocks()
  window.localStorage.removeItem('avito-hack.bukovki.help-dismissed')
  window.localStorage.removeItem('avito-hack.moreless.help-dismissed')
  window.localStorage.removeItem('avito-hack.raccoonjump.help-dismissed')
  state.mockResolvedValue(makeGameState())
})

describe('GamePlayScreen', () => {
  it('shows a dismissible help panel next to the game', async () => {
    renderGame('bukovki')

    expect(await screen.findByText('Буква на своём месте')).toBeInTheDocument()

    await userEvent.click(screen.getByRole('button', { name: 'Понятно' }))

    expect(screen.queryByText('Буква на своём месте')).not.toBeInTheDocument()
    expect(
      screen.getByRole('button', { name: 'Как играть' }),
    ).toBeInTheDocument()
  })

  it('gives every game the same help panel', async () => {
    renderGame('higher-lower')

    expect(await screen.findByText('Слева цена уже известна')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Понятно' })).toBeInTheDocument()
  })

  it('reopens the help after it was dismissed', async () => {
    renderGame('bukovki')

    await userEvent.click(
      await screen.findByRole('button', { name: 'Понятно' }),
    )
    await userEvent.click(screen.getByRole('button', { name: 'Как играть' }))

    expect(await screen.findByText('Буква на своём месте')).toBeInTheDocument()
  })

  it('keeps the help dismissed per game', async () => {
    window.localStorage.setItem('avito-hack.bukovki.help-dismissed', '1')

    renderGame('bukovki')

    expect(
      await screen.findByRole('button', { name: 'Как играть' }),
    ).toBeInTheDocument()
    expect(screen.queryByText('Буква на своём месте')).not.toBeInTheDocument()
  })

  it('keeps the global streak card out of the game screen', async () => {
    state.mockResolvedValue(
      makeGameState({
        slug: 'bukovki',
        streak: { current_days: 7, best_days: 7, reward_ready: true },
      }),
    )

    renderGame('bukovki')
    await screen.findByRole('button', { name: 'Понятно' })

    expect(screen.queryByText('Недельный стрик')).not.toBeInTheDocument()
    expect(
      screen.queryByRole('button', { name: 'Забрать награду' }),
    ).not.toBeInTheDocument()
    expect(claimReward).not.toHaveBeenCalled()
  })

  it('resolves the raccoon jump game from the registry by its path', async () => {
    renderGame('noti-jump')

    const heading = await screen.findByRole('heading', {
      name: 'Прыжки Ноти',
      level: 1,
    })

    expect(heading).toBeInTheDocument()
    expect(
      within(heading.closest('header') as HTMLElement).getByText(
        'Прыгайте по платформам как можно выше',
      ),
    ).toBeInTheDocument()
  })

  it('resolves the raccoon jump game from the registry by its slug', async () => {
    renderGame('raccoonjump')

    expect(
      await screen.findByRole('heading', { name: 'Прыжки Ноти', level: 1 }),
    ).toBeInTheDocument()
  })

  it('renders the raccoon jump help panel in the shared style', async () => {
    renderGame('noti-jump')

    expect(await screen.findByText('ПК: ← → или A / D')).toBeInTheDocument()
    expect(
      screen.getByText('Телефон: касайтесь левой или правой стороны'),
    ).toBeInTheDocument()

    await userEvent.click(screen.getByRole('button', { name: 'Понятно' }))

    expect(screen.queryByText('ПК: ← → или A / D')).not.toBeInTheDocument()
    expect(
      screen.getByRole('button', { name: 'Как играть' }),
    ).toBeInTheDocument()
  })

  it('keeps the raccoon jump control hints only in the help panel', async () => {
    window.localStorage.setItem('avito-hack.raccoonjump.help-dismissed', '1')

    renderGame('noti-jump')

    await screen.findByRole('button', { name: 'Как играть' })

    expect(screen.queryByText('ПК: ← → или A / D')).not.toBeInTheDocument()
    expect(
      screen.queryByText('Телефон: касайтесь левой или правой стороны'),
    ).not.toBeInTheDocument()
  })

  it('falls back when the game is unknown', async () => {
    renderGame('nope')

    expect(await screen.findByText('Игр пока нет')).toBeInTheDocument()
  })
})
