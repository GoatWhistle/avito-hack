import { useEffect, useState } from 'react';

const IDLE_TIMEOUT_MS = 45_000;
const ACTIVITY_EVENTS = ['pointerdown', 'pointermove', 'keydown', 'wheel'] as const;

export function useIdleSleep(): boolean {
  const [isAsleep, setIsAsleep] = useState(false);

  useEffect(() => {
    let timer: ReturnType<typeof setTimeout> | null = null;

    const schedule = (): void => {
      if (timer !== null) {
        clearTimeout(timer);
      }
      timer = setTimeout(() => {
        setIsAsleep(true);
      }, IDLE_TIMEOUT_MS);
    };

    const wake = (): void => {
      setIsAsleep(false);
      schedule();
    };

    const onVisibility = (): void => {
      if (document.visibilityState === 'visible') {
        wake();
      }
    };

    for (const event of ACTIVITY_EVENTS) {
      window.addEventListener(event, wake, { passive: true });
    }
    document.addEventListener('visibilitychange', onVisibility);
    schedule();

    return () => {
      if (timer !== null) {
        clearTimeout(timer);
      }
      for (const event of ACTIVITY_EVENTS) {
        window.removeEventListener(event, wake);
      }
      document.removeEventListener('visibilitychange', onVisibility);
    };
  }, []);

  return isAsleep;
}
