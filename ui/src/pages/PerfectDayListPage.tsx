import React, { useState } from 'react';
import {
  CalendarIcon,
  ActivityIcon,
  AddCircleIcon,
  EditActionIcon,
} from '../components/ui/Icons';
import { DateSelectionModal } from '../components/modals/DateSelectionModal';
import type { PerfectDay, Activity } from '../models';
import { cn } from '../lib/utils';


interface PerfectDayListPageProps {
  perfectDays: PerfectDay[];
  onSelectPerfectDay: (perfectDay: PerfectDay) => void;
  onCreateNew: (date: string) => void;
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
  user,
}) => {
  const [isDateModalOpen, setIsDateModalOpen] = useState(false);

  const handleCreateClick = () => {
    setIsDateModalOpen(true);
  };

  const handleDateSelect = (date: string) => {
    onCreateNew(date);
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
                <button
                  key={perfectDay.id}
                  onClick={() => onSelectPerfectDay(perfectDay)}
                  className={cn(
                    'w-full p-4 rounded-lg border text-left transition-colors',
                    'bg-card hover:bg-accent border-border',
                    'focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2'
                  )}
                >
                  <div className="space-y-3">
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
                      <EditActionIcon size="sm" className="text-muted-foreground ml-2 flex-shrink-0" />
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
              );
            })
          )}
        </div>
      </div>

      {/* Date Selection Modal */}
      <DateSelectionModal
        isOpen={isDateModalOpen}
        onClose={() => setIsDateModalOpen(false)}
        onSelectDate={handleDateSelect}
      />
    </div>
  );
};