import { Select } from 'antd';
import { useTranslation } from 'react-i18next';

import { isSupportedLocale, SUPPORTED_LOCALES } from '@/shared/i18n';

const SELECT_WIDTH = 120;

export function LanguageSwitcher() {
  const { t, i18n } = useTranslation('common');

  const current = isSupportedLocale(i18n.language) ? i18n.language : SUPPORTED_LOCALES[0];

  return (
    <Select
      value={current}
      style={{ width: SELECT_WIDTH }}
      aria-label={t('language.ru')}
      onChange={(value) => {
        void i18n.changeLanguage(value);
      }}
      options={SUPPORTED_LOCALES.map((locale) => ({
        value: locale,
        label: t(`language.${locale}`),
      }))}
    />
  );
}
