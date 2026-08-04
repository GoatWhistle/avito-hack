import { Form } from 'antd';
import type { PropsWithChildren } from 'react';

interface FormFieldProps {
  name: string;
  label: string;
  error?: string | undefined;
}

export function FormField({ name, label, error, children }: PropsWithChildren<FormFieldProps>) {
  return (
    <Form.Item
      label={label}
      htmlFor={name}
      {...(error === undefined ? {} : { validateStatus: 'error' as const, help: error })}
    >
      {children}
    </Form.Item>
  );
}
