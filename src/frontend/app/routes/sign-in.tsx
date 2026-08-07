import { AuthScreen, SignInForm } from '#/features/auth/components'

export default function SignInRoute() {
  return (
    <AuthScreen mode="signIn">
      <SignInForm />
    </AuthScreen>
  )
}
