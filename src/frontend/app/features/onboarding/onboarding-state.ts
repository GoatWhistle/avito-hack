export const onboardingStorageKey = 'avito-hack.onboarding'

export interface OnboardingState {
  dismissed: boolean
  tourSeen: boolean
  tourRequested: boolean
}

export const defaultOnboardingState: OnboardingState = {
  dismissed: false,
  tourSeen: false,
  tourRequested: false,
}

export const readOnboardingState = (): OnboardingState => {
  if (typeof window === 'undefined') return defaultOnboardingState

  try {
    const raw = window.localStorage.getItem(onboardingStorageKey)
    if (!raw) return defaultOnboardingState

    const parsed: unknown = JSON.parse(raw)
    if (typeof parsed !== 'object' || parsed === null) {
      return defaultOnboardingState
    }

    return {
      ...defaultOnboardingState,
      ...(parsed as Partial<OnboardingState>),
    }
  } catch {
    return defaultOnboardingState
  }
}

export const writeOnboardingState = (state: OnboardingState) => {
  if (typeof window === 'undefined') return

  try {
    window.localStorage.setItem(onboardingStorageKey, JSON.stringify(state))
  } catch {
    return
  }
}

const patchOnboardingState = (patch: Partial<OnboardingState>) => {
  writeOnboardingState({ ...readOnboardingState(), ...patch })
}

export const requestPlatformTour = () => {
  const state = readOnboardingState()
  if (state.tourSeen) {
    if (state.tourRequested) {
      writeOnboardingState({ ...state, tourRequested: false })
    }

    return
  }

  writeOnboardingState({ ...state, tourRequested: true })
}

export const completePlatformTour = () => {
  patchOnboardingState({ tourSeen: true, tourRequested: false })
}

export const shouldShowPlatformTour = (): boolean => {
  const state = readOnboardingState()

  return state.tourRequested && !state.tourSeen
}
