import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it } from 'vitest'
import { PetHintBanner } from './PetHintBanner'
import { renderWithProviders } from './test-utils'
import { petHintStorageKey } from '#/features/items/lib'

beforeEach(() => {
  window.localStorage.removeItem(petHintStorageKey)
})

describe('PetHintBanner', () => {
  it('explains how listing activity feeds the pet', () => {
    renderWithProviders(<PetHintBanner />)

    expect(screen.getByTestId('pet-hint-banner')).toBeInTheDocument()
    expect(screen.getByText('Зачем здесь питомец')).toBeInTheDocument()
  })

  it('closes on the dismiss control and remembers the choice', async () => {
    renderWithProviders(<PetHintBanner />)

    await userEvent.click(screen.getByRole('button', { name: 'Закрыть' }))

    expect(screen.queryByTestId('pet-hint-banner')).not.toBeInTheDocument()
    expect(window.localStorage.getItem(petHintStorageKey)).toBe('dismissed')
  })

  it('stays hidden on a later render once dismissed', async () => {
    const first = renderWithProviders(<PetHintBanner />)
    await userEvent.click(screen.getByRole('button', { name: 'Закрыть' }))
    first.unmount()

    renderWithProviders(<PetHintBanner />)

    expect(screen.queryByTestId('pet-hint-banner')).not.toBeInTheDocument()
  })

  it('does not render when the choice was stored before mount', () => {
    window.localStorage.setItem(petHintStorageKey, 'dismissed')

    renderWithProviders(<PetHintBanner />)

    expect(screen.queryByTestId('pet-hint-banner')).not.toBeInTheDocument()
  })

  it('reaches the dismiss control from the keyboard', async () => {
    renderWithProviders(<PetHintBanner />)

    await userEvent.tab()
    await userEvent.keyboard('{Enter}')

    expect(screen.queryByTestId('pet-hint-banner')).not.toBeInTheDocument()
  })
})
