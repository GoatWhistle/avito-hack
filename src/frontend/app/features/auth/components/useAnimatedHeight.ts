import {
  useEffect,
  useLayoutEffect,
  useRef,
  useState,
  type CSSProperties,
  type RefObject,
} from 'react'
import { useReducedMotion } from '#/lib/useReducedMotion'

interface AnimatedHeight {
  ref: RefObject<HTMLDivElement | null>
  style: CSSProperties | undefined
}

const measure = (node: HTMLDivElement): number => {
  const panel = node.querySelector<HTMLElement>('[data-slot="tabs-content"]')

  return panel ? panel.offsetHeight : node.offsetHeight
}

export function useAnimatedHeight(key: string): AnimatedHeight {
  const ref = useRef<HTMLDivElement>(null)
  const reducedMotion = useReducedMotion()
  const [height, setHeight] = useState<number | null>(null)
  const previousKey = useRef(key)

  useLayoutEffect(() => {
    const node = ref.current

    if (!node) return

    const next = measure(node)

    if (previousKey.current === key) {
      setHeight((current) => (current === null ? null : next))

      return
    }

    previousKey.current = key
    setHeight(next)
  }, [key])

  useEffect(() => {
    const node = ref.current

    if (!node || typeof ResizeObserver !== 'function') return

    const observer = new ResizeObserver(() => {
      setHeight((current) => (current === null ? null : measure(node)))
    })

    const panel = node.querySelector<HTMLElement>('[data-slot="tabs-content"]')

    if (panel) observer.observe(panel)

    return () => {
      observer.disconnect()
    }
  }, [key])

  if (reducedMotion || height === null) {
    return { ref, style: undefined }
  }

  return { ref, style: { height } }
}
