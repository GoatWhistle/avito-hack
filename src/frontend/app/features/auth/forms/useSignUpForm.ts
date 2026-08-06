import { useSignUp } from '#/features/auth/hooks'
import { SignUpSchema } from '#/features/auth/schemas'
import type { SignUpRequest } from '#/features/auth/types'
import { useForm } from '@tanstack/react-form'

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
