import { Delete, CornerDownLeft } from 'lucide-react'
import { cn } from '#/lib/utils'
import type {
  BukovkiLetter,
  BukovkiLetterStatus,
} from '#/features/games/types'

const ROWS = [
  [...'йцукенгшщзхъ'],
  [...'фывапролджэ'],
  [...'ячсмитьбю'],
]

const RANK: Record<BukovkiLetterStatus, number> = {
  absent: 0,
  present: 1,
  correct: 2,
}

interface BukovkiKeyboardProps {
  history: BukovkiLetter[][]
  disabled: boolean
  canSubmit: boolean
  onKey: (char: string) => void
  onBackspace: () => void
  onSubmit: () => void
  submitLabel: string
  backspaceLabel: string
}

export function BukovkiKeyboard({
  history,
  disabled,
  canSubmit,
  onKey,
  onBackspace,
  onSubmit,
  submitLabel,
  backspaceLabel,
}: BukovkiKeyboardProps) {
  const statuses = new Map<string, BukovkiLetterStatus>()

  for (const row of history) {
    for (const cell of row) {
      const char = cell.char.toLowerCase()
      const known = statuses.get(char)
      if (!known || RANK[cell.status] > RANK[known]) {
        statuses.set(char, cell.status)
      }
    }
  }

  return (
    <div className="mx-auto flex w-full max-w-lg flex-col items-center gap-1.5">
      {ROWS.map((row, index) => (
        <div key={index} className="flex w-full justify-center gap-1 sm:gap-1.5">
          {index === 2 && (
            <KeyButton
              wide
              disabled={disabled || !canSubmit}
              label={submitLabel}
              onClick={onSubmit}
            >
              <CornerDownLeft className="size-4" aria-hidden="true" />
            </KeyButton>
          )}

          {row.map((char) => (
            <KeyButton
              key={char}
              disabled={disabled}
              status={statuses.get(char)}
              label={char}
              onClick={() => onKey(char)}
            >
              {char}
            </KeyButton>
          ))}

          {index === 2 && (
            <KeyButton
              wide
              disabled={disabled}
              label={backspaceLabel}
              onClick={onBackspace}
            >
              <Delete className="size-4" aria-hidden="true" />
            </KeyButton>
          )}
        </div>
      ))}
    </div>
  )
}

interface KeyButtonProps {
  children: React.ReactNode
  label: string
  disabled: boolean
  status?: BukovkiLetterStatus
  wide?: boolean
  onClick: () => void
}

function KeyButton({
  children,
  label,
  disabled,
  status,
  wide = false,
  onClick,
}: KeyButtonProps) {
  return (
    <button
      type="button"
      aria-label={label}
      disabled={disabled}
      onClick={onClick}
      className={cn(
        'flex h-12 min-w-0 flex-1 items-center justify-center rounded-lg text-sm font-semibold uppercase transition-colors select-none sm:h-13 sm:text-base',
        'focus-visible:ring-3 focus-visible:ring-ring/50 focus-visible:outline-none',
        'disabled:opacity-40',
        wide ? 'max-w-16 px-2' : 'max-w-10 sm:max-w-11',
        status === 'correct' && 'bg-success text-success-foreground',
        status === 'present' && 'bg-warning text-warning-foreground',
        status === 'absent' &&
          'bg-destructive/70 text-destructive-foreground opacity-80',
        !status && 'bg-card text-foreground ring-1 ring-border hover:bg-muted',
      )}
    >
      {children}
    </button>
  )
}
