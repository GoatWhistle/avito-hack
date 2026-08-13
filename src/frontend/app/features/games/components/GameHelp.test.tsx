import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it } from 'vitest'
import { GameHelp } from './GameHelp'
import { renderWithProviders } from './test-utils'

const props = {
  storageKey: 'sample',
  title: 'Как играть',
  intro: 'Короткое описание',
  closeLabel: 'Понятно',
  hideLabel: 'Скрыть подсказку',
  openLabel: 'Как играть',
}

beforeEach(() => {
  window.localStorage.removeItem('avito-hack.sample.help-dismissed')
  window.localStorage.removeItem('avito-hack.other.help-dismissed')
})

describe('GameHelp', () => {
  it('shows the panel with its rows until it is dismissed', async () => {
    renderWithProviders(
      <GameHelp
        {...props}
        rows={[
          {
            tone: 'bg-success text-success-foreground',
            sample: 'д',
            label: 'Первое правило',
          },
        ]}
        outro="Финальная подсказка"
      />,
    )

    expect(screen.getByText('Короткое описание')).toBeInTheDocument()
    expect(screen.getByText('Первое правило')).toBeInTheDocument()
    expect(screen.getByText('Финальная подсказка')).toBeInTheDocument()

    await userEvent.click(screen.getByRole('button', { name: 'Понятно' }))

    expect(screen.queryByText('Короткое описание')).not.toBeInTheDocument()
    expect(
      screen.getByRole('button', { name: 'Как играть' }),
    ).toBeInTheDocument()
  })

  it('remembers the dismissal and reopens on demand', async () => {
    const { unmount } = renderWithProviders(<GameHelp {...props} />)

    await userEvent.click(screen.getByRole('button', { name: 'Понятно' }))
    unmount()

    renderWithProviders(<GameHelp {...props} />)
    expect(screen.queryByText('Короткое описание')).not.toBeInTheDocument()

    await userEvent.click(screen.getByRole('button', { name: 'Как играть' }))
    expect(screen.getByText('Короткое описание')).toBeInTheDocument()
  })

  it('gives the corner control its own accessible name', async () => {
    renderWithProviders(<GameHelp {...props} />)

    await userEvent.click(
      screen.getByRole('button', { name: 'Скрыть подсказку' }),
    )

    expect(screen.queryByText('Короткое описание')).not.toBeInTheDocument()
  })

  it('keeps the dismissed state separate per game', async () => {
    const { unmount } = renderWithProviders(<GameHelp {...props} />)
    await userEvent.click(screen.getByRole('button', { name: 'Понятно' }))
    unmount()

    renderWithProviders(
      <GameHelp {...props} storageKey="other" intro="Другая игра" />,
    )

    expect(screen.getByText('Другая игра')).toBeInTheDocument()
  })
})
