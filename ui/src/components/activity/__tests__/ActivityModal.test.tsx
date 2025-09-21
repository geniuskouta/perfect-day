import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import { ActivityModal } from '../ActivityModal';

describe('ActivityModal Component', () => {
  const mockProps = {
    isOpen: true,
    onClose: vi.fn(),
    onSave: vi.fn(),
    timeSlot: '10:00',
    initialActivity: null,
  };

  const existingActivity = {
    id: 'act1',
    name: 'Morning Coffee',
    location: { type: 'custom_text', name: 'Local Cafe', area: 'Shibuya' },
    startTime: '09:00',
    duration: 60,
    description: 'Start the day right',
    commentary: 'Best coffee in town',
  };

  it('should render modal when isOpen is true', () => {
    render(<ActivityModal {...mockProps} />);

    expect(screen.getByRole('dialog')).toBeInTheDocument();
    expect(screen.getByText('Add Activity')).toBeInTheDocument();
  });

  it('should not render modal when isOpen is false', () => {
    render(<ActivityModal {...mockProps} isOpen={false} />);

    expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
  });

  it('should show edit mode when initialActivity is provided', () => {
    render(
      <ActivityModal {...mockProps} initialActivity={existingActivity} />
    );

    expect(screen.getByText('Edit Activity')).toBeInTheDocument();
    expect(screen.getByDisplayValue('Morning Coffee')).toBeInTheDocument();
  });

  it('should have form fields for activity details', () => {
    render(<ActivityModal {...mockProps} />);

    expect(screen.getByLabelText(/activity name/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/location/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/start time/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/duration/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/description/i)).toBeInTheDocument();
  });

  it('should pre-fill start time from timeSlot prop', () => {
    render(<ActivityModal {...mockProps} timeSlot="14:30" />);

    const timeInput = screen.getByLabelText(/start time/i);
    expect(timeInput).toHaveValue('14:30');
  });

  it('should call onClose when close button is clicked', async () => {
    const user = userEvent.setup();
    render(<ActivityModal {...mockProps} />);

    const closeButton = screen.getByRole('button', { name: /close/i });
    await user.click(closeButton);

    expect(mockProps.onClose).toHaveBeenCalled();
  });

  it('should call onClose when clicking outside modal', async () => {
    const user = userEvent.setup();
    render(<ActivityModal {...mockProps} />);

    const overlay = screen.getByTestId('modal-overlay');
    await user.click(overlay);

    expect(mockProps.onClose).toHaveBeenCalled();
  });

  it('should validate required fields before saving', async () => {
    const user = userEvent.setup();
    render(<ActivityModal {...mockProps} />);

    const saveButton = screen.getByRole('button', { name: /save/i });
    await user.click(saveButton);

    // Should show validation errors
    expect(screen.getByText(/activity name is required/i)).toBeInTheDocument();
    expect(mockProps.onSave).not.toHaveBeenCalled();
  });

  it('should call onSave with form data when valid', async () => {
    const user = userEvent.setup();
    render(<ActivityModal {...mockProps} />);

    // Fill in form
    await user.type(screen.getByLabelText(/activity name/i), 'Test Activity');
    await user.type(screen.getByLabelText(/location/i), 'Test Location');
    await user.selectOptions(screen.getByLabelText(/duration/i), '60');
    await user.type(screen.getByLabelText(/description/i), 'Test description');

    const saveButton = screen.getByRole('button', { name: /save/i });
    await user.click(saveButton);

    expect(mockProps.onSave).toHaveBeenCalledWith({
      name: 'Test Activity',
      location: expect.objectContaining({
        name: 'Test Location',
        type: 'google_places',
        area: 'Unknown',
      }),
      startTime: '10:00',
      duration: 60,
      description: 'Test description',
      commentary: '',
    });
  });

  it('should support location search with suggestions', async () => {
    const user = userEvent.setup();
    render(<ActivityModal {...mockProps} />);

    const locationInput = screen.getByLabelText(/location/i);
    await user.type(locationInput, 'Starbucks');

    // Should show search suggestions
    expect(screen.getByText(/searching locations/i)).toBeInTheDocument();
  });

  it('should have mobile-optimized form layout', () => {
    render(<ActivityModal {...mockProps} />);

    const modal = screen.getByRole('dialog');
    expect(modal).toHaveClass(/bottom-sheet|slide-up/); // Mobile bottom sheet style
  });

  it('should support duration picker with common durations', () => {
    render(<ActivityModal {...mockProps} />);

    const durationSelect = screen.getByLabelText(/duration/i);
    expect(screen.getByRole('option', { name: /30 minutes/i })).toBeInTheDocument();
    expect(screen.getByRole('option', { name: /1 hour/i })).toBeInTheDocument();
    expect(screen.getByRole('option', { name: /2 hours/i })).toBeInTheDocument();
  });

  it('should show camera icon for photo attachment', () => {
    render(<ActivityModal {...mockProps} />);

    expect(screen.getByTestId('camera-icon')).toBeInTheDocument();
  });

  it('should allow custom location entry', async () => {
    const user = userEvent.setup();
    render(<ActivityModal {...mockProps} />);

    const locationInput = screen.getByLabelText(/location/i);
    await user.type(locationInput, 'Custom Location');

    const customLocationButton = screen.getByRole('button', {
      name: /use custom location/i,
    });
    await user.click(customLocationButton);

    expect(screen.getByText('Custom location mode enabled')).toBeInTheDocument();
  });
});