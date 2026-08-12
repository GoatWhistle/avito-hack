import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { PawPrint, SearchX } from 'lucide-react'
import { Button } from '#/components/ui'

const pawTrail = [
  { top: '18%', left: '12%', rotate: '-18deg', size: 'size-3.5' },
  { top: '32%', left: '22%', rotate: '10deg', size: 'size-4' },
  { top: '24%', left: '78%', rotate: '24deg', size: 'size-3.5' },
  { top: '58%', left: '86%', rotate: '-6deg', size: 'size-4' },
  { top: '74%', left: '16%', rotate: '14deg', size: 'size-3' },
]

export function NotFoundScreen() {
  const { t } = useTranslation('errors')

  return (
    <div className="mx-auto flex max-w-content flex-1 items-center justify-center px-4 py-10 sm:py-16">
      <div className="flex w-full max-w-md flex-col items-center gap-6 text-center">
        <div
          className="relative flex aspect-square w-full max-w-64 items-center justify-center rounded-2xl border-2 border-dashed border-foreground/15 bg-card"
          aria-hidden="true"
        >
          {pawTrail.map((paw, index) => (
            <PawPrint
              key={index}
              className={`absolute text-muted-foreground/30 ${paw.size}`}
              style={{
                top: paw.top,
                left: paw.left,
                transform: `rotate(${paw.rotate})`,
              }}
            />
          ))}
          <div className="relative flex size-20 items-center justify-center rounded-full bg-primary-subtle text-primary-subtle-foreground ring-1 ring-foreground/10">
            <SearchX className="size-9" />
          </div>
        </div>

        <div className="flex flex-col gap-2">
          <h1 className="text-xl font-semibold text-balance">
            {t('page.notFoundTitle')}
          </h1>
          <p className="max-w-prose text-sm text-pretty text-muted-foreground">
            {t('page.notFoundHint')}
          </p>
        </div>

        <Button render={<Link to="/" />}>{t('page.goHome')}</Button>
      </div>
    </div>
  )
}
