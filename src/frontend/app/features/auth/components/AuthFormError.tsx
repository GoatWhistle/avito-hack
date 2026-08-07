import { AlertCircle } from 'lucide-react'

interface AuthFormErrorProps {
  message: string | null
}

export function AuthFormError({ message }: AuthFormErrorProps) {
  return (
    <div role="alert" aria-live="polite" data-testid="auth-form-error">
      {message && (
        <p className="flex items-start gap-2 rounded-lg bg-destructive-subtle px-3 py-2 text-sm text-destructive-subtle-foreground">
          <AlertCircle className="mt-0.5 size-4 shrink-0" aria-hidden="true" />
          <span className="min-w-0">{message}</span>
        </p>
      )}
    </div>
  )
}
