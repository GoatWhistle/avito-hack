import { ArrowLeft } from 'lucide-react'
import { Link } from 'react-router'
import { Button } from '#/components/ui'
import { cn } from '#/lib/utils'
import type { ReactNode } from 'react'

export const GAMES_PATH = '/pet/games'

interface GameScreenShellProps {
  title: string
  subtitle?: string
  backLabel: string
  help?: ReactNode
  wide?: boolean
  children: ReactNode
}

export function GameScreenShell({
  title,
  subtitle,
  backLabel,
  help,
  wide = false,
  children,
}: GameScreenShellProps) {
  return (
    <div
      className={cn(
        'mx-auto flex w-full max-w-6xl flex-col px-3 sm:px-6',
        wide ? 'gap-2 py-2' : 'gap-4 py-4',
      )}
    >
      <header className="flex items-center gap-3">
        <Button
          variant="ghost"
          size="icon-sm"
          aria-label={backLabel}
          render={<Link to={GAMES_PATH} />}
        >
          <ArrowLeft className="size-4" aria-hidden="true" />
        </Button>
        <div className="flex min-w-0 flex-col">
          <h1 className="truncate text-xl font-semibold">{title}</h1>
          {subtitle && (
            <p className="truncate text-xs text-muted-foreground">{subtitle}</p>
          )}
        </div>
      </header>

      {help ? (
        <div className="relative flex flex-col gap-4">
          <div
            className={cn(
              'flex min-w-0 flex-col gap-4',
              wide ? 'w-full' : 'xl:mx-auto xl:w-full xl:max-w-3xl',
            )}
          >
            {children}
          </div>
          <div className="flex flex-col gap-3 xl:absolute xl:end-0 xl:top-0 xl:w-64">
            {help}
          </div>
        </div>
      ) : (
        children
      )}
    </div>
  )
}
