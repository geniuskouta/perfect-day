import React, { useState } from 'react';
import { CloseModalIcon, CalendarIcon } from '../ui/Icons';
import { cn } from '../../lib/utils';

interface DateSelectionModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSelectDate: (date: string) => void;
}

const formatDate = (date: Date): string => {
  return date.toLocaleDateString('en-US', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
};

const formatDateForInput = (date: Date): string => {
  return date.toISOString().split('T')[0];
};

export const DateSelectionModal: React.FC<DateSelectionModalProps> = ({
  isOpen,
  onClose,
  onSelectDate,
}) => {
  const [selectedDate, setSelectedDate] = useState<string>(
    formatDateForInput(new Date())
  );

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSelectDate(selectedDate);
    onClose();
  };

  const handleOverlayClick = (e: React.MouseEvent) => {
    if (e.target === e.currentTarget) {
      onClose();
    }
  };

  const getQuickDateOptions = () => {
    const today = new Date();
    const yesterday = new Date(today);
    yesterday.setDate(today.getDate() - 1);
    const dayBeforeYesterday = new Date(today);
    dayBeforeYesterday.setDate(today.getDate() - 2);

    return [
      { label: 'Today', date: formatDateForInput(today), formatted: formatDate(today) },
      { label: 'Yesterday', date: formatDateForInput(yesterday), formatted: formatDate(yesterday) },
      { label: 'Day Before Yesterday', date: formatDateForInput(dayBeforeYesterday), formatted: formatDate(dayBeforeYesterday) },
    ];
  };

  const quickOptions = getQuickDateOptions();

  if (!isOpen) return null;

  return (
    <div
      className="fixed inset-0 bg-black/50 z-50 flex items-end md:items-center justify-center p-4"
      data-testid="date-selection-overlay"
      onClick={handleOverlayClick}
    >
      <div
        className={cn(
          'bg-background shadow-lg',
          'w-full md:w-full md:max-w-md',
          'max-h-[85vh] md:max-h-[90vh] overflow-y-auto',
          'rounded-t-xl md:rounded-xl',
          'bottom-sheet slide-up'
        )}
        role="dialog"
        aria-labelledby="date-modal-title"
      >
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-border">
          <h2 id="date-modal-title" className="text-lg font-semibold text-foreground">
            When was your perfect day?
          </h2>
          <button
            onClick={onClose}
            className="p-2 hover:bg-accent rounded-md transition-colors"
            aria-label="Close"
          >
            <CloseModalIcon size="sm" />
          </button>
        </div>

        {/* Content */}
        <form onSubmit={handleSubmit} className="p-4 space-y-6">
          {/* Quick Options */}
          <div>
            <h3 className="text-sm font-medium text-foreground mb-3">Quick Options</h3>
            <div className="space-y-2">
              {quickOptions.map((option) => (
                <button
                  key={option.label}
                  type="button"
                  onClick={() => setSelectedDate(option.date)}
                  className={cn(
                    'w-full p-3 rounded-md border text-left transition-colors',
                    'focus:outline-none focus:ring-2 focus:ring-ring',
                    selectedDate === option.date
                      ? 'bg-primary text-primary-foreground border-primary'
                      : 'bg-background border-border hover:bg-accent'
                  )}
                >
                  <div className="flex items-center space-x-3">
                    <CalendarIcon size="sm" className="flex-shrink-0" />
                    <div>
                      <div className="font-medium">{option.label}</div>
                      <div className="text-sm opacity-75">{option.formatted}</div>
                    </div>
                  </div>
                </button>
              ))}
            </div>
          </div>

          {/* Custom Date */}
          <div>
            <h3 className="text-sm font-medium text-foreground mb-3">Or Choose Custom Date</h3>
            <div className="relative">
              <input
                type="date"
                value={selectedDate}
                onChange={(e) => setSelectedDate(e.target.value)}
                className={cn(
                  'w-full p-3 rounded-md border border-input',
                  'bg-background text-foreground',
                  'focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2'
                )}
              />
            </div>
            {selectedDate && !quickOptions.some(opt => opt.date === selectedDate) && (
              <p className="text-sm text-muted-foreground mt-2">
                {formatDate(new Date(selectedDate))}
              </p>
            )}
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
              Document This Day
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};