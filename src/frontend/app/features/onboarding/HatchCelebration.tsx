import { useEffect, useRef } from 'react'
import { useTranslation } from 'react-i18next'
import { Button } from '#/components/ui'
import { PetAvatar } from '#/components/pet-avatar'
import { useReducedMotion } from '#/components/pet-avatar/hooks/useReducedMotion'

import './confetti.css'

interface HatchCelebrationProps {
  onClose: () => void
}

const confettiPieces = Array.from({ length: 14 }, (_, index) => index)

export function HatchCelebration({ onClose }: HatchCelebrationProps) {
  const { t } = useTranslation('pet')
  const reducedMotion = useReducedMotion()
  const closeRef = useRef<HTMLButtonElement>(null)

  useEffect(() => {
    closeRef.current?.focus()

    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') onClose()
    }
    document.addEventListener('keydown', onKeyDown)
    return () => document.removeEventListener('keydown', onKeyDown)
  }, [onClose])

  return (
    <div
      role="dialog"
      aria-modal="true"
      aria-labelledby="hatch-title"
      className="fixed inset-0 z-modal flex items-center justify-center bg-foreground/40 p-4 backdrop-blur-sm"
    >
      <div className="relative w-full max-w-sm overflow-hidden rounded-xl bg-card p-6 text-center text-card-foreground shadow-lg ring-1 ring-foreground/10">
        {!reducedMotion && (
          <div
            className="pointer-events-none absolute inset-0"
            aria-hidden="true"
          >
            {confettiPieces.map((piece) => (
              <span
                key={piece}
                className="onboarding-confetti absolute top-0 size-1.5 rounded-xs"
                style={{
                  left: `${(piece * 7 + 4) % 100}%`,
                  backgroundColor: `var(--color-chart-${(piece % 5) + 1})`,
                  animationDelay: `${(piece % 7) * 120}ms`,
                  animationDuration: `${1100 + (piece % 4) * 260}ms`,
                  ['--confetti-drift' as string]: `${((piece % 5) - 2) * 18}px`,
                  ['--confetti-spin' as string]: `${180 + (piece % 3) * 140}deg`,
                }}
              />
            ))}
          </div>
        )}

        <div className="relative flex justify-center">
          <PetAvatar
            stage="baby"
            emotion="hatching"
            satiety={90}
            happiness={95}
            energy={90}
            size="lg"
          />
        </div>

        <h2 id="hatch-title" className="mt-4 text-xl">
          {t('hatch.title')}
        </h2>
        <p className="mt-2 text-sm text-muted-foreground">
          {t('hatch.description')}
        </p>

        <Button ref={closeRef} onClick={onClose} className="mt-5 w-full">
          {t('hatch.cta')}
        </Button>
      </div>
    </div>
  )
}
