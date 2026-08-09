import { screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { PetSpeech } from './PetSpeech'
import { renderWithProviders } from './test-utils'

describe('PetSpeech', () => {
  it('speaks in the first person about hunger', () => {
    renderWithProviders(<PetSpeech speech="hungry" petName="Ноти" />)

    expect(screen.getByTestId('pet-speech')).toHaveTextContent(/голодный/i)
  })

  it('announces the speaker for assistive technology', () => {
    renderWithProviders(<PetSpeech speech="happy" petName="Ноти" />)

    expect(screen.getByText('Ноти говорит:')).toBeInTheDocument()
  })

  it('is announced politely on change', () => {
    renderWithProviders(<PetSpeech speech="checkIn" petName="Ноти" />)

    expect(screen.getByTestId('pet-speech')).toHaveAttribute(
      'aria-live',
      'polite',
    )
  })
})
