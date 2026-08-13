import { useEffect, useState } from 'react'
import { Collapsible } from '@base-ui/react/collapsible'
import { HelpCircle, X } from 'lucide-react'
import { Button } from '#/components/ui'
import { cn } from '#/lib/utils'
import type { ReactNode } from 'react'

export interface GameHelpRow {
  tone: string
  sample: ReactNode
  label: string
}

export interface GameHelpProps {
  storageKey: string
  title: string
  intro: string
  rows?: GameHelpRow[]
  outro?: string
  closeLabel: string
  openLabel: string
  hideLabel: string
}

const storagePath = (storageKey: string) =>
  `avito-hack.${storageKey}.help-dismissed`

export function GameHelp({
  storageKey,
  title,
  intro,
  rows = [],
  outro,
  closeLabel,
  openLabel,
  hideLabel,
}: GameHelpProps) {
  const [open, setOpen] = useState(false)
  const [hydrated, setHydrated] = useState(false)

  useEffect(() => {
    if (hydrated) return

    setOpen(window.localStorage.getItem(storagePath(storageKey)) !== '1')
    setHydrated(true)
  }, [storageKey, hydrated])

  const onOpenChange = (next: boolean) => {
    setOpen(next)

    if (next) window.localStorage.removeItem(storagePath(storageKey))
    else window.localStorage.setItem(storagePath(storageKey), '1')
  }

  return (
    <Collapsible.Root
      open={open}
      onOpenChange={onOpenChange}
      className="flex flex-col"
    >
      {!open && (
        <Collapsible.Trigger
          render={
            <Button
              type="button"
              variant="ghost"
              size="sm"
              className="self-end"
            >
              <HelpCircle className="size-4" aria-hidden="true" />
              {openLabel}
            </Button>
          }
        />
      )}

      <Collapsible.Panel
        className="game-help-panel"
        render={
          <aside className="relative flex flex-col gap-3 rounded-2xl bg-card p-4 ring-1 ring-border" />
        }
      >
        <Collapsible.Trigger
          render={
            <Button
              type="button"
              variant="ghost"
              size="icon-sm"
              aria-label={hideLabel}
              className="absolute end-2 top-2"
            >
              <X className="size-4" aria-hidden="true" />
            </Button>
          }
        />

        <div className="flex flex-col gap-1 pe-8">
          <h2 className="text-sm font-semibold">{title}</h2>
          <p className="text-xs text-muted-foreground">{intro}</p>
        </div>

        {rows.length > 0 && (
          <ul className="flex list-none flex-col gap-2">
            {rows.map((row) => (
              <li key={row.label} className="flex items-center gap-2.5">
                <span
                  aria-hidden="true"
                  className={cn(
                    'flex size-8 shrink-0 items-center justify-center rounded-lg text-sm font-bold uppercase',
                    row.tone,
                  )}
                >
                  {row.sample}
                </span>
                <span className="text-xs text-muted-foreground">
                  {row.label}
                </span>
              </li>
            ))}
          </ul>
        )}

        {outro && <p className="text-xs text-muted-foreground">{outro}</p>}

        <Collapsible.Trigger
          render={
            <Button type="button" size="sm">
              {closeLabel}
            </Button>
          }
        />
      </Collapsible.Panel>
    </Collapsible.Root>
  )
}
