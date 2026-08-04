import { existsSync, readdirSync, readFileSync } from 'node:fs';
import { join } from 'node:path';

const ROOT = 'src/frontend/src/shared/i18n/locales';
const BASE = 'ru';
const PLURAL_SUFFIXES = ['_zero', '_one', '_two', '_few', '_many', '_other'];

if (!existsSync(join(ROOT, BASE))) {
  console.error(`Base locale not found: ${join(ROOT, BASE)}`);
  process.exit(1);
}

const locales = readdirSync(ROOT).filter((entry) => entry !== BASE);

const flatten = (obj, prefix = '') =>
  Object.entries(obj).flatMap(([key, value]) =>
    typeof value === 'object' && value !== null
      ? flatten(value, `${prefix}${key}.`)
      : [`${prefix}${key}`],
  );

const stripPlural = (key) => {
  const suffix = PLURAL_SUFFIXES.find((candidate) => key.endsWith(candidate));

  return suffix ? key.slice(0, -suffix.length) : key;
};

const readKeys = (locale, namespace) =>
  new Set(flatten(JSON.parse(readFileSync(join(ROOT, locale, namespace), 'utf8'))).map(stripPlural));

let failed = false;

for (const namespace of readdirSync(join(ROOT, BASE))) {
  const baseKeys = readKeys(BASE, namespace);

  for (const locale of locales) {
    if (!existsSync(join(ROOT, locale, namespace))) {
      console.error(`[${locale}] missing namespace: ${namespace}`);
      failed = true;
      continue;
    }

    const keys = readKeys(locale, namespace);
    const missing = [...baseKeys].filter((key) => !keys.has(key));
    const extra = [...keys].filter((key) => !baseKeys.has(key));

    if (missing.length > 0) {
      console.error(`[${locale}/${namespace}] missing keys: ${missing.join(', ')}`);
      failed = true;
    }
    if (extra.length > 0) {
      console.error(`[${locale}/${namespace}] unexpected keys: ${extra.join(', ')}`);
      failed = true;
    }
  }
}

if (!failed) {
  console.log(`Locales are in sync: ${[BASE, ...locales].join(', ')}`);
}

process.exit(failed ? 1 : 0);
