import { useEffect, useRef, useState, type RefObject } from 'react';

export interface EyeOffset {
  x: number;
  y: number;
}

const MAX_OFFSET_PX = 3;
const REACH_PX = 320;

function clampOffset(value: number): number {
  return Math.max(-MAX_OFFSET_PX, Math.min(MAX_OFFSET_PX, value));
}

export function useEyeTracking<T extends Element>(
  enabled: boolean,
): { ref: RefObject<T>; offset: EyeOffset } {
  const ref = useRef<T>(null);
  const [offset, setOffset] = useState<EyeOffset>({ x: 0, y: 0 });

  useEffect(() => {
    if (!enabled) {
      setOffset({ x: 0, y: 0 });

      return;
    }

    let frame: number | null = null;

    const onMove = (event: PointerEvent): void => {
      if (frame !== null) {
        return;
      }

      frame = requestAnimationFrame(() => {
        frame = null;
        const node = ref.current;

        if (node === null) {
          return;
        }

        const rect = node.getBoundingClientRect();
        const centerX = rect.left + rect.width / 2;
        const centerY = rect.top + rect.height / 2;

        setOffset({
          x: clampOffset(((event.clientX - centerX) / REACH_PX) * MAX_OFFSET_PX),
          y: clampOffset(((event.clientY - centerY) / REACH_PX) * MAX_OFFSET_PX),
        });
      });
    };

    window.addEventListener('pointermove', onMove, { passive: true });

    return () => {
      if (frame !== null) {
        cancelAnimationFrame(frame);
      }
      window.removeEventListener('pointermove', onMove);
    };
  }, [enabled]);

  return { ref, offset };
}

export function useReducedMotion(): boolean {
  const [reduced, setReduced] = useState(false);

  useEffect(() => {
    const query = window.matchMedia('(prefers-reduced-motion: reduce)');
    setReduced(query.matches);

    const onChange = (event: MediaQueryListEvent): void => {
      setReduced(event.matches);
    };

    query.addEventListener('change', onChange);

    return () => {
      query.removeEventListener('change', onChange);
    };
  }, []);

  return reduced;
}
