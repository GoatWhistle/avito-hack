import { renderHook } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import {
  makeCheckedInPet,
  makeCheckInInfo,
  makePet,
} from '#/features/pet/components/test-utils'
import { useAutoCheckIn } from './useAutoCheckIn'
import type { PetCelebration } from './usePetCelebration'

const makeCelebration = (): PetCelebration => ({
  emotion: null,
  xpToasts: [],
  banner: null,
  handleEvent: vi.fn(),
  clearEmotion: vi.fn(),
  dismissBanner: vi.fn(),
  celebrate: vi.fn(),
  showCheckIn: vi.fn(),
})

let celebration: PetCelebration

beforeEach(() => {
  celebration = makeCelebration()
})

describe('useAutoCheckIn', () => {
  it('shows the check-in banner when the backend applied one', () => {
    const pet = makeCheckedInPet({}, { xp_granted: 12 })

    renderHook(() => useAutoCheckIn(pet, celebration))

    expect(celebration.showCheckIn).toHaveBeenCalledTimes(1)
    expect(celebration.showCheckIn).toHaveBeenCalledWith({ xp: 12, days: 5 })
  })

  it('shows the banner only once across repeated renders', () => {
    const pet = makeCheckedInPet()
    const { rerender } = renderHook(() => useAutoCheckIn(pet, celebration))

    rerender()
    rerender()

    expect(celebration.showCheckIn).toHaveBeenCalledTimes(1)
  })

  it('shows the banner again for a later check-in', () => {
    const first = makeCheckedInPet()
    const { rerender } = renderHook(
      ({ pet }: { pet: typeof first }) => useAutoCheckIn(pet, celebration),
      { initialProps: { pet: first } },
    )

    rerender({
      pet: makeCheckedInPet({}, { streak: { ...makeCheckInInfo().streak, days: 6 } }),
    })

    expect(celebration.showCheckIn).toHaveBeenCalledTimes(2)
  })

  it('stays quiet when no check-in was applied', () => {
    renderHook(() => useAutoCheckIn(makePet(), celebration))

    expect(celebration.showCheckIn).not.toHaveBeenCalled()
    expect(celebration.celebrate).not.toHaveBeenCalled()
  })

  it('stays quiet when the flag is set but the payload is missing', () => {
    const pet = makePet({ checkin_applied: true, checkin: null })

    renderHook(() => useAutoCheckIn(pet, celebration))

    expect(celebration.showCheckIn).not.toHaveBeenCalled()
  })

  it('stays quiet without a pet', () => {
    renderHook(() => useAutoCheckIn(undefined, celebration))

    expect(celebration.showCheckIn).not.toHaveBeenCalled()
  })

  it('celebrates a level up on top of the check-in', () => {
    const pet = makeCheckedInPet({}, { level: 4, previous_level: 3 })

    renderHook(() => useAutoCheckIn(pet, celebration))

    expect(celebration.showCheckIn).toHaveBeenCalledTimes(1)
    expect(celebration.celebrate).toHaveBeenCalledWith('levelup')
  })

  it('skips the level up celebration when the level did not change', () => {
    const pet = makeCheckedInPet({}, { level: 3, previous_level: 3 })

    renderHook(() => useAutoCheckIn(pet, celebration))

    expect(celebration.celebrate).not.toHaveBeenCalled()
  })
})
