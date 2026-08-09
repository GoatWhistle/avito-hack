export const petHintStorageKey = 'avito-hack.items-pet-hint'

export const readPetHintDismissed = (): boolean => {
  if (typeof window === 'undefined') return false

  try {
    return window.localStorage.getItem(petHintStorageKey) === 'dismissed'
  } catch {
    return false
  }
}

export const writePetHintDismissed = () => {
  if (typeof window === 'undefined') return

  try {
    window.localStorage.setItem(petHintStorageKey, 'dismissed')
  } catch {
    return
  }
}
