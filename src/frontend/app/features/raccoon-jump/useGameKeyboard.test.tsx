import { useRef } from 'react'
import { act, render } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { Game } from './game'
import { useGameKeyboard } from './useGameKeyboard'

const press = (key: string, target: EventTarget = document.body) => {
  const event = new KeyboardEvent('keydown', {
    key,
    bubbles: true,
    cancelable: true,
  })

  act(() => {
    target.dispatchEvent(event)
  })

  return event
}

let game: Game
const onStart = vi.fn()

function Harness() {
  const gameRef = useRef<Game | null>(game)
  const surfaceRef = useRef<HTMLDivElement>(null)

  useGameKeyboard(gameRef, surfaceRef, onStart)

  return <div ref={surfaceRef} data-testid="surface" />
}

describe('useGameKeyboard', () => {
  beforeEach(() => {
    game = new Game()
    onStart.mockClear()
  })

  it('consumes arrow keys so the page does not scroll', () => {
    render(<Harness />)

    expect(press('ArrowLeft').defaultPrevented).toBe(true)
    expect(game.input.left).toBe(true)

    expect(press('ArrowRight').defaultPrevented).toBe(true)
    expect(game.input.right).toBe(true)
  })

  it('requests a server-backed start on space when the game is not playing', () => {
    render(<Harness />)

    expect(press(' ').defaultPrevented).toBe(true)
    expect(onStart).toHaveBeenCalledTimes(1)
    expect(game.status).toBe('idle')
  })

  it('does not request a start while a run is in progress', () => {
    game.reset(1)
    render(<Harness />)

    press(' ')

    expect(onStart).not.toHaveBeenCalled()
  })

  it('ignores keys aimed at a control outside the game surface', () => {
    render(<Harness />)

    const button = document.createElement('button')
    document.body.appendChild(button)
    button.focus()

    expect(press(' ', button).defaultPrevented).toBe(false)
    expect(onStart).not.toHaveBeenCalled()

    expect(press('ArrowLeft', button).defaultPrevented).toBe(false)
    expect(game.input.left).toBe(false)

    button.remove()
  })

  it('never swallows Tab or Escape', () => {
    render(<Harness />)

    expect(press('Tab').defaultPrevented).toBe(false)
    expect(press('Escape').defaultPrevented).toBe(false)
  })

  it('ignores shortcuts held with a modifier', () => {
    render(<Harness />)

    const event = new KeyboardEvent('keydown', {
      key: 'ArrowLeft',
      ctrlKey: true,
      bubbles: true,
      cancelable: true,
    })
    act(() => {
      document.body.dispatchEvent(event)
    })

    expect(event.defaultPrevented).toBe(false)
    expect(game.input.left).toBe(false)
  })

  it('releases movement when the window loses focus', () => {
    render(<Harness />)

    press('ArrowLeft')
    expect(game.input.left).toBe(true)

    act(() => {
      window.dispatchEvent(new Event('blur'))
    })

    expect(game.input.left).toBe(false)
  })

  it('stops responding after unmount', () => {
    const { unmount } = render(<Harness />)
    unmount()

    expect(press('ArrowLeft').defaultPrevented).toBe(false)
    expect(game.input.left).toBe(false)
  })
})
