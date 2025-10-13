import React, { useState } from 'react';
import {
  CalendarIcon,
  ActivityIcon,
  AddCircleIcon,
  DeleteActionIcon,
} from '../components/ui/Icons';
import { DateSelectionModal } from '../components/modals/DateSelectionModal';
import type { PerfectDay, Activity } from '../models';
import { cn } from '../lib/utils';


interface PerfectDayListPageProps {
  perfectDays: PerfectDay[];
  onSelectPerfectDay: (perfectDay: PerfectDay) => void;
  onCreateNew: (date: string) => void;
  onDelete: (perfectDayId: string) => void;
  user: { username: string } | null;
}

const formatDate = (dateString: string): string => {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
};

const calculateTotalDuration = (activities: Activity[]): number => {
  return activities.reduce((total, activity) => total + activity.duration, 0);
};

const formatDuration = (minutes: number): string => {
  const hours = Math.floor(minutes / 60);
  const mins = minutes % 60;

  if (hours === 0) return `${mins}m`;
  if (mins === 0) return `${hours}h`;
  return `${hours}h ${mins}m`;
};

const getUniqueAreas = (activities: Activity[]): string[] => {
  const areas = activities.map(activity => activity.location.area);
  return [...new Set(areas)].filter(area => area && area !== 'Unknown');
};

