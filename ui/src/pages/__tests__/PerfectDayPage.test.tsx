import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import { BrowserRouter } from 'react-router-dom';
import { PerfectDayPage } from '../PerfectDayPage';

// Mock the Timeline component since we're testing the page integration
vi.mock('../../components/timeline/Timeline', () => ({
  Timeline: ({ onAddActivity, onEditActivity }: any) => (
    <div data-testid="timeline">
      <button onClick={() => onAddActivity({ timeSlot: '10:00' })}>
        Add Activity
      </button>
      <button onClick={() => onEditActivity('act1')}>Edit Activity</button>
    </div>
  ),
}));

// Mock the ActivityModal component
vi.mock('../../components/activity/ActivityModal', () => ({
  ActivityModal: ({ isOpen, onClose, onSave }: any) =>
    isOpen ? (
      <div data-testid="activity-modal">
        <button onClick={() => onSave({ name: 'Test Activity' })}>Save</button>
        <button onClick={onClose}>Close</button>
      </div>
    ) : null,
}));

const renderWithRouter = (component: React.ReactElement) => {
  return render(<BrowserRouter>{component}</BrowserRouter>);
};

describe('PerfectDayPage', () => {
  const mockPerfectDay = {
    id: 'pd1',
    title: 'Amazing Tokyo Day',
    description: 'Perfect exploration of Tokyo',
    date: '2025-01-15',
    username: 'testuser',
    activities: [
      {
        id: 'act1',
        name: 'Morning Coffee',
        location: { type: 'custom_text', name: 'Local Cafe', area: 'Shibuya' },
        startTime: '09:00',
        duration: 60,
        description: 'Start the day right',
        commentary: 'Best coffee in town',
      },
    ],
  };

  it('should render perfect day title and description', () => {
    renderWithRouter(<PerfectDayPage perfectDay={mockPerfectDay} />);

    expect(screen.getByText('Amazing Tokyo Day')).toBeInTheDocument();
    expect(screen.getByText('Perfect exploration of Tokyo')).toBeInTheDocument();
  });

  it('should render perfect day date', () => {
    renderWithRouter(<PerfectDayPage perfectDay={mockPerfectDay} />);

    expect(screen.getByText('January 15, 2025')).toBeInTheDocument();
  });

  it('should render timeline component', () => {
    renderWithRouter(<PerfectDayPage perfectDay={mockPerfectDay} />);

    expect(screen.getByTestId('timeline')).toBeInTheDocument();
  });

  it('should show share button', () => {
    renderWithRouter(<PerfectDayPage perfectDay={mockPerfectDay} />);

    expect(screen.getByRole('button', { name: /share/i })).toBeInTheDocument();
  });

  it('should show edit mode toggle when user owns the perfect day', () => {
    renderWithRouter(
      <PerfectDayPage perfectDay={mockPerfectDay} isOwner={true} />
    );

    expect(screen.getByRole('button', { name: 'Edit' })).toBeInTheDocument();
  });

  it('should not show edit mode when user does not own the perfect day', () => {
    renderWithRouter(
      <PerfectDayPage perfectDay={mockPerfectDay} isOwner={false} />
    );

    expect(screen.queryByRole('button', { name: 'Edit' })).not.toBeInTheDocument();
  });

  it('should open activity modal when add activity is triggered', async () => {
    const user = userEvent.setup();
    renderWithRouter(
      <PerfectDayPage perfectDay={mockPerfectDay} isOwner={true} />
    );

    // First enable edit mode
    const editButton = screen.getByRole('button', { name: 'Edit' });
    await user.click(editButton);

    const addButton = screen.getByText('Add Activity');
    await user.click(addButton);

    expect(screen.getByTestId('activity-modal')).toBeInTheDocument();
  });

  it('should save new activity and update timeline', async () => {
    const user = userEvent.setup();
    const mockOnSave = vi.fn();
    renderWithRouter(
      <PerfectDayPage
        perfectDay={mockPerfectDay}
        isOwner={true}
        onSave={mockOnSave}
      />
    );

    // First enable edit mode
    const editButton = screen.getByRole('button', { name: 'Edit' });
    await user.click(editButton);

    // Trigger add activity
    const addButton = screen.getByText('Add Activity');
    await user.click(addButton);

    // Save activity
    const saveButton = screen.getByText('Save');
    await user.click(saveButton);

    expect(mockOnSave).toHaveBeenCalledWith(
      expect.objectContaining({
        activities: expect.arrayContaining([
          expect.objectContaining({ name: 'Test Activity' }),
        ]),
      })
    );
  });

  it('should show activity count', () => {
    renderWithRouter(<PerfectDayPage perfectDay={mockPerfectDay} />);

    expect(screen.getByText('1 activity')).toBeInTheDocument();
  });

  it('should show estimated total duration', () => {
    renderWithRouter(<PerfectDayPage perfectDay={mockPerfectDay} />);

    expect(screen.getByText('1h total')).toBeInTheDocument();
  });

  it('should handle share action', async () => {
    const user = userEvent.setup();
    const mockOnShare = vi.fn();
    renderWithRouter(
      <PerfectDayPage perfectDay={mockPerfectDay} onShare={mockOnShare} />
    );

    const shareButton = screen.getByRole('button', { name: /share/i });
    await user.click(shareButton);

    expect(mockOnShare).toHaveBeenCalledWith(mockPerfectDay.id);
  });

  it('should be responsive for mobile layout', () => {
    renderWithRouter(<PerfectDayPage perfectDay={mockPerfectDay} />);

    const pageContainer = screen.getByTestId('perfect-day-page');
    expect(pageContainer).toHaveClass(/mobile|responsive/);
  });

  it('should show loading state while saving', () => {
    renderWithRouter(
      <PerfectDayPage perfectDay={mockPerfectDay} isSaving={true} />
    );

    expect(screen.getByText(/saving/i)).toBeInTheDocument();
  });

  it('should display area tags from activities', () => {
    renderWithRouter(<PerfectDayPage perfectDay={mockPerfectDay} />);

    expect(screen.getByText('Shibuya')).toBeInTheDocument();
  });

  it('should handle edit activity flow', async () => {
    const user = userEvent.setup();
    renderWithRouter(
      <PerfectDayPage perfectDay={mockPerfectDay} isOwner={true} />
    );

    // First enable edit mode
    const editModeButton = screen.getByRole('button', { name: 'Edit' });
    await user.click(editModeButton);

    const editActivityButton = screen.getByText('Edit Activity');
    await user.click(editActivityButton);

    expect(screen.getByTestId('activity-modal')).toBeInTheDocument();
  });
});