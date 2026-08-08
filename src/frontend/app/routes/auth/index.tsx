import { SignInForm } from '#/features/auth/components/SignInForm'
import { SignUpForm } from '#/features/auth/components/SignUpForm'
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from '#/shared/components/ui/card'
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '#/shared/components/ui/tabs'
import { useState } from 'react'

type AuthTab = 'login' | 'register'

export default function Index() {
  const [tab, setTab] = useState<AuthTab>('login')

  return (
    <div className="flex min-h-svh items-center justify-center bg-muted/40 px-4 py-8">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <CardTitle className="text-2xl font-semibold">
            {tab === 'login' ? 'Вход' : 'Регистрация'}
          </CardTitle>
          <CardDescription>
            {tab === 'register'
              ? 'Войдите, чтобы продолжить'
              : 'Создайте аккаунт, чтобы начать'}
          </CardDescription>
        </CardHeader>

        <Tabs value={tab} onValueChange={value => setTab(value as AuthTab)}>
          <div className="px-6">
            <TabsList className="grid w-full grid-cols-2">
              <TabsTrigger value="signin">Войти</TabsTrigger>
              <TabsTrigger value="signup">Регистрация</TabsTrigger>
            </TabsList>
          </div>

          <TabsContent value="signin">
            <SignInForm />
          </TabsContent>

          <TabsContent value="signup">
            <SignUpForm />
          </TabsContent>
        </Tabs>
      </Card>
    </div>
  )
}
