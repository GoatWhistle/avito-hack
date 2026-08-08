import { useSignUp } from '#/features/auth/hooks/useSignUp'
import { SignUpSchema } from '#/features/auth/schemas/sign-up.schema'
import type { SignUpRequest } from '#/features/auth/types/sign-up.request'
import { useForm } from '@tanstack/react-form'
import { useNavigate } from 'react-router'
import { toast } from 'sonner'

const defaultValues: SignUpRequest = {
  email: '',
  password: '',
  full_name: '',
}

export const useSignUpForm = () => {
  const navigate = useNavigate()
  const { mutateAsync } = useSignUp()

  return useForm({
    formId: 'signUp',
    defaultValues,
    validators: {
      onChange: SignUpSchema,
    },
    onSubmit: async ({ value }) => {
      try {
        const response = await mutateAsync(value)

        navigate('/raccoon')

        return response
      } catch {
        toast.error('Ошибка входа')
      }
    },
  })
}
