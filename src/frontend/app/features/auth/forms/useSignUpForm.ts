import { useForm } from '@tanstack/react-form'
import { useSignUp } from '#/features/auth/hooks'
import { SignUpSchema } from '#/features/auth/schemas'
import type { SignUpRequest } from '#/features/auth/types'
import type { Session } from '#/types'
import { useAuthFieldFocus } from './useAuthFieldFocus'

const defaultValues: SignUpRequest = {
  email: '',
  password: '',
  fullName: '',
}

const fieldOrder = ['fullName', 'email', 'password'] as const

interface UseSignUpFormOptions {
  onSuccess?: (session: Session) => void
}

export const useSignUpForm = ({ onSuccess }: UseSignUpFormOptions = {}) => {
  const { mutateAsync, isPending, error, reset } = useSignUp()
  const { register, focusFirstInvalid } = useAuthFieldFocus(fieldOrder)

  const form = useForm({
    formId: 'signUp',
    defaultValues,
    validators: {
      onChange: SignUpSchema,
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
