import type { Locale as AntdLocale } from 'antd/es/locale';
import enUS from 'antd/locale/en_US';
import ruRU from 'antd/locale/ru_RU';

export function antdLocale(language: string): AntdLocale {
  return language.startsWith('en') ? enUS : ruRU;
}
