import { useForm } from '@tanstack/react-form'
import { useSignUp } from '#/features/auth/hooks'
import { SignUpSchema } from '#/features/auth/schemas'
import type { SignUpRequest } from '#/features/auth/types'
import type { Session } from '#/types'

const defaultValues: SignUpRequest = {
  email: '',
  password: '',
  fullName: '',
}

interface UseSignUpFormOptions {
  onSuccess?: (session: Session) => void
}

export const useSignUpForm = ({ onSuccess }: UseSignUpFormOptions = {}) => {
  const { mutateAsync, isPending, error, reset } = useSignUp()

  const form = useForm({
    formId: 'signUp',
    defaultValues,
    validators: {
      onChange: SignUpSchema,
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
