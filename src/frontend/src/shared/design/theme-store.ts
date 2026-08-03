import { createEvent, createStore } from 'effector';

import type { ThemeMode } from './tokens/colors';

const STORAGE_KEY = 'theme';

function readStoredMode(): ThemeMode {
  try {
    const stored = window.localStorage.getItem(STORAGE_KEY);

    return stored === 'dark' ? 'dark' : 'light';
  } catch {
    return 'light';
  }
}

export const themeToggled = createEvent();
export const themeSet = createEvent<ThemeMode>();

export const $themeMode = createStore<ThemeMode>(readStoredMode())
  .on(themeToggled, (mode) => (mode === 'dark' ? 'light' : 'dark'))
  .on(themeSet, (_, mode) => mode);

$themeMode.watch((mode) => {
  try {
    window.localStorage.setItem(STORAGE_KEY, mode);
  } catch {
    return;
  }
});
