// Material Icons for Perfect Day PWA
import AddIcon from '@mui/icons-material/Add';
import ScheduleIcon from '@mui/icons-material/Schedule';
import LocationOnIcon from '@mui/icons-material/LocationOn';
import CameraAltIcon from '@mui/icons-material/CameraAlt';
import CloseIcon from '@mui/icons-material/Close';
import ShareIcon from '@mui/icons-material/Share';
import CalendarTodayIcon from '@mui/icons-material/CalendarToday';
import LocalActivityIcon from '@mui/icons-material/LocalActivity';
import SearchIcon from '@mui/icons-material/Search';
import FilterListIcon from '@mui/icons-material/FilterList';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import HomeIcon from '@mui/icons-material/Home';
import ExploreIcon from '@mui/icons-material/Explore';
import PersonIcon from '@mui/icons-material/Person';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import SaveIcon from '@mui/icons-material/Save';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';

// Icon wrapper with consistent sizing
interface IconProps {
  className?: string;
  size?: 'sm' | 'md' | 'lg';
}

const iconSizes = {
  sm: 'text-base', // 16px
  md: 'text-xl', // 20px
  lg: 'text-2xl', // 24px
};

// Timeline icons
export const AddCircleIcon = ({ className = '', size = 'md' }: IconProps) => (
  <AddIcon className={`${iconSizes[size]} ${className}`} />
);

export const TimeIcon = ({ className = '', size = 'md' }: IconProps) => (
  <ScheduleIcon className={`${iconSizes[size]} ${className}`} data-testid="time-icon" />
);

// Activity modal icons
export const LocationIcon = ({ className = '', size = 'md' }: IconProps) => (
  <LocationOnIcon className={`${iconSizes[size]} ${className}`} data-testid="location-icon" />
);

export const CameraIcon = ({ className = '', size = 'md' }: IconProps) => (
  <CameraAltIcon className={`${iconSizes[size]} ${className}`} data-testid="camera-icon" />
);

export const CloseModalIcon = ({ className = '', size = 'md' }: IconProps) => (
  <CloseIcon className={`${iconSizes[size]} ${className}`} />
);

// Perfect day card icons
export const ShareActionIcon = ({ className = '', size = 'md' }: IconProps) => (
  <ShareIcon className={`${iconSizes[size]} ${className}`} />
);

export const CalendarIcon = ({ className = '', size = 'md' }: IconProps) => (
  <CalendarTodayIcon className={`${iconSizes[size]} ${className}`} />
);

export const ActivityIcon = ({ className = '', size = 'md' }: IconProps) => (
  <LocalActivityIcon className={`${iconSizes[size]} ${className}`} />
);

// Discovery icons
export const SearchActionIcon = ({
  className = '',
  size = 'md',
}: IconProps) => <SearchIcon className={`${iconSizes[size]} ${className}`} />;

export const FilterIcon = ({ className = '', size = 'md' }: IconProps) => (
  <FilterListIcon className={`${iconSizes[size]} ${className}`} />
);

// Share dialog icons
export const CopyIcon = ({ className = '', size = 'md' }: IconProps) => (
  <ContentCopyIcon className={`${iconSizes[size]} ${className}`} />
);

// Navigation icons
export const HomeNavIcon = ({ className = '', size = 'md' }: IconProps) => (
  <HomeIcon className={`${iconSizes[size]} ${className}`} />
);

export const ExploreNavIcon = ({ className = '', size = 'md' }: IconProps) => (
  <ExploreIcon className={`${iconSizes[size]} ${className}`} />
);

export const ProfileIcon = ({ className = '', size = 'md' }: IconProps) => (
  <PersonIcon className={`${iconSizes[size]} ${className}`} />
);

// Action icons
export const EditActionIcon = ({ className = '', size = 'md' }: IconProps) => (
  <EditIcon className={`${iconSizes[size]} ${className}`} />
);

export const DeleteActionIcon = ({
  className = '',
  size = 'md',
}: IconProps) => <DeleteIcon className={`${iconSizes[size]} ${className}`} />;

export const SaveActionIcon = ({ className = '', size = 'md' }: IconProps) => (
  <SaveIcon className={`${iconSizes[size]} ${className}`} />
);

export const ArrowBackActionIcon = ({ className = '', size = 'md' }: IconProps) => (
  <ArrowBackIcon className={`${iconSizes[size]} ${className}`} />
);

// Export all icons for easy importing
export {
  AddIcon,
  ScheduleIcon,
  LocationOnIcon,
  CameraAltIcon,
  CloseIcon,
  ShareIcon,
  CalendarTodayIcon,
  LocalActivityIcon,
  SearchIcon,
  FilterListIcon,
  ContentCopyIcon,
  HomeIcon,
  ExploreIcon,
  PersonIcon,
  EditIcon,
  DeleteIcon,
  SaveIcon,
};
