import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import { Timeline } from '../Timeline';

describe('Timeline Component', () => {
  const mockActivities = [
    {
      id: 'act1',
      name: 'Morning Coffee',
      location: { type: 'custom_text', name: 'Local Cafe', area: 'Shibuya' },
      startTime: '09:00',
      duration: 60,
      description: 'Start the day right',
      commentary: 'Best coffee in town',
    },
    {
      id: 'act2',
      name: 'Lunch Meeting',
      location: { type: 'custom_text', name: 'Restaurant', area: 'Shinjuku' },
      startTime: '12:00',
      duration: 90,
      description: 'Business lunch',
      commentary: 'Great conversation',
    },
  ];

  const mockProps = {
    activities: mockActivities,
    onAddActivity: vi.fn(),
    onEditActivity: vi.fn(),
    onDeleteActivity: vi.fn(),
  };

  it('should render timeline with hourly slots from 9 AM to 11 PM', () => {
    render(<Timeline {...mockProps} />);

    // Should show time slots for each hour
    expect(screen.getByText('9:00 AM')).toBeInTheDocument();
    expect(screen.getByText('12:00 PM')).toBeInTheDocument();
    expect(screen.getByText('6:00 PM')).toBeInTheDocument();
    expect(screen.getByText('11:00 PM')).toBeInTheDocument();
  });

  it('should display activities in correct time slots', () => {
    render(<Timeline {...mockProps} />);

    // Should show activities at their scheduled times
    expect(screen.getByText('Morning Coffee')).toBeInTheDocument();
    expect(screen.getByText('Local Cafe')).toBeInTheDocument();
    expect(screen.getByText('Lunch Meeting')).toBeInTheDocument();
    expect(screen.getByText('Restaurant')).toBeInTheDocument();
  });

  it('should show add buttons on empty time slots', () => {
    render(<Timeline {...mockProps} />);

    // Should have add buttons for empty slots
    const addButtons = screen.getAllByRole('button', { name: /add activity/i });
    expect(addButtons.length).toBeGreaterThan(0);
  });

  it('should call onAddActivity when add button is clicked', async () => {
    const user = userEvent.setup();
    render(<Timeline {...mockProps} />);

    const addButton = screen.getAllByRole('button', { name: /add activity/i })[0];
    await user.click(addButton);

    expect(mockProps.onAddActivity).toHaveBeenCalledWith(
      expect.objectContaining({
        timeSlot: expect.any(String),
      })
    );
  });

  it('should show current time indicator', () => {
    render(<Timeline {...mockProps} />);

    // Should have a visual indicator for current time
    expect(screen.getByTestId('current-time-indicator')).toBeInTheDocument();
  });

  it('should be optimized for mobile touch interactions', () => {
    render(<Timeline {...mockProps} />);

    // Touch targets should have proper size classes
    const addButtons = screen.getAllByRole('button', { name: /add activity/i });
    addButtons.forEach(button => {
      expect(button.className).toContain('min-h-[44px]');
    });
  });

  it('should display activity duration visually', () => {
    render(<Timeline {...mockProps} />);

    // Should show duration for activities
    expect(screen.getByText('1h')).toBeInTheDocument(); // 60 min
    expect(screen.getByText('1h 30m')).toBeInTheDocument(); // 90 min
  });

  it('should call onEditActivity when activity is clicked', async () => {
    const user = userEvent.setup();
    render(<Timeline {...mockProps} />);

    const activityCard = screen.getByText('Morning Coffee').closest('button');
    expect(activityCard).toBeInTheDocument();

    await user.click(activityCard!);
    expect(mockProps.onEditActivity).toHaveBeenCalledWith('act1');
  });
});