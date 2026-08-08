import { useSignIn } from '#/features/auth/hooks/useSignIn'
import { SignInSchema } from '#/features/auth/schemas/sign-in.schema'
import type { SignInRequest } from '#/features/auth/types/sign-in.request'
import { useForm } from '@tanstack/react-form'
import { useNavigate } from 'react-router'
import { toast } from 'sonner'

const defaultValues: SignInRequest = {
  email: '',
  password: '',
}

export const useSignInForm = () => {
  const navigate = useNavigate()
  const { mutateAsync } = useSignIn()

  return useForm({
    formId: 'signIn',
    defaultValues,
    validators: {
      onChange: SignInSchema,
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
