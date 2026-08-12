import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { ServerCrash } from 'lucide-react'
import { Button } from '#/components/ui'

export function CrashScreen() {
  const { t } = useTranslation('errors')

  return (
    <div className="mx-auto flex max-w-content flex-1 items-center justify-center px-4 py-10 sm:py-16">
      <div className="flex w-full max-w-md flex-col items-center gap-6 text-center">
        <div className="flex size-20 items-center justify-center rounded-full bg-destructive/10 text-destructive ring-1 ring-foreground/10">
          <ServerCrash className="size-9" aria-hidden="true" />
        </div>

        <div className="flex flex-col gap-2">
          <h1 className="text-xl font-semibold text-balance">
            {t('page.crashTitle')}
          </h1>
          <p className="max-w-prose text-sm text-pretty text-muted-foreground">
            {t('page.crashHint')}
          </p>
        </div>

        <Button render={<Link to="/" />}>{t('page.goHome')}</Button>
      </div>
    </div>
  )
}
