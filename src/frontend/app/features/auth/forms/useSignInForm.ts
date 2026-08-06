import { useSignIn } from '#/features/auth/hooks'
import { SignInSchema } from '#/features/auth/schemas'
import type { SignInRequest } from '#/features/auth/types'
import { useForm } from '@tanstack/react-form'

const defaultValues: SignInRequest = {
  email: '',
  password: '',
}

export const useSignInForm = () => {
  const { mutateAsync } = useSignIn()

  return useForm({
    formId: 'signIn',
    defaultValues,
    validators: {
      onChange: SignInSchema,
    },
    onSubmit: ({ value }) => mutateAsync(value),
  })
}
