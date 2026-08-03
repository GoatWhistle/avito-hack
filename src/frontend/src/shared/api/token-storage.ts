const STORAGE_KEY = 'access_token';

let inMemoryToken: string | null = null;

export const tokenStorage = {
  get(): string | null {
    if (inMemoryToken !== null) {
      return inMemoryToken;
    }

    try {
      inMemoryToken = window.localStorage.getItem(STORAGE_KEY);
    } catch {
      inMemoryToken = null;
    }

    return inMemoryToken;
  },

  set(token: string): void {
    inMemoryToken = token;

    try {
      window.localStorage.setItem(STORAGE_KEY, token);
    } catch {
      inMemoryToken = token;
    }
  },

  clear(): void {
    inMemoryToken = null;

    try {
      window.localStorage.removeItem(STORAGE_KEY);
    } catch {
      inMemoryToken = null;
    }
  },
};
