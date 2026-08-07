import { act, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { PetAvatar } from './PetAvatar'
import type { PetStage } from './types'

const HEALTHY = { satiety: 80, happiness: 80, energy: 80 }

function mockReducedMotion(reduced: boolean): void {
  vi.stubGlobal(
    'matchMedia',
    vi.fn().mockImplementation((query: string) => ({
      matches: reduced && query.includes('prefers-reduced-motion'),
      media: query,
      onchange: null,
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
      addListener: vi.fn(),
      removeListener: vi.fn(),
      dispatchEvent: vi.fn(),
    })),
  )
}

beforeEach(() => {
  mockReducedMotion(false)
})

afterEach(() => {
  vi.unstubAllGlobals()
  vi.useRealTimers()
})

describe('PetAvatar stages', () => {
  const stages: PetStage[] = ['egg', 'baby', 'teen', 'adult', 'legend']

  it.each(stages)('renders stage %s', (stage) => {
    render(<PetAvatar stage={stage} {...HEALTHY} />)

    const avatar = screen.getByRole('img')
    expect(avatar).toBeInTheDocument()
    expect(avatar).toHaveAttribute('data-stage', stage)
  })

  it('describes the pet in the aria-label', () => {
    render(<PetAvatar stage="adult" {...HEALTHY} />)

    expect(screen.getByRole('img')).toHaveAttribute(
      'aria-label',
      expect.stringContaining('Ноти'),
    )
  })

  it('accepts label overrides', () => {
    render(<PetAvatar stage="adult" {...HEALTHY} labels={{ name: 'Проша' }} />)

    expect(screen.getByRole('img')).toHaveAttribute(
      'aria-label',
      expect.stringContaining('Проша'),
    )
  })
})

describe('PetAvatar background moods', () => {
  it.each([
    [{ satiety: 10, happiness: 80, energy: 80 }, 'hungry'],
    [{ satiety: 80, happiness: 10, energy: 80 }, 'sad'],
    [{ satiety: 80, happiness: 80, energy: 5 }, 'sleeping'],
    [{ satiety: 80, happiness: 80, energy: 80 }, 'happy'],
  ])('derives mood %#', (stats, expected) => {
    render(<PetAvatar stage="teen" {...stats} idleSleepEnabled={false} />)

    expect(screen.getByRole('img')).toHaveAttribute('data-mood', expected)
  })
})

describe('PetAvatar emotions', () => {
  it('applies the emotion modifier class', () => {
    render(<PetAvatar stage="baby" {...HEALTHY} emotion="celebrate" />)

    const avatar = screen.getByRole('img')
    expect(avatar).toHaveAttribute('data-emotion', 'celebrate')
    expect(avatar.getAttribute('class')).toContain(
      'pet-avatar--emotion-celebrate',
    )
  })

  it('returns to the background state after the reset delay', () => {
    vi.useFakeTimers()
    const onEmotionEnd = vi.fn()

    render(
      <PetAvatar
        stage="baby"
        {...HEALTHY}
        emotion="eating"
        emotionResetMs={2_500}
        onEmotionEnd={onEmotionEnd}
      />,
    )

    expect(screen.getByRole('img')).toHaveAttribute('data-emotion', 'eating')

    act(() => {
      vi.advanceTimersByTime(2_600)
    })

    expect(screen.getByRole('img')).toHaveAttribute('data-emotion', 'none')
    expect(onEmotionEnd).toHaveBeenCalledTimes(1)
  })

  it('renders the hatching sequence for the egg stage', () => {
    render(<PetAvatar stage="egg" {...HEALTHY} emotion="hatching" />)

    const avatar = screen.getByRole('img')
    expect(avatar.getAttribute('class')).toContain(
      'pet-avatar--emotion-hatching',
    )
    expect(avatar.querySelector('.pet-avatar__hatchling')).not.toBeNull()
  })
})

describe('PetAvatar stroke interaction', () => {
  it('calls onStroke on click', () => {
    const onStroke = vi.fn()
    render(<PetAvatar stage="teen" {...HEALTHY} onStroke={onStroke} />)

    fireEvent.click(screen.getByRole('img'))

    expect(onStroke).toHaveBeenCalledTimes(1)
  })

  it('calls onStroke on Enter', () => {
    const onStroke = vi.fn()
    render(<PetAvatar stage="teen" {...HEALTHY} onStroke={onStroke} />)

    fireEvent.keyDown(screen.getByRole('img'), { key: 'Enter' })

    expect(onStroke).toHaveBeenCalledTimes(1)
  })

  it('calls onStroke on Space', () => {
    const onStroke = vi.fn()
    render(<PetAvatar stage="teen" {...HEALTHY} onStroke={onStroke} />)

    fireEvent.keyDown(screen.getByRole('img'), { key: ' ' })

    expect(onStroke).toHaveBeenCalledTimes(1)
  })

  it('ignores unrelated keys', () => {
    const onStroke = vi.fn()
    render(<PetAvatar stage="teen" {...HEALTHY} onStroke={onStroke} />)

    fireEvent.keyDown(screen.getByRole('img'), { key: 'a' })

    expect(onStroke).not.toHaveBeenCalled()
  })

  it('shows hearts after being stroked', () => {
    render(<PetAvatar stage="teen" {...HEALTHY} />)

    const avatar = screen.getByRole('img')
    fireEvent.click(avatar)

    expect(avatar.querySelector('.pet-avatar__heart')).not.toBeNull()
  })

  it('is focusable', () => {
    render(<PetAvatar stage="teen" {...HEALTHY} />)

    expect(screen.getByRole('img')).toHaveAttribute('tabindex', '0')
  })
})

describe('PetAvatar reduced motion', () => {
  it('marks the avatar as reduced', () => {
    mockReducedMotion(true)
    render(<PetAvatar stage="adult" {...HEALTHY} />)

    expect(screen.getByRole('img').getAttribute('class')).toContain(
      'pet-avatar--reduced',
    )
  })

  it('keeps pupils centred when motion is reduced', () => {
    mockReducedMotion(true)
    render(<PetAvatar stage="adult" {...HEALTHY} />)

    const pupil = screen
      .getByRole('img')
      .querySelector<SVGCircleElement>('.pet-avatar__pupil')

    expect(pupil?.style.transform).toBe('translate(0px, 0px)')
  })
})
