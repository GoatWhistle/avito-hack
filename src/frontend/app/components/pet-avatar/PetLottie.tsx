import {
  useCallback,
  useRef,
  useState,
  type CSSProperties,
  type KeyboardEvent,
} from 'react'
import { DotLottieReact, setWasmUrl } from '@lottiefiles/dotlottie-react'
import { useGaze } from './hooks/useGaze'
import { useReducedMotion } from './hooks/useReducedMotion'
import { useStroke } from './hooks/useStroke'
import { LOTTIE_WASM_URL } from './lottie-sources'
import { PetPlaceholder } from './PetPlaceholder'

import './pet-lottie.css'

setWasmUrl(LOTTIE_WASM_URL)

export interface PetLottieProps {
  src: string
  ariaLabel: string
  placeholderLabel: string
  sizeClass: string
  stage: string
  mood: string
  emotion: string
  legend?: boolean
  className?: string
  onStroke?: () => void
}

export function PetLottie({
  src,
  ariaLabel,
  placeholderLabel,
  sizeClass,
  stage,
  mood,
  emotion,
  legend = false,
  className,
  onStroke,
}: PetLottieProps) {
  const reducedMotion = useReducedMotion()
  const [isReady, setIsReady] = useState(false)
  const [failed, setFailed] = useState(false)
  const { isPetted, trigger } = useStroke(onStroke)
  const rootRef = useRef<HTMLDivElement>(null)
  const gaze = useGaze(rootRef, !reducedMotion)

  const handleKeyDown = useCallback(
    (event: KeyboardEvent<HTMLDivElement>) => {
      if (event.key === 'Enter' || event.key === ' ') {
        event.preventDefault()
        trigger()
      }
    },
    [trigger],
  )

  const handleLoadError = useCallback(() => {
    setFailed(true)
  }, [])

  if (failed) {
    return (
      <PetPlaceholder
        label={placeholderLabel}
        sizeClass={sizeClass.replace('pet-lottie--', 'pet-placeholder--')}
        stage={stage}
        failed
        className={className}
      />
    )
  }

  const classes = [
    'pet-lottie',
    sizeClass,
    legend ? 'pet-lottie--legend' : '',
    isPetted ? 'pet-lottie--petted' : '',
    reducedMotion ? 'pet-lottie--reduced' : '',
    isReady ? 'pet-lottie--ready' : '',
    className ?? '',
  ]
    .filter((value) => value !== '')
    .join(' ')

  const gazeStyle = {
    '--pet-gaze-x': `${gaze.x}px`,
    '--pet-gaze-y': `${gaze.y}px`,
    '--pet-gaze-tilt': `${gaze.tilt}deg`,
  } as CSSProperties

  return (
    <div
      ref={rootRef}
      className={classes}
      style={gazeStyle}
      role="button"
      tabIndex={0}
      aria-label={ariaLabel}
      data-testid="pet-lottie"
      data-stage={stage}
      data-mood={mood}
      data-emotion={emotion}
      onClick={trigger}
      onKeyDown={handleKeyDown}
    >
      {legend && <span className="pet-lottie__halo" aria-hidden="true" />}

      <div className="pet-lottie__stage">
        {legend && <span className="pet-lottie__crown" aria-hidden="true" />}

        <DotLottieReact
          className="pet-lottie__canvas"
          src={src}
          autoplay={!reducedMotion}
          loop={!reducedMotion}
          dotLottieRefCallback={(instance) => {
            if (instance === null) return
            instance.addEventListener('load', () => setIsReady(true))
            instance.addEventListener('loadError', handleLoadError)
          }}
        />
      </div>

      {!isReady && (
        <span className="pet-lottie__skeleton" aria-hidden="true">
          <span className="pet-placeholder__orb" />
        </span>
      )}
    </div>
  )
}
