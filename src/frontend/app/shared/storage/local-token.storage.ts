export class LocalTokenStorage {
  private readonly KEY = 'token'

  get() {
    return localStorage.getItem(this.KEY)
  }

  set(token: string) {
    localStorage.setItem(this.KEY, token)
  }

  remove() {
    localStorage.removeItem(this.KEY)
  }
}

export const localTokenStorage = new LocalTokenStorage()
