import { useSignUp } from '#/features/auth/hooks'
import { useForm } from '@tanstack/react-form'
import { SignUpSchema } from '#/features/auth/schemas'
import type { SignUpRequest } from '#/features/auth/types'

const defaultValues: SignUpRequest = {
  email: '',
  password: '',
  fullName: '',
}

export const useSignUpForm = () => {
  const { mutateAsync } = useSignUp()

  return useForm({
    formId: 'signUp',
    defaultValues,
    validators: {
      onChange: SignUpSchema,
    },
    onSubmit: ({ value }) => mutateAsync(value),
  })
}
