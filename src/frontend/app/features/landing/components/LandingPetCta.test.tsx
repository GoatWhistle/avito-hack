import { screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { LandingPetCta } from './LandingPetCta'
import { renderWithProviders } from '#/features/items/components/test-utils'

describe('LandingPetCta', () => {
  it('sends a guest to sign-up', () => {
    renderWithProviders(<LandingPetCta />, { userId: null })

    expect(
      screen.getByRole('link', { name: 'Завести питомца' }),
    ).toHaveAttribute('href', '/sign-up')
    expect(
      screen.queryByRole('link', { name: 'К моему питомцу' }),
    ).not.toBeInTheDocument()
  })

  it('sends an authenticated user straight to the pet', () => {
    renderWithProviders(<LandingPetCta />)

    expect(
      screen.getByRole('link', { name: 'К моему питомцу' }),
    ).toHaveAttribute('href', '/pet')
    expect(
      screen.queryByRole('link', { name: 'Завести питомца' }),
    ).not.toBeInTheDocument()
  })
})
