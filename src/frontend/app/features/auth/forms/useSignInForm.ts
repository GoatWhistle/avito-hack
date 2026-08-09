import { useForm } from '@tanstack/react-form'
import { useSignIn } from '#/features/auth/hooks'
import { SignInSchema } from '#/features/auth/schemas'
import type { SignInRequest } from '#/features/auth/types'
import type { Session } from '#/types'
import { useAuthFieldFocus } from './useAuthFieldFocus'

const defaultValues: SignInRequest = {
  email: '',
  password: '',
}

const fieldOrder = ['email', 'password'] as const

interface UseSignInFormOptions {
  onSuccess?: (session: Session) => void
}

export const useSignInForm = ({ onSuccess }: UseSignInFormOptions = {}) => {
  const { mutateAsync, isPending, error, reset } = useSignIn()
  const { register, focusFirstInvalid } = useAuthFieldFocus(fieldOrder)

  const form = useForm({
    formId: 'signIn',
    defaultValues,
    validators: {
      onChange: SignInSchema,
    },
    onSubmitInvalid: ({ formApi }) => {
      const invalid = fieldOrder.filter(
        (name) => formApi.getFieldMeta(name)?.errors.length,
      )
      focusFirstInvalid(invalid)
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

  return { form, isPending, error, register, focusFirstInvalid }
}
