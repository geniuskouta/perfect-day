import React, { useState, useMemo } from 'react';
import { Timeline } from '../components/timeline/Timeline';
import { ActivityModal } from '../components/activity/ActivityModal';
import {
  ShareActionIcon,
  EditActionIcon,
  CalendarIcon,
  ActivityIcon,
} from '../components/ui/Icons';
import { cn } from '../lib/utils';

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

interface PerfectDay {
  id: string;
  title: string;
  description: string;
  date: string;
  username: string;
  activities: Activity[];
}

interface PerfectDayPageProps {
  perfectDay: PerfectDay;
  isOwner?: boolean;
  isSaving?: boolean;
  onSave?: (updatedPerfectDay: PerfectDay) => void;
  onShare?: (perfectDayId: string) => void;
}

const formatDate = (dateString: string): string => {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
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

export const PerfectDayPage: React.FC<PerfectDayPageProps> = ({
  perfectDay,
  isOwner = false,
  isSaving = false,
  onSave,
  onShare,
}) => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingActivity, setEditingActivity] = useState<Activity | null>(null);
  const [selectedTimeSlot, setSelectedTimeSlot] = useState<string>('');
  const [deleteConfirmation, setDeleteConfirmation] = useState<{
    isOpen: boolean;
    activityId: string | null;
    activityName: string;
  }>({
    isOpen: false,
    activityId: null,
    activityName: '',
  });

  const activityCount = perfectDay.activities.length;
  const totalDuration = useMemo(
    () => calculateTotalDuration(perfectDay.activities),
    [perfectDay.activities]
  );
  const uniqueAreas = useMemo(
    () => getUniqueAreas(perfectDay.activities),
    [perfectDay.activities]
  );

  const handleAddActivity = ({ timeSlot }: { timeSlot: string }) => {
    if (!isOwner) return;
    setSelectedTimeSlot(timeSlot);
    setEditingActivity(null);
    setIsModalOpen(true);
  };

  const handleEditActivity = (activityId: string) => {
    if (!isOwner) return;
    const activity = perfectDay.activities.find(a => a.id === activityId);
    if (activity) {
      setEditingActivity(activity);
      setSelectedTimeSlot(activity.startTime);
      setIsModalOpen(true);
    }
  };

  const handleDeleteActivity = (activityId: string) => {
    if (!isOwner) return;
    const activity = perfectDay.activities.find(a => a.id === activityId);
    if (activity) {
      setDeleteConfirmation({
        isOpen: true,
        activityId,
        activityName: activity.name,
      });
    }
  };

  const confirmDeleteActivity = () => {
    if (!onSave || !deleteConfirmation.activityId) return;
    const updatedActivities = perfectDay.activities.filter(
      a => a.id !== deleteConfirmation.activityId
    );
    onSave({
      ...perfectDay,
      activities: updatedActivities,
    });
    setDeleteConfirmation({
      isOpen: false,
      activityId: null,
      activityName: '',
    });
  };

  const cancelDeleteActivity = () => {
    setDeleteConfirmation({
      isOpen: false,
      activityId: null,
      activityName: '',
    });
  };

  const handleSaveActivity = (activityData: Partial<Activity>) => {
    if (!onSave) return;

    let updatedActivities;
    if (editingActivity) {
      // Update existing activity
      updatedActivities = perfectDay.activities.map(activity =>
        activity.id === editingActivity.id
          ? { ...activity, ...activityData }
          : activity
      );
    } else {
      // Add new activity
      const newActivity: Activity = {
        id: `act_${Date.now()}`,
        name: activityData.name || '',
        location: activityData.location || {
          type: 'custom_text',
          name: '',
          area: 'Unknown',
        },
        startTime: activityData.startTime || selectedTimeSlot,
        duration: activityData.duration || 60,
        description: activityData.description || '',
        commentary: activityData.commentary || '',
      };
      updatedActivities = [...perfectDay.activities, newActivity];
    }

    onSave({
      ...perfectDay,
      activities: updatedActivities,
    });

    setIsModalOpen(false);
    setEditingActivity(null);
  };

  const handleShare = () => {
    if (onShare) {
      onShare(perfectDay.id);
    }
  };

  return (
    <div
      className="min-h-screen bg-background mobile responsive"
      data-testid="perfect-day-page"
    >
      <div className="max-w-md mx-auto">
        {/* Header */}
        <div className="sticky top-0 bg-background/95 backdrop-blur-sm z-40 border-b border-border">
          <div className="p-4">
            <div className="flex items-start justify-between mb-3">
              <div className="flex-1 min-w-0">
                <h1 className="text-2xl font-bold text-foreground truncate">
                  {perfectDay.title}
                </h1>
                <p className="text-muted-foreground mt-1">
                  {perfectDay.description}
                </p>
              </div>
              <div className="flex items-center space-x-2 ml-4">
                <button
                  onClick={handleShare}
                  className="p-2 rounded-md bg-secondary text-secondary-foreground hover:bg-accent transition-colors"
                  aria-label="Share"
                >
                  <ShareActionIcon size="sm" />
                </button>
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
              <div>
                <span>{formatDuration(totalDuration)} total</span>
              </div>
            </div>

            {/* Areas */}
            {uniqueAreas.length > 0 && (
              <div className="flex flex-wrap gap-1 mt-2">
                {uniqueAreas.map(area => (
                  <span
                    key={area}
                    className="px-2 py-1 bg-accent text-accent-foreground text-xs rounded-full"
                  >
                    {area}
                  </span>
                ))}
              </div>
            )}

            {/* Saving indicator */}
            {isSaving && (
              <div className="text-sm text-muted-foreground mt-2">
                Saving...
              </div>
            )}
          </div>
        </div>

        {/* Timeline */}
        <div className="p-4">
          <Timeline
            activities={perfectDay.activities}
            onAddActivity={handleAddActivity}
            onEditActivity={handleEditActivity}
            onDeleteActivity={handleDeleteActivity}
          />
        </div>

        {/* Activity Modal */}
        <ActivityModal
          isOpen={isModalOpen}
          onClose={() => {
            setIsModalOpen(false);
            setEditingActivity(null);
          }}
          onSave={handleSaveActivity}
          timeSlot={selectedTimeSlot}
          initialActivity={editingActivity}
        />

        {/* Delete Confirmation Dialog */}
        {deleteConfirmation.isOpen && (
          <div
            className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4"
            data-testid="delete-confirmation-overlay"
          >
            <div className="bg-background rounded-lg p-6 max-w-sm w-full shadow-lg">
              <h3 className="text-lg font-semibold text-foreground mb-2">
                Delete Activity
              </h3>
              <p className="text-muted-foreground mb-4">
                Are you sure you want to delete "{deleteConfirmation.activityName}"? This action cannot be undone.
              </p>
              <div className="flex space-x-3">
                <button
                  onClick={cancelDeleteActivity}
                  className={cn(
                    'flex-1 p-3 rounded-md border border-input',
                    'text-foreground hover:bg-accent transition-colors'
                  )}
                >
                  Cancel
                </button>
                <button
                  onClick={confirmDeleteActivity}
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
      </div>
    </div>
  );
};