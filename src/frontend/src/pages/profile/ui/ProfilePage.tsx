import { App, Button, Card, Form, Input, Space, Typography } from 'antd';
import { useUnit } from 'effector-react';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';

import { $user, profileUpdated, sessionApi } from '@/entities/session';
import { useApiErrorMessage } from '@/shared/lib/use-api-error';
import { PageSkeleton } from '@/shared/ui';

const FORM_MAX_WIDTH = 480;

export function ProfilePage() {
  const { t } = useTranslation(['auth', 'common']);
  const { notification } = App.useApp();
  const user = useUnit($user);
  const getErrorMessage = useApiErrorMessage();

  const [fullName, setFullName] = useState(user?.fullName ?? '');
  const [isSaving, setIsSaving] = useState(false);

  if (user === null) {
    return <PageSkeleton />;
  }

  const save = async () => {
    setIsSaving(true);

    try {
      const updated = await sessionApi.updateProfile({ fullName });
      profileUpdated(updated);
      notification.success({ message: t('auth:profile.updated') });
    } catch (error) {
      notification.error({ message: getErrorMessage(error) });
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <Space direction="vertical" size="middle" style={{ width: '100%', maxWidth: FORM_MAX_WIDTH }}>
      <Typography.Title level={3}>{t('auth:profile.title')}</Typography.Title>

      <Card>
        <Form layout="vertical" onFinish={() => void save()}>
          <Form.Item label={t('auth:fields.email')}>
            <Input value={user.email} disabled />
          </Form.Item>

          <Form.Item label={t('auth:fields.fullName')}>
            <Input
              value={fullName}
              onChange={(event) => {
                setFullName(event.target.value);
              }}
            />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={isSaving}>
              {t('common:actions.save')}
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </Space>
  );
}
