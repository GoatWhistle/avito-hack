import { writeFileSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'

import { buildTokensCss } from '../app/tokens/to-css.ts'

const here = dirname(fileURLToPath(import.meta.url))
const target = join(here, '..', 'app', 'styles', 'tokens.generated.css')

writeFileSync(target, buildTokensCss(), 'utf8')
console.log(`Generated ${target}`)
