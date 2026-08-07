import { useForm } from '@tanstack/react-form'
import { useSignIn } from '#/features/auth/hooks'
import { SignInSchema } from '#/features/auth/schemas'
import type { SignInRequest } from '#/features/auth/types'
import type { Session } from '#/types'

const defaultValues: SignInRequest = {
  email: '',
  password: '',
}

interface UseSignInFormOptions {
  onSuccess?: (session: Session) => void
}

export const useSignInForm = ({ onSuccess }: UseSignInFormOptions = {}) => {
  const { mutateAsync, isPending, error, reset } = useSignIn()

  const form = useForm({
    formId: 'signIn',
    defaultValues,
    validators: {
      onChange: SignInSchema,
    },
    onSubmit: async ({ value }) => {
      reset()
      try {
        const session = await mutateAsync(value)
        onSuccess?.(session)
      } catch {
        return
      }
    },
  })

  return { form, isPending, error }
}
