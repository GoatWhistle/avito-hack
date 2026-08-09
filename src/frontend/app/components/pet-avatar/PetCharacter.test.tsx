import { act, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { PetCharacter } from './PetCharacter'
import { RACCOON_ADULT_SRC, RACCOON_TEEN_SRC } from './lottie-sources'
import type { PetStage } from './types'

type Listener = (...args: unknown[]) => void

const listeners = new Map<string, Listener[]>()

const emit = (event: string): void => {
  act(() => {
    for (const listener of listeners.get(event) ?? []) {
      listener()
    }
  })
}

vi.mock('@lottiefiles/dotlottie-react', () => ({
  DotLottieReact: ({
    src,
    dotLottieRefCallback,
  }: {
    src: string
    dotLottieRefCallback?: (instance: unknown) => void
  }) => {
    dotLottieRefCallback?.({
      addEventListener: (event: string, listener: Listener) => {
        listeners.set(event, [...(listeners.get(event) ?? []), listener])
      },
      removeEventListener: () => {},
    })

    return <canvas data-testid="dotlottie-canvas" data-src={src} />
  },
  setWasmUrl: vi.fn(),
}))

const HEALTHY = { satiety: 80, happiness: 80, energy: 80 }
const ALL_STAGES: PetStage[] = ['baby', 'teen', 'adult', 'legend']

beforeEach(() => {
  listeners.clear()
  vi.stubGlobal(
    'matchMedia',
    vi.fn().mockImplementation((query: string) => ({
      matches: false,
      media: query,
      onchange: null,
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
      addListener: vi.fn(),
      removeListener: vi.fn(),
      dispatchEvent: vi.fn(),
    })),
  )
})

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('PetCharacter renders the designer raccoon everywhere', () => {
  it.each<[PetStage, string]>([
    ['baby', RACCOON_TEEN_SRC],
    ['teen', RACCOON_TEEN_SRC],
    ['adult', RACCOON_ADULT_SRC],
    ['legend', RACCOON_ADULT_SRC],
  ])('uses a designer animation on stage %s', (stage, src) => {
    render(<PetCharacter stage={stage} {...HEALTHY} />)

    expect(screen.getByTestId('pet-lottie')).toHaveAttribute(
      'data-stage',
      stage,
    )
    expect(screen.getByTestId('dotlottie-canvas')).toHaveAttribute(
      'data-src',
      src,
    )
  })

  it('never falls back to a hand drawn svg avatar', () => {
    for (const stage of ALL_STAGES) {
      const view = render(<PetCharacter stage={stage} {...HEALTHY} />)

      expect(view.container.querySelector('svg.pet-avatar')).toBeNull()
      view.unmount()
    }
  })

  it('never renders an egg shell on any stage', () => {
    for (const stage of ALL_STAGES) {
      const view = render(<PetCharacter stage={stage} {...HEALTHY} />)

      expect(view.container.querySelector('.pet-lottie__shell')).toBeNull()
      view.unmount()
    }
  })

  it('crowns only the legend stage', () => {
    const { container, rerender } = render(
      <PetCharacter stage="legend" {...HEALTHY} />,
    )
    expect(container.querySelector('.pet-lottie__crown')).not.toBeNull()
    expect(container.querySelector('.pet-lottie__halo')).not.toBeNull()

    rerender(<PetCharacter stage="adult" {...HEALTHY} />)
    expect(container.querySelector('.pet-lottie__crown')).toBeNull()
  })

  it('exposes the mood on the wrapper', () => {
    render(
      <PetCharacter stage="teen" satiety={10} happiness={80} energy={80} />,
    )

    expect(screen.getByTestId('pet-lottie')).toHaveAttribute(
      'data-mood',
      'hungry',
    )
  })
})

describe('PetCharacter loading and error states', () => {
  it('shows a skeleton until the animation reports load', () => {
    const { container } = render(<PetCharacter stage="teen" {...HEALTHY} />)

    expect(container.querySelector('.pet-lottie__skeleton')).not.toBeNull()

    emit('load')

    expect(container.querySelector('.pet-lottie__skeleton')).toBeNull()
  })

  it('shows the placeholder instead of the animation on load error', () => {
    render(<PetCharacter stage="teen" {...HEALTHY} />)

    emit('loadError')

    expect(screen.queryByTestId('pet-lottie')).not.toBeInTheDocument()
    const placeholder = screen.getByTestId('pet-placeholder')
    expect(placeholder).toHaveAttribute('data-failed', 'true')
    expect(placeholder).toHaveAttribute('data-stage', 'teen')
  })

  it('labels the placeholder with a translated string', () => {
    render(
      <PetCharacter
        stage="baby"
        {...HEALTHY}
        labels={{ loading: 'Загружаем питомца' }}
      />,
    )

    emit('loadError')

    expect(screen.getByTestId('pet-placeholder')).toHaveAttribute(
      'aria-label',
      'Загружаем питомца',
    )
  })
})

describe('PetCharacter stroke interaction', () => {
  it.each<[string, (target: HTMLElement) => void]>([
    ['click', (target) => fireEvent.click(target)],
    ['Enter', (target) => fireEvent.keyDown(target, { key: 'Enter' })],
    ['Space', (target) => fireEvent.keyDown(target, { key: ' ' })],
  ])('strokes the raccoon on %s', (_name, interact) => {
    const onStroke = vi.fn()
    render(<PetCharacter stage="teen" {...HEALTHY} onStroke={onStroke} />)

    interact(screen.getByTestId('pet-lottie'))

    expect(onStroke).toHaveBeenCalledTimes(1)
  })

  it.each<PetStage>(ALL_STAGES)('is strokable on stage %s', (stage) => {
    const onStroke = vi.fn()
    render(<PetCharacter stage={stage} {...HEALTHY} onStroke={onStroke} />)

    fireEvent.click(screen.getByTestId('pet-lottie'))
    fireEvent.keyDown(screen.getByTestId('pet-lottie'), { key: 'Enter' })

    expect(onStroke).toHaveBeenCalledTimes(2)
  })

  it('ignores unrelated keys', () => {
    const onStroke = vi.fn()
    render(<PetCharacter stage="teen" {...HEALTHY} onStroke={onStroke} />)

    fireEvent.keyDown(screen.getByTestId('pet-lottie'), { key: 'a' })

    expect(onStroke).not.toHaveBeenCalled()
  })

  it('marks the raccoon as petted for visual feedback', () => {
    render(<PetCharacter stage="teen" {...HEALTHY} />)

    const character = screen.getByTestId('pet-lottie')
    fireEvent.click(character)

    expect(character.getAttribute('class')).toContain('pet-lottie--petted')
  })
})

describe('PetCharacter accessibility', () => {
  it.each<PetStage>(ALL_STAGES)(
    'is a focusable button with a descriptive label on stage %s',
    (stage) => {
      render(
        <PetCharacter
          stage={stage}
          {...HEALTHY}
          labels={{ name: 'Проша', strokeHint: 'нажмите, чтобы погладить' }}
          onStroke={vi.fn()}
        />,
      )

      const character = screen.getByRole('button')
      expect(character).toHaveAttribute('tabindex', '0')
      expect(character.getAttribute('aria-label')).toContain('Проша')
      expect(character.getAttribute('aria-label')).toContain(
        'нажмите, чтобы погладить',
      )
    },
  )
})
