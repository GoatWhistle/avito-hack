import { existsSync } from 'node:fs'
import { fileURLToPath } from 'node:url'

export function resolve(specifier, context, next) {
  if (specifier.startsWith('.') && !/\.[cm]?[jt]sx?$/.test(specifier)) {
    const base = new URL(specifier, context.parentURL)
    for (const ext of ['.ts', '.tsx', '/index.ts', '/index.tsx']) {
      const candidate = new URL(base.href + ext)
      if (existsSync(fileURLToPath(candidate))) {
        return next(specifier + ext, context)
      }
    }
  }
  return next(specifier, context)
}
