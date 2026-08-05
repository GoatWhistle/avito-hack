import { Segmented } from 'antd';
import { useTranslation } from 'react-i18next';

import { isSupportedLocale, SUPPORTED_LOCALES } from '@/shared/i18n';

import './shared-ui.css';

export function LanguageSwitcher() {
  const { t, i18n } = useTranslation('common');

  const current = isSupportedLocale(i18n.language) ? i18n.language : SUPPORTED_LOCALES[0];

  return (
    <Segmented
      className="language-switcher"
      size="small"
      value={current}
      aria-label={t('language.label')}
      onChange={(value: string) => {
        void i18n.changeLanguage(value);
      }}
      options={SUPPORTED_LOCALES.map((locale) => ({
        value: locale,
        label: locale.toUpperCase(),
      }))}
    />
  );
}
