import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router'
import { LogOut } from 'lucide-react'
import {
  Button,
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  Separator,
} from '#/components/ui'
import { useSession } from '#/features/auth/session'
import { ProfileForm } from './ProfileForm'
import { ProfileStats } from './ProfileStats'

const formatDate = (value: string, locale: string) => {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value

  return new Intl.DateTimeFormat(locale, { dateStyle: 'long' }).format(date)
}

export function ProfileScreen() {
  const { t, i18n } = useTranslation(['common', 'auth'])
  const { user, signOut } = useSession()
  const navigate = useNavigate()

  if (!user) return null

  const handleSignOut = () => {
    signOut()
    void navigate('/sign-in', { replace: true })
  }

  return (
    <main
      className="mx-auto flex w-full max-w-prose flex-col gap-4"
      data-testid="profile-screen"
    >
      <h1 className="font-heading text-xl font-semibold">
        {t('common:nav.profile')}
      </h1>

      <Card>
        <CardHeader>
          <CardTitle>{t('common:profile.account')}</CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col gap-3">
          <ProfileForm user={user} />
          <Separator />
          <div>
            <p className="text-xs text-muted-foreground">
              {t('auth:fields.email')}
            </p>
            <p className="truncate text-sm font-medium">{user.email}</p>
          </div>
          <div>
            <p className="text-xs text-muted-foreground">
              {t('common:profile.memberSince')}
            </p>
            <p className="text-sm font-medium">
              {formatDate(user.createdAt, i18n.language)}
            </p>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>{t('common:profile.stats')}</CardTitle>
        </CardHeader>
        <CardContent>
          <ProfileStats />
        </CardContent>
      </Card>

      <Button
        type="button"
        variant="outline"
        onClick={handleSignOut}
        className="self-start text-destructive"
      >
        <LogOut className="size-4" aria-hidden="true" />
        {t('auth:signOut')}
      </Button>
    </main>
  )
}
