import { semanticColors, type ThemeMode } from './tokens/colors';
import { radii } from './tokens/radii';
import { shadows } from './tokens/shadows';
import { layout, spacing } from './tokens/spacing';
import { fontFamily, fontSize, fontWeight, lineHeight } from './tokens/typography';

function kebab(value: string): string {
  return value.replace(/([a-z0-9])([A-Z])/g, '$1-$2').toLowerCase();
}

function toEntries(
  prefix: string,
  source: Record<string, string | number>,
  unit = '',
): [string, string][] {
  return Object.entries(source).map(([key, value]): [string, string] => [
    `--${prefix}-${kebab(key)}`,
    typeof value === 'number' && unit !== '' ? `${value}${unit}` : String(value),
  ]);
}

export function buildCssVars(mode: ThemeMode): Record<string, string> {
  const colors = semanticColors(mode);

  return Object.fromEntries<string>([
    ...Object.entries(colors).map(([key, value]): [string, string] => [`--${kebab(key)}`, value]),
    ...toEntries('spacing', spacing, 'px'),
    ...toEntries('radius', radii, 'px'),
    ...toEntries('shadow', shadows),
    ...toEntries('font-size', fontSize, 'px'),
    ...toEntries('font-weight', fontWeight),
    ...toEntries('line-height', lineHeight),
    ...toEntries('layout', layout, 'px'),
    ['--font-family-base', fontFamily.base],
    ['--font-family-mono', fontFamily.mono],
  ]);
}

export function applyCssVars(mode: ThemeMode, target: HTMLElement): void {
  const vars = buildCssVars(mode);

  for (const [name, value] of Object.entries(vars)) {
    target.style.setProperty(name, value);
  }

  target.dataset.theme = mode;
}
