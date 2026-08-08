export class LocalTokenStorage {
  private static readonly KEY = 'token'

  get() {
    return localStorage.getItem(LocalTokenStorage.KEY)
  }

  set(token: string) {
    localStorage.setItem(LocalTokenStorage.KEY, token)
  }

  remove() {
    localStorage.removeItem(LocalTokenStorage.KEY)
  }
}

export const localTokenStorage = new LocalTokenStorage()
