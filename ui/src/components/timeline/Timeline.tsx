import React from 'react';
import { AddCircleIcon, TimeIcon } from '../ui/Icons';
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

const formatDuration = (minutes: number): string => {
  const hours = Math.floor(minutes / 60);
  const mins = minutes % 60;

  if (hours === 0) return `${mins}m`;
  if (mins === 0) return `${hours}h`;
  return `${hours}h ${mins}m`;
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
                  /* Activity card */
                  <button
                    className={cn(
                      'w-full p-3 rounded-lg border text-left transition-colors',
                      'bg-card hover:bg-accent border-border',
                      'min-h-[44px] touch-manipulation'
                    )}
                    onClick={() => onEditActivity(activity.id)}
                  >
                    <div className="space-y-1">
                      <h3 className="font-medium text-card-foreground">
                        {activity.name}
                      </h3>
                      <p className="text-sm text-muted-foreground">
                        {activity.location.name}
                      </p>
                      <p className="text-xs text-muted-foreground">
                        {activity.location.area}
                      </p>
                      <div className="flex items-center justify-between">
                        <span className="text-xs text-muted-foreground">
                          {formatDuration(activity.duration)}
                        </span>
                      </div>
                      {activity.description && (
                        <p className="text-xs text-muted-foreground mt-1">
                          {activity.description}
                        </p>
                      )}
                    </div>
                  </button>
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