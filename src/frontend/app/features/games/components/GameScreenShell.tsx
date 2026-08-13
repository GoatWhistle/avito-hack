import { ArrowLeft } from 'lucide-react'
import { Link } from 'react-router'
import { Button } from '#/components/ui'
import type { ReactNode } from 'react'

export const GAMES_PATH = '/pet/games'

interface GameScreenShellProps {
  title: string
  subtitle?: string
  backLabel: string
  help?: ReactNode
  children: ReactNode
}

export function GameScreenShell({
  title,
  subtitle,
  backLabel,
  help,
  children,
}: GameScreenShellProps) {
  return (
    <div className="mx-auto flex w-full max-w-6xl flex-col gap-4 px-3 py-4 sm:px-6">
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
          <div className="flex min-w-0 flex-col gap-4 xl:mx-auto xl:w-full xl:max-w-3xl">
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
