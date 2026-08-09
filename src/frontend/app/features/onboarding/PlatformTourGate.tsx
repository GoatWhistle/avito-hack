import { useCallback, useEffect, useState } from 'react'
import { PlatformTour } from './PlatformTour'
import {
  completePlatformTour,
  shouldShowPlatformTour,
} from './onboarding-state'

export function PlatformTourGate() {
  const [open, setOpen] = useState(false)

  useEffect(() => {
    setOpen(shouldShowPlatformTour())
  }, [])

  const close = useCallback(() => {
    completePlatformTour()
    setOpen(false)
  }, [])

  if (!open) return null

  return <PlatformTour onFinish={close} onSkip={close} />
}
