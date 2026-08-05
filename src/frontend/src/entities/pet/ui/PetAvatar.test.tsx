import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';

import { renderWithProviders } from '@/shared/lib/test-utils';

import { PetAvatar } from './PetAvatar';

describe('PetAvatar', () => {
  it('describes the current stage and emotion for screen readers', () => {
    renderWithProviders(<PetAvatar stage="teen" emotion="happy" satiety={80} happiness={90} />);

    const avatar = screen.getByRole('img');

    expect(avatar).toHaveAttribute('aria-label', expect.stringContaining('teen'));
    expect(avatar).toHaveAttribute('aria-label', expect.stringContaining('happy'));
  });

  it('carries the emotion as a modifier class so CSS can animate it', () => {
    renderWithProviders(<PetAvatar stage="adult" emotion="levelup" satiety={70} happiness={70} />);

    expect(screen.getByRole('img')).toHaveClass('pet-avatar--levelup');
  });

  it('calls back when the pet is clicked', async () => {
    const onStroke = vi.fn();
    renderWithProviders(
      <PetAvatar stage="baby" emotion="idle" satiety={50} happiness={50} onStroke={onStroke} />,
    );

    await userEvent.click(screen.getByRole('img'));

    expect(onStroke).toHaveBeenCalledTimes(1);
  });

  it('is reachable and strokeable from the keyboard', async () => {
    const onStroke = vi.fn();
    renderWithProviders(
      <PetAvatar stage="baby" emotion="idle" satiety={50} happiness={50} onStroke={onStroke} />,
    );

    const avatar = screen.getByRole('img');
    avatar.focus();
    await userEvent.keyboard('{Enter}');

    expect(onStroke).toHaveBeenCalledTimes(1);
  });

  it('renders the egg silhouette before hatching', () => {
    renderWithProviders(<PetAvatar stage="egg" emotion="hatching" satiety={0} happiness={0} />);

    expect(screen.getByRole('img')).toHaveClass('pet-avatar--hatching');
  });
});
