import { useEffect, useRef, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Link, useNavigate } from 'react-router'
import { LogOut, User as UserIcon } from 'lucide-react'
import { Button } from '#/components/ui'
import { useSession } from '#/features/auth/session'

const initialsOf = (fullName: string) =>
  fullName
    .split(/\s+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase() ?? '')
    .join('') || '?'

export function UserMenu() {
  const { t } = useTranslation(['common', 'auth'])
  const { user, signOut } = useSession()
  const navigate = useNavigate()
  const [open, setOpen] = useState(false)
  const containerRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!open) return

    const onPointerDown = (event: MouseEvent) => {
      if (!containerRef.current?.contains(event.target as Node)) setOpen(false)
    }
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') setOpen(false)
    }

    document.addEventListener('mousedown', onPointerDown)
    document.addEventListener('keydown', onKeyDown)
    return () => {
      document.removeEventListener('mousedown', onPointerDown)
      document.removeEventListener('keydown', onKeyDown)
    }
  }, [open])

  if (!user) return null

  const handleSignOut = () => {
    setOpen(false)
    signOut()
    void navigate('/sign-in')
  }

  return (
    <div className="relative" ref={containerRef}>
      <Button
        type="button"
        variant="ghost"
        size="icon-sm"
        onClick={() => setOpen((value) => !value)}
        aria-haspopup="menu"
        aria-expanded={open}
        aria-label={t('common:nav.profile')}
      >
        <span className="flex size-6 items-center justify-center rounded-full bg-primary-subtle text-[0.6875rem] font-semibold text-primary-subtle-foreground">
          {initialsOf(user.fullName)}
        </span>
      </Button>

      {open && (
        <div
          role="menu"
          aria-label={t('common:nav.profile')}
          className="absolute end-0 z-dropdown mt-2 w-56 overflow-hidden rounded-xl bg-popover p-1 text-popover-foreground shadow-lg ring-1 ring-foreground/10"
        >
          <div className="px-2.5 py-2">
            <p className="truncate text-sm font-medium">{user.fullName}</p>
            <p className="truncate text-xs text-muted-foreground">
              {user.email}
            </p>
          </div>
          <div className="my-1 h-px bg-border" />
          <Link
            to="/profile"
            role="menuitem"
            onClick={() => setOpen(false)}
            className="flex items-center gap-2 rounded-lg px-2.5 py-1.5 text-sm no-underline hover:bg-muted focus-visible:bg-muted"
          >
            <UserIcon className="size-4" aria-hidden="true" />
            {t('common:nav.profile')}
          </Link>
          <button
            type="button"
            role="menuitem"
            onClick={handleSignOut}
            className="flex w-full items-center gap-2 rounded-lg px-2.5 py-1.5 text-start text-sm text-destructive hover:bg-destructive-subtle focus-visible:bg-destructive-subtle"
          >
            <LogOut className="size-4" aria-hidden="true" />
            {t('auth:signOut')}
          </button>
        </div>
      )}
    </div>
  )
}
