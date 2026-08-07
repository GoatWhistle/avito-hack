import { screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { StreakCard } from './StreakCard'
import { renderWithProviders } from './test-utils'

describe('StreakCard', () => {
  it('encourages starting a streak when there are no days', () => {
    renderWithProviders(<StreakCard days={0} freezes={0} atRisk={false} />)

    expect(screen.getByText('Серия ещё не начата')).toBeInTheDocument()
    expect(
      screen.getByText('Отметьтесь сегодня, чтобы начать серию'),
    ).toBeInTheDocument()
  })

  it('shows an active streak and the next day hint', () => {
    renderWithProviders(<StreakCard days={4} freezes={0} atRisk={false} />)

    expect(screen.getByText('4 дня подряд')).toBeInTheDocument()
    expect(screen.getByText('Завтра будет 5-й день')).toBeInTheDocument()
  })

  it('warns when the streak is at risk', () => {
    renderWithProviders(<StreakCard days={9} freezes={0} atRisk />)

    expect(screen.getByText('9 дней подряд')).toBeInTheDocument()
    expect(screen.getByText(/сгорит/i)).toBeInTheDocument()
  })

  it('warns about the risk even when the streak is empty', () => {
    renderWithProviders(<StreakCard days={0} freezes={0} atRisk />)

    expect(screen.getByText('Серия ещё не начата')).toBeInTheDocument()
    expect(screen.getByText(/сгорит/i)).toBeInTheDocument()
  })

  it('hides the freeze badge when none are available', () => {
    renderWithProviders(<StreakCard days={3} freezes={0} atRisk={false} />)

    expect(screen.queryByText(/заморозк/i)).not.toBeInTheDocument()
  })

  it('shows the freeze badge when freezes remain', () => {
    renderWithProviders(<StreakCard days={3} freezes={2} atRisk={false} />)

    expect(screen.getByText(/заморозк/i)).toBeInTheDocument()
  })

  it('labels the section for assistive technology', () => {
    renderWithProviders(<StreakCard days={2} freezes={0} atRisk={false} />)

    expect(screen.getByRole('region', { name: 'Серия' })).toBeInTheDocument()
  })
})
