import { useSignIn } from '#/features/auth/hooks/useSignIn'
import { SignInSchema } from '#/features/auth/schemas/sign-in.schema'
import type { SignInRequest } from '#/features/auth/types/sign-in.request.type'
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