export const PerfectDayListPage: React.FC<PerfectDayListPageProps> = ({
  perfectDays,
  onSelectPerfectDay,
  onCreateNew,
  onDelete,
  user,
}) => {
  const [isDateModalOpen, setIsDateModalOpen] = useState(false);
  const [deleteConfirmation, setDeleteConfirmation] = useState<{
    isOpen: boolean;
    perfectDayId: string | null;
    perfectDayTitle: string;
  }>({
    isOpen: false,
    perfectDayId: null,
    perfectDayTitle: '',
  });

  const handleCreateClick = () => {
    setIsDateModalOpen(true);
  };

  const handleDateSelect = (date: string) => {
    onCreateNew(date);
  };

  const handleDeleteClick = (perfectDay: PerfectDay) => {
    if (!perfectDay.id) return;

    setDeleteConfirmation({
      isOpen: true,
      perfectDayId: perfectDay.id,
      perfectDayTitle: perfectDay.title,
    });
  };

  const confirmDelete = () => {
    if (!deleteConfirmation.perfectDayId) return;
    onDelete(deleteConfirmation.perfectDayId);
    setDeleteConfirmation({
      isOpen: false,
      perfectDayId: null,
      perfectDayTitle: '',
    });
  };

  const cancelDelete = () => {
    setDeleteConfirmation({
      isOpen: false,
      perfectDayId: null,
      perfectDayTitle: '',
    });
  };
  return (
    <div className="min-h-screen bg-background">
      <div className="max-w-md mx-auto">
        {/* Header */}
        <div className="sticky top-0 bg-background/95 backdrop-blur-sm z-40 border-b border-border">
          <div className="p-4">
            <div className="flex items-center justify-between mb-3">
              <div>
                <h1 className="text-2xl font-bold text-foreground">
                  Perfect Days
                </h1>
                <p className="text-muted-foreground">
                  Welcome back, {user?.username}
                </p>
              </div>
              <button
                onClick={handleCreateClick}
                className={cn(
                  'p-3 rounded-md bg-primary text-primary-foreground',
                  'hover:bg-primary/90 transition-colors',
                  'min-h-[44px] min-w-[44px] flex items-center justify-center'
                )}
                aria-label="Create new perfect day"
              >
                <AddCircleIcon size="sm" />
              </button>
            </div>
          </div>
        </div>

        {/* Perfect Days List */}
        <div className="p-4 space-y-4">
          {perfectDays.length === 0 ? (
            /* Empty state */
            <div className="text-center py-12">
              <div className="mb-4">
                <CalendarIcon size="lg" className="mx-auto text-muted-foreground" />
              </div>
              <h3 className="text-lg font-medium text-foreground mb-2">
                No perfect days yet
              </h3>
              <p className="text-muted-foreground mb-6">
                Create your first perfect day to get started!
              </p>
              <button
                onClick={handleCreateClick}
                className={cn(
                  'px-6 py-3 rounded-md bg-primary text-primary-foreground',
                  'hover:bg-primary/90 transition-colors font-medium'
                )}
              >
                Create Perfect Day
              </button>
            </div>
          ) : (
            /* Perfect Days Cards */
            perfectDays.map(perfectDay => {
              const activityCount = perfectDay.activities.length;
              const totalDuration = calculateTotalDuration(perfectDay.activities);
              const uniqueAreas = getUniqueAreas(perfectDay.activities);

              return (
                <div
                  key={perfectDay.id}
                  className={cn(
                    'relative w-full p-4 rounded-lg border',
                    'bg-card border-border'
                  )}
                >
                  {/* Delete button positioned absolutely */}
                  <button
                    onClick={() => handleDeleteClick(perfectDay)}
                    className={cn(
                      'absolute top-2 right-2 p-2 rounded-md bg-background/80 backdrop-blur-sm',
                      'hover:bg-accent transition-colors z-10',
                      'min-h-[44px] min-w-[44px] flex items-center justify-center'
                    )}
                    aria-label="Delete perfect day"
                  >
                    <DeleteActionIcon size="sm" className="text-muted-foreground" />
                  </button>

                  {/* Clickable card content */}
                  <button
                    onClick={() => onSelectPerfectDay(perfectDay)}
                    className={cn(
                      'w-full text-left transition-colors rounded-lg',
                      'hover:bg-accent/50 p-2 -m-2', // Negative margin to maintain spacing
                      'focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2'
                    )}
                  >
                    <div className="space-y-3 pr-12"> {/* Add right padding for delete button */}
                      {/* Header */}
                      <div className="flex items-start justify-between">
                        <div className="flex-1 min-w-0">
                          <h3 className="text-lg font-semibold text-card-foreground truncate">
                            {perfectDay.title}
                          </h3>
                          {perfectDay.description && (
                            <p className="text-muted-foreground text-sm mt-1">
                              {perfectDay.description}
                            </p>
                          )}
                        </div>
                      </div>

                    {/* Metadata */}
                    <div className="flex items-center space-x-4 text-sm text-muted-foreground">
                      <div className="flex items-center space-x-1">
                        <CalendarIcon size="sm" />
                        <span>{formatDate(perfectDay.date)}</span>
                      </div>
                      <div className="flex items-center space-x-1">
                        <ActivityIcon size="sm" />
                        <span>
                          {activityCount} {activityCount === 1 ? 'activity' : 'activities'}
                        </span>
                      </div>
                      {totalDuration > 0 && (
                        <div>
                          <span>{formatDuration(totalDuration)} total</span>
                        </div>
                      )}
                    </div>

                    {/* Areas */}
                    {uniqueAreas.length > 0 && (
                      <div className="flex flex-wrap gap-1">
                        {uniqueAreas.slice(0, 3).map(area => (
                          <span
                            key={area}
                            className="px-2 py-1 bg-accent text-accent-foreground text-xs rounded-full"
                          >
                            {area}
                          </span>
                        ))}
                        {uniqueAreas.length > 3 && (
                          <span className="px-2 py-1 bg-accent text-accent-foreground text-xs rounded-full">
                            +{uniqueAreas.length - 3} more
                          </span>
                        )}
                      </div>
                    )}

                    {/* Preview of first few activities */}
                    {perfectDay.activities.length > 0 && (
                      <div className="space-y-1 pt-2 border-t border-border">
                        {perfectDay.activities.slice(0, 2).map(activity => (
                          <div key={activity.id} className="flex items-center space-x-2">
                            <span className="text-xs font-medium text-muted-foreground">
                              {activity.startTime}
                            </span>
                            <span className="text-xs text-card-foreground truncate">
                              {activity.name}
                            </span>
                          </div>
                        ))}
                        {perfectDay.activities.length > 2 && (
                          <div className="text-xs text-muted-foreground">
                            +{perfectDay.activities.length - 2} more activities
                          </div>
                        )}
                      </div>
                    )}
                    </div>
                  </button>
                </div>
              );
            })
          )}
        </div>
      </div>

      {/* Delete Confirmation Dialog */}
      {deleteConfirmation.isOpen && (
        <div
          className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4"
          data-testid="delete-confirmation-overlay"
        >
          <div className="bg-background rounded-lg p-6 max-w-sm w-full shadow-lg">
            <h3 className="text-lg font-semibold text-foreground mb-2">
              Delete "{deleteConfirmation.perfectDayTitle}"
            </h3>
            <p className="text-muted-foreground mb-4">
              Are you sure you want to delete this perfect day? This action cannot be undone.
            </p>
            <div className="flex space-x-3">
              <button
                onClick={cancelDelete}
                className={cn(
                  'flex-1 p-3 rounded-md border border-input',
                  'text-foreground hover:bg-accent transition-colors'
                )}
              >
                Cancel
              </button>
              <button
                onClick={confirmDelete}
                className={cn(
                  'flex-1 p-3 rounded-md bg-destructive text-destructive-foreground',
                  'hover:bg-destructive/90 transition-colors font-medium'
                )}
              >
                Delete
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Date Selection Modal */}
      <DateSelectionModal
        isOpen={isDateModalOpen}
        onClose={() => setIsDateModalOpen(false)}
        onSelectDate={handleDateSelect}
      />
    </div>
  );
};