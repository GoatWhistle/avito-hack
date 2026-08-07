import { useTranslation } from 'react-i18next'
import { Moon, Sun } from 'lucide-react'
import { Button } from '#/components/ui'
import { useTheme } from '#/features/layout/theme'

export function ThemeToggle() {
  const { t } = useTranslation('common')
  const { resolved, toggle } = useTheme()

  const nextLabel = resolved === 'dark' ? t('theme.light') : t('theme.dark')

  return (
    <Button
      type="button"
      variant="ghost"
      size="icon-sm"
      onClick={toggle}
      aria-label={`${t('theme.label')}: ${nextLabel}`}
      title={nextLabel}
    >
      {resolved === 'dark' ? (
        <Sun className="size-4" aria-hidden="true" />
      ) : (
        <Moon className="size-4" aria-hidden="true" />
      )}
    </Button>
  )
}
