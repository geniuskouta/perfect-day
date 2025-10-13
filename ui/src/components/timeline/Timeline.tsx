import React from 'react';
import { AddCircleIcon, TimeIcon, DeleteActionIcon } from '../ui/Icons';
import { ActivityCard } from '../activity/ActivityCard';
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

interface TimelineProps {
  activities: Activity[];
  onAddActivity: (data: { timeSlot: string }) => void;
  onEditActivity: (activityId: string) => void;
  onDeleteActivity: (activityId: string) => void;
}

const formatTime = (hour: number): string => {
  if (hour === 0) return '12:00 AM';
  if (hour < 12) return `${hour}:00 AM`;
  if (hour === 12) return '12:00 PM';
  return `${hour - 12}:00 PM`;
};

const getCurrentTimeIndicatorPosition = (): number => {
  const now = new Date();
  const hours = now.getHours();
  const minutes = now.getMinutes();
  return ((hours - 9) * 60 + minutes) / (14 * 60); // 9 AM to 11 PM = 14 hours
};

export const Timeline: React.FC<TimelineProps> = ({
  activities,
  onAddActivity,
  onEditActivity,
  onDeleteActivity,
}) => {
  const timeSlots = Array.from({ length: 15 }, (_, i) => i + 9); // 9 AM to 11 PM
  const currentTimePosition = getCurrentTimeIndicatorPosition();

  const getActivityForSlot = (hour: number) => {
    return activities.find(activity => {
      const activityHour = parseInt(activity.startTime.split(':')[0]);
      return activityHour === hour;
    });
  };

  return (
    <div className="relative w-full max-w-md mx-auto bg-background">
      {/* Current time indicator */}
      <div
        data-testid="current-time-indicator"
        className="absolute left-0 right-0 h-0.5 bg-red-500 z-10"
        style={{
          top: `${currentTimePosition * 100}%`,
          opacity: currentTimePosition >= 0 && currentTimePosition <= 1 ? 1 : 0,
        }}
      />

      <div className="space-y-2 py-4">
        {timeSlots.map(hour => {
          const timeString = `${hour.toString().padStart(2, '0')}:00`;
          const activity = getActivityForSlot(hour);

          return (
            <div key={hour} className="flex items-start space-x-3">
              {/* Time label */}
              <div className="flex items-center space-x-1 w-20 flex-shrink-0">
                <TimeIcon size="sm" className="text-muted-foreground" />
                <span className="text-sm font-medium text-foreground">
                  {formatTime(hour)}
                </span>
              </div>

              {/* Content area */}
              <div className="flex-1 min-h-[60px]">
                {activity ? (
                  /* Activity card with overlaid delete button */
                  <div className="relative w-full">
                    <ActivityCard
                      activity={activity}
                      onEdit={onEditActivity}
                      onClick={onEditActivity}
                    />
                    <button
                      className={cn(
                        'absolute top-2 right-2 p-2 rounded-md bg-background/80 backdrop-blur-sm',
                        'hover:bg-accent transition-colors',
                        'min-h-[44px] min-w-[44px] flex items-center justify-center z-10'
                      )}
                      onClick={() => onDeleteActivity(activity.id)}
                      aria-label="Delete activity"
                    >
                      <DeleteActionIcon size="sm" className="text-muted-foreground" />
                    </button>
                  </div>
                ) : (
                  /* Add activity button */
                  <button
                    className={cn(
                      'w-full h-12 rounded-lg border-2 border-dashed',
                      'border-border hover:border-primary/50 hover:bg-accent/50',
                      'flex items-center justify-center space-x-2',
                      'text-muted-foreground hover:text-foreground',
                      'transition-colors min-h-[44px] touch-manipulation'
                    )}
                    onClick={() => onAddActivity({ timeSlot: timeString })}
                    aria-label="Add activity"
                  >
                    <AddCircleIcon size="sm" />
                    <span className="text-sm">Add activity</span>
                  </button>
                )}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};