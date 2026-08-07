import { screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { UserMenu } from './UserMenu'
import { makeUser, renderWithShell } from '#/features/layout/test-utils'

const authed = (overrides = {}) => ({
  session: { user: makeUser(overrides), isAuthenticated: true },
})

const openMenu = async () => {
  const trigger = screen.getByRole('button', { expanded: false })
  await userEvent.click(trigger)

  return trigger
}

describe('UserMenu', () => {
  it('renders nothing without a signed in user', () => {
    const { container } = renderWithShell(<UserMenu />)

    expect(container).toBeEmptyDOMElement()
  })

  it('shows the initials of the user full name', () => {
    renderWithShell(<UserMenu />, authed())

    expect(screen.getByRole('button')).toHaveTextContent('ИИ')
  })

  it('falls back to a placeholder when the name is blank', () => {
    renderWithShell(<UserMenu />, authed({ fullName: '   ' }))

    expect(screen.getByRole('button')).toHaveTextContent('?')
  })

  it('uses a single initial for a one word name', () => {
    renderWithShell(<UserMenu />, authed({ fullName: 'Ноти' }))

    expect(screen.getByRole('button')).toHaveTextContent('Н')
  })

  it('opens a menu with the account details', async () => {
    renderWithShell(<UserMenu />, authed())
    const trigger = await openMenu()

    expect(trigger).toHaveAttribute('aria-expanded', 'true')
    const menu = screen.getByRole('menu')
    expect(menu).toHaveTextContent('Иван Иванов')
    expect(menu).toHaveTextContent('demo@example.com')
    expect(screen.getAllByRole('menuitem')).toHaveLength(2)
  })

  it('toggles closed when the trigger is pressed again', async () => {
    renderWithShell(<UserMenu />, authed())
    const trigger = await openMenu()

    await userEvent.click(trigger)

    expect(screen.queryByRole('menu')).not.toBeInTheDocument()
  })

  it('closes on Escape', async () => {
    renderWithShell(<UserMenu />, authed())
    await openMenu()

    await userEvent.keyboard('{Escape}')

    expect(screen.queryByRole('menu')).not.toBeInTheDocument()
  })

  it('stays open for unrelated keys', async () => {
    renderWithShell(<UserMenu />, authed())
    await openMenu()

    await userEvent.keyboard('{ArrowDown}')

    expect(screen.getByRole('menu')).toBeInTheDocument()
  })

  it('closes on a pointer press outside the menu', async () => {
    renderWithShell(<UserMenu />, authed())
    await openMenu()

    await userEvent.click(document.body)

    expect(screen.queryByRole('menu')).not.toBeInTheDocument()
  })

  it('stays open when clicking inside the menu container', async () => {
    renderWithShell(<UserMenu />, authed())
    await openMenu()

    await userEvent.click(screen.getByRole('menu'))

    expect(screen.getByRole('menu')).toBeInTheDocument()
  })

  it('closes when the profile link is chosen', async () => {
    renderWithShell(<UserMenu />, authed())
    await openMenu()

    const [profile] = screen.getAllByRole('menuitem')
    await userEvent.click(profile)

    expect(screen.queryByRole('menu')).not.toBeInTheDocument()
  })

  it('signs the user out and leaves the menu', async () => {
    const signOut = vi.fn()
    renderWithShell(<UserMenu />, {
      session: { user: makeUser(), isAuthenticated: true, signOut },
    })
    await openMenu()

    const [, signOutItem] = screen.getAllByRole('menuitem')
    await userEvent.click(signOutItem)

    expect(signOut).toHaveBeenCalledTimes(1)
    expect(screen.queryByRole('menu')).not.toBeInTheDocument()
  })
})
