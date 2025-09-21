import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import { ActivityCard } from '../ActivityCard';

describe('ActivityCard Component', () => {
  const mockActivity = {
    id: 'act1',
    name: 'Morning Coffee',
    location: { type: 'custom_text', name: 'Local Cafe', area: 'Shibuya' },
    startTime: '09:00',
    duration: 60,
    description: 'Start the day right',
    commentary: 'Best coffee in town',
  };

  const mockProps = {
    activity: mockActivity,
    onEdit: vi.fn(),
    onDelete: vi.fn(),
    onClick: vi.fn(),
  };

  it('should display activity name and location', () => {
    render(<ActivityCard {...mockProps} />);

    expect(screen.getByText('Morning Coffee')).toBeInTheDocument();
    expect(screen.getByText('Local Cafe')).toBeInTheDocument();
    expect(screen.getByText('Shibuya')).toBeInTheDocument();
  });

  it('should display time range with start time and duration', () => {
    render(<ActivityCard {...mockProps} />);

    expect(screen.getByText('9:00 AM')).toBeInTheDocument();
    expect(screen.getByText('1h')).toBeInTheDocument();
  });

  it('should show location icon', () => {
    render(<ActivityCard {...mockProps} />);

    expect(screen.getByTestId('location-icon')).toBeInTheDocument();
  });

  it('should call onClick when card is clicked', async () => {
    const user = userEvent.setup();
    render(<ActivityCard {...mockProps} />);

    const card = screen.getByRole('button', { name: /morning coffee activity/i });
    await user.click(card);

    expect(mockProps.onClick).toHaveBeenCalledWith('act1');
  });

  it('should call onEdit when edit button is clicked', async () => {
    const user = userEvent.setup();
    render(<ActivityCard {...mockProps} />);

    const editButton = screen.getByRole('button', { name: /edit activity/i });
    await user.click(editButton);

    expect(mockProps.onEdit).toHaveBeenCalledWith('act1');
  });

  it('should call onDelete when delete button is clicked', async () => {
    const user = userEvent.setup();
    render(<ActivityCard {...mockProps} />);

    const deleteButton = screen.getByRole('button', { name: /delete activity/i });
    await user.click(deleteButton);

    expect(mockProps.onDelete).toHaveBeenCalledWith('act1');
  });

  it('should display activity description when provided', () => {
    render(<ActivityCard {...mockProps} />);

    expect(screen.getByText('Start the day right')).toBeInTheDocument();
  });

  it('should format duration correctly for different values', () => {
    const { rerender } = render(<ActivityCard {...mockProps} />);

    // 60 minutes = 1h
    expect(screen.getByText('1h')).toBeInTheDocument();

    // 90 minutes = 1h 30m
    rerender(
      <ActivityCard
        {...mockProps}
        activity={{ ...mockActivity, duration: 90 }}
      />
    );
    expect(screen.getByText('1h 30m')).toBeInTheDocument();

    // 30 minutes = 30m
    rerender(
      <ActivityCard
        {...mockProps}
        activity={{ ...mockActivity, duration: 30 }}
      />
    );
    expect(screen.getByText('30m')).toBeInTheDocument();
  });

  it('should be touch-friendly with proper button sizes', () => {
    render(<ActivityCard {...mockProps} />);

    const editButton = screen.getByRole('button', { name: /edit activity/i });
    const deleteButton = screen.getByRole('button', { name: /delete activity/i });

    // Touch targets should have proper size classes
    [editButton, deleteButton].forEach(button => {
      expect(button.className).toContain('min-h-[44px]');
      expect(button.className).toContain('min-w-[44px]');
    });
  });

  it('should show time icon', () => {
    render(<ActivityCard {...mockProps} />);

    expect(screen.getByTestId('time-icon')).toBeInTheDocument();
  });

  it('should handle activities without description gracefully', () => {
    const activityWithoutDescription = {
      ...mockActivity,
      description: '',
    };

    render(
      <ActivityCard
        {...mockProps}
        activity={activityWithoutDescription}
      />
    );

    expect(screen.getByText('Morning Coffee')).toBeInTheDocument();
    expect(screen.queryByText('Start the day right')).not.toBeInTheDocument();
  });
});