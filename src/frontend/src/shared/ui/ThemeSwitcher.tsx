import { MoonOutlined, SunOutlined } from '@ant-design/icons';
import { Button, Tooltip } from 'antd';
import { useUnit } from 'effector-react';
import { useTranslation } from 'react-i18next';

import { $themeMode, themeToggled } from '@/shared/design';

import './shared-ui.css';

export function ThemeSwitcher() {
  const { t } = useTranslation('common');
  const mode = useUnit($themeMode);

  const label = mode === 'dark' ? t('theme.light') : t('theme.dark');

  return (
    <Tooltip title={label}>
      <Button
        className="theme-switcher"
        type="text"
        aria-label={label}
        icon={mode === 'dark' ? <SunOutlined /> : <MoonOutlined />}
        onClick={() => {
          themeToggled();
        }}
      />
    </Tooltip>
  );
}
