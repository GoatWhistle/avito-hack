export interface PetActionsProps {
  strokeError: string | null
}

export function PetActions({ strokeError }: PetActionsProps) {
  return (
    <div className="flex flex-col gap-2">
      <p
        aria-live="polite"
        role="status"
        className="min-h-4 text-center text-xs text-destructive"
      >
        {strokeError ?? ''}
      </p>
    </div>
  )
}
