import React, { useState, useEffect } from 'react';
import { CloseModalIcon, LocationIcon } from '../ui/Icons';
import { cn } from '../../lib/utils';

interface Activity {
  id: string;
  name: string;
  location: {
    type: string;
    name: string;
    area: string;
  };
  startTime: string;
  duration: number;
  description: string;
  commentary: string;
}

interface ActivityModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSave: (activityData: Partial<Activity>) => void;
  timeSlot: string;
  initialActivity: Activity | null;
}

interface FormData {
  name: string;
  location: string;
  area: string;
  startTime: string;
  duration: number;
  description: string;
  commentary: string;
}

interface FormErrors {
  name?: string;
  location?: string;
}

const durationOptions = [
  { value: 15, label: '15 minutes' },
  { value: 30, label: '30 minutes' },
  { value: 45, label: '45 minutes' },
  { value: 60, label: '1 hour' },
  { value: 90, label: '1.5 hours' },
  { value: 120, label: '2 hours' },
  { value: 180, label: '3 hours' },
  { value: 240, label: '4 hours' },
];

export const ActivityModal: React.FC<ActivityModalProps> = ({
  isOpen,
  onClose,
  onSave,
  timeSlot,
  initialActivity,
}) => {
  const [formData, setFormData] = useState<FormData>({
    name: '',
    location: '',
    area: '',
    startTime: timeSlot,
    duration: 60,
    description: '',
    commentary: '',
  });

  const [errors, setErrors] = useState<FormErrors>({});
  const [isSearchingLocation, setIsSearchingLocation] = useState(false);
  const [showCustomLocation, setShowCustomLocation] = useState(false);

  useEffect(() => {
    if (initialActivity) {
      setFormData({
        name: initialActivity.name,
        location: initialActivity.location.name,
        area: initialActivity.location.area,
        startTime: initialActivity.startTime,
        duration: initialActivity.duration,
        description: initialActivity.description,
        commentary: initialActivity.commentary,
      });
    } else {
      setFormData(prev => ({
        ...prev,
        startTime: timeSlot,
      }));
    }
  }, [initialActivity, timeSlot]);

  const validateForm = (): boolean => {
    const newErrors: FormErrors = {};

    if (!formData.name.trim()) {
      newErrors.name = 'Activity name is required';
    }

    if (!formData.location.trim()) {
      newErrors.location = 'Location is required';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) {
      return;
    }

    const activityData = {
      name: formData.name,
      location: {
        type: showCustomLocation ? 'custom_text' : 'google_places',
        name: formData.location,
        area: formData.area || 'Unknown',
      },
      startTime: formData.startTime,
      duration: formData.duration,
      description: formData.description,
      commentary: formData.commentary,
    };

    onSave(activityData);
  };

  const handleLocationSearch = (value: string) => {
    setFormData(prev => ({ ...prev, location: value }));
    if (value.length > 2) {
      setIsSearchingLocation(true);
      // Simulate search delay
      setTimeout(() => {
        setIsSearchingLocation(false);
      }, 500);
    }
  };

  const handleOverlayClick = (e: React.MouseEvent) => {
    if (e.target === e.currentTarget) {
      onClose();
    }
  };

  if (!isOpen) return null;

  return (
    <div
      className="fixed inset-0 bg-black/50 z-50 flex items-end md:items-center justify-center p-4"
      data-testid="modal-overlay"
      onClick={handleOverlayClick}
    >
      <div
        className={cn(
          'bg-background shadow-lg',
          'w-full md:w-full md:max-w-md',
          'max-h-[85vh] md:max-h-[90vh] overflow-y-auto',
          // Mobile: bottom sheet with rounded top corners
          'rounded-t-xl md:rounded-xl',
          // Animation classes
          'bottom-sheet slide-up'
        )}
        role="dialog"
        aria-labelledby="modal-title"
      >
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-border">
          <h2 id="modal-title" className="text-lg font-semibold text-foreground">
            {initialActivity ? 'Edit Activity' : 'Add Activity'}
          </h2>
          <button
            onClick={onClose}
            className="p-2 hover:bg-accent rounded-md transition-colors"
            aria-label="Close"
          >
            <CloseModalIcon size="sm" />
          </button>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="p-4 space-y-4">
          {/* Activity Name */}
          <div>
            <label
              htmlFor="activity-name"
              className="block text-sm font-medium text-foreground mb-1"
            >
              Activity Name
            </label>
            <input
              id="activity-name"
              type="text"
              value={formData.name}
              onChange={e =>
                setFormData(prev => ({ ...prev, name: e.target.value }))
              }
              className={cn(
                'w-full p-3 rounded-md border border-input',
                'bg-background text-foreground',
                'focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2',
                errors.name && 'border-destructive'
              )}
              placeholder="e.g., Morning Coffee"
            />
            {errors.name && (
              <p className="text-sm text-destructive mt-1">{errors.name}</p>
            )}
          </div>

          {/* Location */}
          <div>
            <label
              htmlFor="location"
              className="block text-sm font-medium text-foreground mb-1"
            >
              Location
            </label>
            <div className="relative">
              <div className="absolute left-3 top-1/2 transform -translate-y-1/2">
                <LocationIcon size="sm" className="text-muted-foreground" />
              </div>
              <input
                id="location"
                type="text"
                value={formData.location}
                onChange={e => handleLocationSearch(e.target.value)}
                className={cn(
                  'w-full p-3 pl-10 rounded-md border border-input',
                  'bg-background text-foreground',
                  'focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2',
                  errors.location && 'border-destructive'
                )}
                placeholder="Search for a place or enter custom location"
              />
            </div>
            {isSearchingLocation && (
              <p className="text-sm text-muted-foreground mt-1">
                Searching locations...
              </p>
            )}
            {errors.location && (
              <p className="text-sm text-destructive mt-1">{errors.location}</p>
            )}
            {formData.location && (
              <button
                type="button"
                onClick={() => setShowCustomLocation(!showCustomLocation)}
                className="text-sm text-primary hover:underline mt-1"
              >
                Use custom location
              </button>
            )}
            {showCustomLocation && (
              <p className="text-sm text-muted-foreground mt-1">
                Custom location mode enabled
              </p>
            )}
          </div>

          {/* Start Time */}
          <div>
            <label
              htmlFor="start-time"
              className="block text-sm font-medium text-foreground mb-1"
            >
              Start Time
            </label>
            <input
              id="start-time"
              type="time"
              value={formData.startTime}
              onChange={e =>
                setFormData(prev => ({ ...prev, startTime: e.target.value }))
              }
              className={cn(
                'w-full p-3 rounded-md border border-input',
                'bg-background text-foreground',
                'focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2'
              )}
            />
          </div>

          {/* Duration */}
          <div>
            <label
              htmlFor="duration"
              className="block text-sm font-medium text-foreground mb-1"
            >
              Duration
            </label>
            <select
              id="duration"
              value={formData.duration}
              onChange={e =>
                setFormData(prev => ({
                  ...prev,
                  duration: parseInt(e.target.value),
                }))
              }
              className={cn(
                'w-full p-3 rounded-md border border-input',
                'bg-background text-foreground',
                'focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2'
              )}
            >
              {durationOptions.map(option => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          </div>

          {/* Description */}
          <div>
            <label
              htmlFor="description"
              className="block text-sm font-medium text-foreground mb-1"
            >
              Description
            </label>
            <textarea
              id="description"
              value={formData.description}
              onChange={e =>
                setFormData(prev => ({ ...prev, description: e.target.value }))
              }
              rows={3}
              className={cn(
                'w-full p-3 rounded-md border border-input',
                'bg-background text-foreground resize-none',
                'focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2'
              )}
              placeholder="What makes this activity special?"
            />
          </div>


          {/* Action Buttons */}
          <div className="flex space-x-3 pt-4">
            <button
              type="button"
              onClick={onClose}
              className={cn(
                'flex-1 p-3 rounded-md border border-input',
                'text-foreground hover:bg-accent transition-colors'
              )}
            >
              Cancel
            </button>
            <button
              type="submit"
              className={cn(
                'flex-1 p-3 rounded-md bg-primary text-primary-foreground',
                'hover:bg-primary/90 transition-colors font-medium'
              )}
            >
              Save
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};