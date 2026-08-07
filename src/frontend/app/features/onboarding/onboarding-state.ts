export const onboardingStorageKey = 'avito-hack.onboarding'

export interface OnboardingState {
  dismissed: boolean
  hatched: boolean
  celebrated: boolean
}

export const defaultOnboardingState: OnboardingState = {
  dismissed: false,
  hatched: false,
  celebrated: false,
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
