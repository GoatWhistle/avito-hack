import { useCallback, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Button, OverlayPortal } from '#/components/ui'
import { useReducedMotion } from '#/lib/useReducedMotion'
import { cn } from '#/lib/utils'
import { tourSteps } from './tour-steps'
import { useFocusTrap } from './useFocusTrap'

import './platform-tour.css'

interface PlatformTourProps {
  onFinish: () => void
  onSkip: () => void
}

export function PlatformTour({ onFinish, onSkip }: PlatformTourProps) {
  const { t } = useTranslation('onboarding')
  const reducedMotion = useReducedMotion()
  const [index, setIndex] = useState(0)

  const skip = useCallback(() => onSkip(), [onSkip])
  const dialogRef = useFocusTrap(skip)

  const step = tourSteps[index]
  const Icon = step.icon
  const isFirst = index === 0
  const isLast = index === tourSteps.length - 1

  return (
    <OverlayPortal>
      <div className="fixed inset-0 z-modal flex items-center justify-center bg-foreground/40 p-4 backdrop-blur-sm">
        <div
          ref={dialogRef}
          role="dialog"
          aria-modal="true"
          aria-labelledby="platform-tour-title"
          aria-describedby="platform-tour-description"
          tabIndex={-1}
          data-testid="platform-tour"
          className="w-full max-w-md overflow-hidden rounded-xl bg-card p-6 text-card-foreground shadow-lg ring-1 ring-foreground/10 outline-none"
        >
          <div className="flex items-center justify-between gap-3">
            <p
              className="text-xs font-medium text-muted-foreground"
              data-testid="platform-tour-counter"
            >
              {t('tour.counter', {
                current: index + 1,
                total: tourSteps.length,
              })}
            </p>
            <Button variant="ghost" size="sm" onClick={skip}>
              {t('tour.skip')}
            </Button>
          </div>

          <div
            key={step.key}
            className={cn(
              'mt-4 flex flex-col items-center text-center',
              !reducedMotion && 'platform-tour-step',
            )}
          >
            <span
              className="flex size-12 items-center justify-center rounded-full bg-primary/10 text-primary"
              aria-hidden="true"
            >
              <Icon className="size-6" />
            </span>

            <h2
              id="platform-tour-title"
              className="mt-4 font-heading text-lg font-semibold"
            >
              {t(`tour.steps.${step.key}.title`)}
            </h2>
            <p
              id="platform-tour-description"
              className="mt-2 text-sm text-muted-foreground"
            >
              {t(`tour.steps.${step.key}.description`)}
            </p>
          </div>

          <ol
            className="mt-5 flex justify-center gap-1.5"
            aria-label={t('tour.progressLabel')}
          >
            {tourSteps.map((item, position) => (
              <li
                key={item.key}
                aria-current={position === index ? 'step' : undefined}
                className={cn(
                  'h-1.5 rounded-full transition-all',
                  position === index
                    ? 'w-6 bg-primary'
                    : 'w-1.5 bg-muted-foreground/30',
                )}
              >
                <span className="sr-only">
                  {t(`tour.steps.${item.key}.title`)}
                </span>
              </li>
            ))}
          </ol>

          <div className="mt-6 flex items-center gap-2">
            <Button
              variant="outline"
              className="flex-1"
              disabled={isFirst}
              onClick={() => setIndex((value) => Math.max(0, value - 1))}
            >
              {t('tour.back')}
            </Button>
            <Button
              className="flex-1"
              onClick={() =>
                isLast ? onFinish() : setIndex((value) => value + 1)
              }
            >
              {isLast ? t('tour.finish') : t('tour.next')}
            </Button>
          </div>
        </div>
      </div>
    </OverlayPortal>
  )
}
