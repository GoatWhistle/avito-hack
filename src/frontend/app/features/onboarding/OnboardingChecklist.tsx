import { useTranslation } from 'react-i18next'
import { Link } from 'react-router'
import { Check, Circle } from 'lucide-react'
import { cn } from '#/lib/utils'

export interface ChecklistStep {
  key: 'publish' | 'comeback' | 'meet'
  done: boolean
  to?: string
}

interface OnboardingChecklistProps {
  steps: ChecklistStep[]
}

export function OnboardingChecklist({ steps }: OnboardingChecklistProps) {
  const { t } = useTranslation('common')

  return (
    <ol className="flex flex-col gap-2" data-testid="onboarding-checklist">
      {steps.map((step) => {
        const content = (
          <>
            <span
              className={cn(
                'flex size-5 shrink-0 items-center justify-center rounded-full',
                step.done
                  ? 'bg-success text-success-foreground'
                  : 'bg-muted text-muted-foreground',
              )}
              aria-hidden="true"
            >
              {step.done ? (
                <Check className="size-3" />
              ) : (
                <Circle className="size-2" />
              )}
            </span>
            <span className="min-w-0">
              <span
                className={cn(
                  'block text-sm font-medium',
                  step.done && 'text-muted-foreground line-through',
                )}
              >
                {t(`onboarding.steps.${step.key}.title`)}
              </span>
              <span className="block text-xs text-muted-foreground">
                {t(`onboarding.steps.${step.key}.hint`)}
              </span>
            </span>
          </>
        )

        return (
          <li key={step.key}>
            {step.to && !step.done ? (
              <Link
                to={step.to}
                className="flex items-start gap-2.5 rounded-lg p-2 no-underline transition-colors hover:bg-muted"
              >
                {content}
              </Link>
            ) : (
              <div className="flex items-start gap-2.5 p-2">{content}</div>
            )}
            <span className="sr-only">
              {step.done
                ? t('onboarding.stepDone')
                : t('onboarding.stepPending')}
            </span>
          </li>
        )
      })}
    </ol>
  )
}
