import { cn } from '#/lib/utils'
import type { BukovkiLetter } from '#/features/games/types'

interface BukovkiGridProps {
  wordLength: number
  maxTries: number
  history: BukovkiLetter[][]
  draft: string
  active?: boolean
}

export function BukovkiGrid({
  wordLength,
  maxTries,
  history,
  draft,
  active = true,
}: BukovkiGridProps) {
  const draftLetters = [...draft]
  const rows = Array.from({ length: maxTries }, (_, rowIndex) => {
    if (rowIndex < history.length) return history[rowIndex]
    if (rowIndex === history.length) {
      return Array.from({ length: wordLength }, (_, index) => ({
        char: draftLetters[index] ?? '',
        status: null,
      }))
    }

    return Array.from({ length: wordLength }, () => ({ char: '', status: null }))
  })

  return (
    <div className="flex flex-col items-center gap-1.5 sm:gap-2">
      {rows.map((row, rowIndex) => {
        const isCurrent = active && rowIndex === history.length

        return (
          <div key={rowIndex} className="flex gap-1.5 sm:gap-2">
            {row.map((cell, cellIndex) => (
              <span
                key={cellIndex}
                className={cn(
                  'flex size-10 items-center justify-center rounded-lg text-base font-bold uppercase transition-all sm:size-13 sm:text-xl',
                  cell.status === 'correct' &&
                    'bg-success text-success-foreground',
                  cell.status === 'present' &&
                    'bg-warning text-warning-foreground',
                  cell.status === 'absent' &&
                    'bg-destructive/85 text-destructive-foreground',
                  cell.status === null &&
                    (cell.char
                      ? 'scale-105 bg-card text-foreground ring-2 ring-primary/50'
                      : isCurrent && cellIndex === draftLetters.length
                        ? 'bg-card text-foreground ring-2 ring-primary/30'
                        : 'bg-card text-foreground ring-1 ring-border'),
                )}
              >
                {cell.char}
              </span>
            ))}
          </div>
        )
      })}
    </div>
  )
}
