import { AuthScreen, SignUpForm } from '#/features/auth/components'

export default function SignUpRoute() {
  return (
    <AuthScreen mode="signUp">
      <SignUpForm />
    </AuthScreen>
  )
}
