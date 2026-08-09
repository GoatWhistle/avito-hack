import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { NavLink } from 'react-router'
import { Menu, X } from 'lucide-react'
import { Button, OverlayPortal } from '#/components/ui'
import { cn } from '#/lib/utils'
import { navItems } from '#/features/layout/nav-items'

export function MobileMenu() {
  const { t } = useTranslation('common')
  const [open, setOpen] = useState(false)

  useEffect(() => {
    if (!open) return

    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') setOpen(false)
    }
    document.addEventListener('keydown', onKeyDown)
    return () => document.removeEventListener('keydown', onKeyDown)
  }, [open])

  return (
    <div className="lg:hidden">
      <Button
        type="button"
        variant="ghost"
        size="icon-sm"
        onClick={() => setOpen((value) => !value)}
        aria-expanded={open}
        aria-label={open ? t('actions.close') : t('nav.menu')}
      >
        {open ? (
          <X className="size-4" aria-hidden="true" />
        ) : (
          <Menu className="size-4" aria-hidden="true" />
        )}
      </Button>

      {open && (
        <OverlayPortal>
          <div
            className="fixed inset-0 z-overlay bg-foreground/20"
            onClick={() => setOpen(false)}
            aria-hidden="true"
          />
          <nav
            aria-label={t('nav.primary')}
            className="fixed inset-x-0 top-14 z-modal mx-2 rounded-xl bg-popover p-2 text-popover-foreground shadow-lg ring-1 ring-foreground/10"
          >
            <ul className="flex flex-col gap-0.5">
              {navItems.map((item) => (
                <li key={item.key}>
                  <NavLink
                    to={item.to}
                    onClick={() => setOpen(false)}
                    className={({ isActive }) =>
                      cn(
                        'flex items-center gap-2 rounded-lg px-2.5 py-2 text-sm font-medium no-underline transition-colors',
                        isActive
                          ? 'bg-primary-subtle text-primary-subtle-foreground'
                          : 'hover:bg-muted',
                      )
                    }
                  >
                    <item.icon className="size-4" aria-hidden="true" />
                    {t(`nav.${item.key}`)}
                  </NavLink>
                </li>
              ))}
            </ul>
          </nav>
        </OverlayPortal>
      )}
    </div>
  )
}
