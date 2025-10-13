import { useState, useEffect } from 'react';
import { PerfectDayPage } from './pages/PerfectDayPage';
import { PerfectDayListPage } from './pages/PerfectDayListPage';
import { LoginPage } from './components/auth/LoginPage';
import { apiService } from './services/apiService';
import type { PerfectDay } from './models';
import './App.css';

function App() {
  const [user, setUser] = useState<{ username: string } | null>(null);
  const [authLoading, setAuthLoading] = useState(true);
  const [authError, setAuthError] = useState<string | null>(null);
  const [perfectDays, setPerfectDays] = useState<PerfectDay[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [currentView, setCurrentView] = useState<'list' | 'detail'>('list');
  const [selectedPerfectDay, setSelectedPerfectDay] = useState<PerfectDay | null>(null);

  // Check authentication status on app load
  useEffect(() => {
    checkAuthStatus();
  }, []);

  // Load perfect days when user is authenticated
  useEffect(() => {
    if (user) {
      loadPerfectDays();
    }
  }, [user]);

  const checkAuthStatus = async () => {
    try {
      setAuthLoading(true);
      const result = await apiService.auth.getCurrentUser();
      if (result.success && result.data) {
        setUser(result.data.user);
      }
    } catch (err) {
      // Not authenticated, which is fine
    } finally {
      setAuthLoading(false);
    }
  };

  const handleLogin = async (username: string) => {
    try {
      setAuthLoading(true);
      setAuthError(null);
      const result = await apiService.auth.login(username);
      if (result.success && result.data) {
        setUser(result.data.user);
      } else {
        setAuthError(result.error?.message || 'Login failed');
      }
    } catch (err) {
      setAuthError('Network error during login');
    } finally {
      setAuthLoading(false);
    }
  };

  // const handleLogout = async () => {
  //   try {
  //     await apiService.auth.logout();
  //     setUser(null);
  //     setPerfectDays([]);
  //   } catch (err) {
  //     console.error('Logout error:', err);
  //   }
  // };

  const loadPerfectDays = async () => {
    try {
      setLoading(true);
      const result = await apiService.perfectDays.list();
      if (result.success && result.data) {
        setPerfectDays(result.data.perfect_days || []);
      } else {
        setError(result.error?.message || 'Failed to load perfect days');
      }
    } catch (err) {
      setError('Network error while loading perfect days');
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async (updatedPerfectDay: PerfectDay) => {
    try {
      if (updatedPerfectDay.id) {
        // Update existing
        const result = await apiService.perfectDays.update(updatedPerfectDay.id, updatedPerfectDay);
        if (result.success && result.data) {
          // Update the working copy with server response
          setSelectedPerfectDay(result.data);
        }
      } else {
        // Create new
        const result = await apiService.perfectDays.create(updatedPerfectDay);
        if (result.success && result.data) {
          // Update the working copy with server response (now has ID)
          setSelectedPerfectDay(result.data);
        }
      }
    } catch (err) {
      console.error('Failed to save perfect day:', err);
    }
  };

  const handleShare = (perfectDayId: string) => {
    if (!perfectDayId) {
      console.warn('Cannot share unsaved perfect day');
      return;
    }
    const url = `${window.location.origin}/day/${perfectDayId}`;
    navigator.clipboard.writeText(url);
    console.log('Shared URL:', url);
  };

  const handleSelectPerfectDay = (perfectDay: PerfectDay) => {
    setSelectedPerfectDay(perfectDay);
    setCurrentView('detail');
  };

  const handleCreateNew = (date: string) => {
    const newPerfectDay: PerfectDay = {
      id: '',
      title: 'My Perfect Day',
      description: 'Create your perfect day timeline',
      date: date,
      username: user?.username || 'me',
      activities: [],
    };
    setSelectedPerfectDay(newPerfectDay);
    setCurrentView('detail');
  };

  const handleBackToList = async () => {
    // Refresh the list when going back
    await loadPerfectDays();
    setCurrentView('list');
    setSelectedPerfectDay(null);
  };

  // Show auth loading
  if (authLoading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <div className="text-lg font-medium text-foreground">Loading...</div>
          <div className="text-sm text-muted-foreground mt-1">Checking authentication</div>
        </div>
      </div>
    );
  }

  // Show login page if not authenticated
  if (!user) {
    return <LoginPage onLogin={handleLogin} loading={authLoading} error={authError} />;
  }

  // Show data loading or error state for authenticated users
  if (loading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <div className="text-lg font-medium text-foreground">Loading...</div>
          <div className="text-sm text-muted-foreground mt-1">Fetching your perfect days</div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <div className="text-lg font-medium text-destructive">Error</div>
          <div className="text-sm text-muted-foreground mt-1">{error}</div>
          <button
            onClick={loadPerfectDays}
            className="mt-4 px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  // Show list or detail view based on current state
  if (currentView === 'list') {
    return (
      <div className="bg-background text-foreground">
        <PerfectDayListPage
          perfectDays={perfectDays}
          onSelectPerfectDay={handleSelectPerfectDay}
          onCreateNew={handleCreateNew}
          user={user}
        />
      </div>
    );
  }

  // Detail view
  const currentPerfectDay = selectedPerfectDay || {
    id: '',
    title: 'My Perfect Day',
    description: 'Create your perfect day timeline',
    date: new Date().toISOString().split('T')[0],
    username: user?.username || 'me',
    activities: [],
  };

  return (
    <div className="bg-background text-foreground">
      <PerfectDayPage
        perfectDay={currentPerfectDay}
        isOwner={true}
        onSave={handleSave}
        onShare={handleShare}
        onBack={handleBackToList}
      />
    </div>
  );
}

export default App;
