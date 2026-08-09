import { useCallback, useEffect, useRef } from 'react'

export const useAuthFieldFocus = (order: readonly string[]) => {
  const refs = useRef<Record<string, HTMLInputElement | null>>({})

  const register = useCallback(
    (name: string) => (element: HTMLInputElement | null) => {
      refs.current[name] = element
    },
    [],
  )

  const focusFirstInvalid = useCallback(
    (invalid: readonly string[]) => {
      const target = order.find((name) => invalid.includes(name))
      if (!target) return

      refs.current[target]?.focus()
    },
    [order],
  )

  return { register, focusFirstInvalid }
}

export const useFocusOnServerField = (
  field: string | null,
  focus: (invalid: readonly string[]) => void,
) => {
  useEffect(() => {
    if (!field) return

    focus([field])
  }, [field, focus])
}
