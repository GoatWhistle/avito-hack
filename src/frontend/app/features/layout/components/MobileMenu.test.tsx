import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it } from 'vitest'
import { MobileMenu } from './MobileMenu'
import { navItems } from '#/features/layout/nav-items'
import { renderWithShell } from '#/features/layout/test-utils'

const openMenu = async () => {
  const trigger = screen.getByRole('button', { expanded: false })
  await userEvent.click(trigger)

  return trigger
}

describe('MobileMenu', () => {
  it('starts collapsed and exposes the toggle state', () => {
    renderWithShell(<MobileMenu />)

    const trigger = screen.getByRole('button')
    expect(trigger).toHaveAttribute('aria-expanded', 'false')
    expect(screen.queryByRole('navigation')).not.toBeInTheDocument()
  })

  it('opens the panel with every navigation entry', async () => {
    renderWithShell(<MobileMenu />)
    const trigger = await openMenu()

    expect(trigger).toHaveAttribute('aria-expanded', 'true')
    const nav = screen.getByRole('navigation')
    expect(screen.getAllByRole('link', { hidden: false }).length).toBe(
      navItems.length,
    )
    expect(nav).toBeInTheDocument()
  })

  it('collapses again when the trigger is pressed twice', async () => {
    renderWithShell(<MobileMenu />)
    const trigger = await openMenu()

    await userEvent.click(trigger)

    expect(trigger).toHaveAttribute('aria-expanded', 'false')
    expect(screen.queryByRole('navigation')).not.toBeInTheDocument()
  })

  it('closes when the backdrop is clicked', async () => {
    const { container } = renderWithShell(<MobileMenu />)
    await openMenu()

    const backdrop = container.querySelector('[aria-hidden="true"].fixed')
    expect(backdrop).not.toBeNull()
    await userEvent.click(backdrop as Element)

    expect(screen.queryByRole('navigation')).not.toBeInTheDocument()
  })

  it('closes on Escape', async () => {
    renderWithShell(<MobileMenu />)
    await openMenu()

    await userEvent.keyboard('{Escape}')

    expect(screen.queryByRole('navigation')).not.toBeInTheDocument()
  })

  it('ignores unrelated keys while open', async () => {
    renderWithShell(<MobileMenu />)
    await openMenu()

    await userEvent.keyboard('{ArrowDown}')

    expect(screen.getByRole('navigation')).toBeInTheDocument()
  })

  it('closes after following a navigation link', async () => {
    renderWithShell(<MobileMenu />)
    await openMenu()

    const [first] = screen.getAllByRole('link')
    await userEvent.click(first)

    expect(screen.queryByRole('navigation')).not.toBeInTheDocument()
  })

  it('marks the active route link inside the panel', async () => {
    renderWithShell(<MobileMenu />, { initialEntries: ['/pet'] })
    await openMenu()

    const active = screen
      .getAllByRole('link')
      .find((link) => link.getAttribute('href') === '/pet')

    expect(active).toHaveAttribute('aria-current', 'page')
  })
})
