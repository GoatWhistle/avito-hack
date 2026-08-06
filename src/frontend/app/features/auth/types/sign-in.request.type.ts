import type { SignUpRequest } from './sign-up.request.type'

export type SignInRequest = Omit<SignUpRequest, 'fullName'>
