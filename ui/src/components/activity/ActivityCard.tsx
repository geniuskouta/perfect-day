import React from 'react';
import { LocationIcon, TimeIcon } from '../ui/Icons';
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

interface ActivityCardProps {
  activity: Activity;
  onEdit: (activityId: string) => void;
  onClick: (activityId: string) => void;
}

const formatTime = (timeString: string): string => {
  const [hours, minutes] = timeString.split(':').map(Number);
  const period = hours >= 12 ? 'PM' : 'AM';
  const displayHours = hours === 0 ? 12 : hours > 12 ? hours - 12 : hours;
  return `${displayHours}:${minutes.toString().padStart(2, '0')} ${period}`;
};

const formatDuration = (minutes: number): string => {
  const hours = Math.floor(minutes / 60);
  const mins = minutes % 60;

  if (hours === 0) return `${mins}m`;
  if (mins === 0) return `${hours}h`;
  return `${hours}h ${mins}m`;
};

export const ActivityCard: React.FC<ActivityCardProps> = ({
  activity,
  onEdit,
  onClick,
}) => {
  const handleCardClick = (e: React.MouseEvent) => {
    onClick(activity.id);
  };

  return (
    <div
      className={cn(
        'w-full p-4 rounded-lg border text-left transition-colors',
        'bg-card hover:bg-accent border-border cursor-pointer',
        'focus-within:ring-2 focus-within:ring-ring focus-within:ring-offset-2',
        'touch-manipulation'
      )}
      onClick={handleCardClick}
      role="button"
      tabIndex={0}
      aria-label={`${activity.name} activity`}
      onKeyDown={(e) => {
        if (e.key === 'Enter' || e.key === ' ') {
          onClick(activity.id);
        }
      }}
    >
      <div className="space-y-3">
        {/* Header with title */}
        <div>
          <h3 className="font-semibold text-card-foreground text-base">
            {activity.name}
          </h3>
        </div>

        {/* Location */}
        <div className="flex items-center space-x-2">
          <LocationIcon
            size="sm"
            className="text-muted-foreground"
            data-testid="location-icon"
          />
          <div className="flex-1 min-w-0">
            <p className="text-sm text-card-foreground truncate">
              {activity.location.name}
            </p>
            <p className="text-xs text-muted-foreground">
              {activity.location.area}
            </p>
          </div>
        </div>

        {/* Time and duration */}
        <div className="flex items-center space-x-2">
          <TimeIcon
            size="sm"
            className="text-muted-foreground"
            data-testid="time-icon"
          />
          <div className="flex items-center space-x-2">
            <span className="text-sm font-medium text-card-foreground">
              {formatTime(activity.startTime)}
            </span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-sm text-muted-foreground">
              {formatDuration(activity.duration)}
            </span>
          </div>
        </div>

        {/* Description */}
        {activity.description && (
          <div className="pt-1">
            <p className="text-sm text-muted-foreground leading-relaxed">
              {activity.description}
            </p>
          </div>
        )}

        {/* Commentary */}
        {activity.commentary && (
          <div className="pt-1 border-t border-border">
            <p className="text-xs text-muted-foreground italic leading-relaxed">
              "{activity.commentary}"
            </p>
          </div>
        )}
      </div>
    </div>
  );
};