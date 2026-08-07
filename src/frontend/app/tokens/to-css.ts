import { palette, semanticDark, semanticLight } from './colors'
import { durations, easings, zLayers } from './motion'
import { radii, radiusBase } from './radii'
import { shadows } from './shadows'
import { breakpoints, containerWidths, spacingBase } from './spacing'
import {
  fontFamilies,
  fontSizes,
  fontWeights,
  letterSpacings,
} from './typography'

const declarations = (entries: [string, string][], indent = '  ') =>
  entries.map(([name, value]) => `${indent}${name}: ${value};`).join('\n')

export const themeBlock = () => {
  const entries: [string, string][] = [
    ...Object.entries(palette).map(
      ([token, value]) => [`--color-${token}`, value] as [string, string],
    ),
    ...Object.entries(fontFamilies).map(
      ([token, value]) => [`--font-${token}`, value] as [string, string],
    ),
    ...Object.entries(fontSizes).flatMap(
      ([token, { size, lineHeight }]) =>
        [
          [`--text-${token}`, size],
          [`--text-${token}--line-height`, lineHeight],
        ] as [string, string][],
    ),
    ...Object.entries(fontWeights).map(
      ([token, value]) => [`--font-weight-${token}`, value] as [string, string],
    ),
    ...Object.entries(letterSpacings).map(
      ([token, value]) => [`--tracking-${token}`, value] as [string, string],
    ),
    ['--spacing', spacingBase],
    ...Object.entries(containerWidths).map(
      ([token, value]) => [`--container-${token}`, value] as [string, string],
    ),
    ...Object.entries(radii).map(
      ([token, value]) => [`--radius-${token}`, value] as [string, string],
    ),
    ...Object.entries(shadows).map(
      ([token, value]) => [`--shadow-${token}`, value] as [string, string],
    ),
    ...Object.entries(zLayers).map(
      ([token, value]) => [`--z-${token}`, value] as [string, string],
    ),
    ...Object.entries(durations).map(
      ([token, value]) => [`--duration-${token}`, value] as [string, string],
    ),
    ...Object.entries(easings).map(
      ([token, value]) => [`--ease-${token}`, value] as [string, string],
    ),
    ...Object.entries(breakpoints).map(
      ([token, value]) => [`--breakpoint-${token}`, value] as [string, string],
    ),
  ]

  return `@theme {\n${declarations(entries)}\n}`
}

export const rootBlock = () => {
  const entries: [string, string][] = [
    ['--radius', radiusBase],
    ...Object.entries(semanticLight).map(
      ([token, value]) => [`--${token}`, value] as [string, string],
    ),
  ]

  return `:root {\n${declarations(entries)}\n}`
}

export const darkBlock = () => {
  const entries = Object.entries(semanticDark).map(
    ([token, value]) => [`--${token}`, value] as [string, string],
  )

  return `.dark {\n${declarations(entries)}\n}`
}

export const themeInlineBlock = () => {
  const entries = Object.keys(semanticLight).map(
    (token) => [`--color-${token}`, `var(--${token})`] as [string, string],
  )

  return `@theme inline {\n${declarations(entries)}\n}`
}

export const buildTokensCss = () =>
  [themeBlock(), rootBlock(), darkBlock(), themeInlineBlock()].join('\n\n') +
  '\n'
