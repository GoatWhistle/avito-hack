import { useSignUp } from '#/features/auth/hooks/useSignUp'
import { SignUpSchema } from '#/features/auth/schemas/sign-up.schema'
import type { SignUpRequest } from '#/features/auth/types/sign-up.request.type'
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
