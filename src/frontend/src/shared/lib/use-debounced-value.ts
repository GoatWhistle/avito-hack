import { useEffect, useState } from 'react';

const DEFAULT_DELAY_MS = 300;

export function useDebouncedValue<T>(value: T, delay: number = DEFAULT_DELAY_MS): T {
  const [debounced, setDebounced] = useState(value);

  useEffect(() => {
    const timer = setTimeout(() => {
      setDebounced(value);
    }, delay);

    return () => {
      clearTimeout(timer);
    };
  }, [value, delay]);

  return debounced;
}
