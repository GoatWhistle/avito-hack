import { existsSync, readdirSync, readFileSync, statSync } from 'node:fs';
import { join } from 'node:path';

const ROOT = 'src/frontend/app/i18n/locales';
const BASE = 'ru';
const PLURAL_SUFFIXES = ['_zero', '_one', '_two', '_few', '_many', '_other'];

if (!existsSync(ROOT)) {
  console.error(`Locales directory not found: ${ROOT}`);
  console.error('Update ROOT in scripts/check-locales.mjs if the frontend layout changed.');
  process.exit(1);
}

if (!existsSync(join(ROOT, BASE))) {
  console.error(`Base locale not found: ${join(ROOT, BASE)}`);
  process.exit(1);
}

const isDir = (path) => statSync(path).isDirectory();
const namespacesOf = (locale) =>
  readdirSync(join(ROOT, locale)).filter((entry) => entry.endsWith('.json'));

const locales = readdirSync(ROOT).filter(
  (entry) => entry !== BASE && isDir(join(ROOT, entry)),
);

const baseNamespaces = namespacesOf(BASE);

if (baseNamespaces.length === 0) {
  console.warn(`Base locale ${BASE} has no namespace files yet, nothing to compare.`);
  process.exit(0);
}

if (locales.length === 0) {
  console.log(`Only the base locale ${BASE} is present, nothing to compare.`);
  process.exit(0);
}

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

const readKeys = (locale, namespace) => {
  const path = join(ROOT, locale, namespace);

  try {
    return new Set(flatten(JSON.parse(readFileSync(path, 'utf8'))).map(stripPlural));
  } catch (error) {
    console.error(`[${locale}/${namespace}] cannot be parsed: ${error.message}`);

    return null;
  }
};

let failed = false;

for (const locale of locales) {
  const extraNamespaces = namespacesOf(locale).filter(
    (namespace) => !baseNamespaces.includes(namespace),
  );

  if (extraNamespaces.length > 0) {
    console.error(`[${locale}] namespaces missing in ${BASE}: ${extraNamespaces.join(', ')}`);
    failed = true;
  }
}

for (const namespace of baseNamespaces) {
  const baseKeys = readKeys(BASE, namespace);

  if (baseKeys === null) {
    failed = true;
    continue;
  }

  for (const locale of locales) {
    if (!existsSync(join(ROOT, locale, namespace))) {
      console.error(`[${locale}] missing namespace: ${namespace}`);
      failed = true;
      continue;
    }

    const keys = readKeys(locale, namespace);

    if (keys === null) {
      failed = true;
      continue;
    }

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
