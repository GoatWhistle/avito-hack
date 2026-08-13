import { Link } from 'react-router'
import { Check, ChevronRight, type LucideIcon } from 'lucide-react'

interface GameListRowProps {
  to: string
  icon: LucideIcon
  name: string
  description: string
  done?: boolean
}

export function GameListRow({
  to,
  icon: Icon,
  name,
  description,
  done = false,
}: GameListRowProps) {
  return (
    <Link
      to={to}
      className="group flex items-center gap-3 rounded-xl bg-card p-3 no-underline ring-1 ring-border transition-colors hover:bg-muted"
    >
      <span className="flex size-9 shrink-0 items-center justify-center rounded-lg bg-primary-subtle">
        <Icon
          className="size-5 text-primary-subtle-foreground"
          aria-hidden="true"
        />
      </span>

      <span className="flex min-w-0 flex-1 flex-col">
        <span className="text-sm font-medium">{name}</span>
        <span className="truncate text-xs text-muted-foreground">
          {description}
        </span>
      </span>

      <span className="flex shrink-0 items-center gap-2">
        {done && <Check className="size-4 text-primary" aria-hidden="true" />}
        <ChevronRight
          className="size-4 text-muted-foreground transition-transform group-hover:translate-x-0.5"
          aria-hidden="true"
        />
      </span>
    </Link>
  )
}
