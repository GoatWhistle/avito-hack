import js from '@eslint/js'
import pluginQuery from '@tanstack/eslint-plugin-query'
import reactHooks from 'eslint-plugin-react-hooks'
import { defineConfig, globalIgnores } from 'eslint/config'
import react from 'typescript-eslint'
import tseslint from 'typescript-eslint'

export default defineConfig([
  globalIgnores([
    '.idea/**',
    '.react-router/**',
    'build/**',
    'node_modules/**',
    '**/env*',
    '.gitignore',
    '.prettierignore',
    '.prettierrc',
    'bun.lock',
    'components.json',
    'public/**',
    '**/*.config.js',
    '**/*.config.ts',
  ]),
  js.configs.recommended,

  tseslint.configs.recommended,
  tseslint.configs.stylistic,

  react.configs.recommended,
  reactHooks.configs.flat.recommended,

  pluginQuery.configs['flat/recommended'],
])
