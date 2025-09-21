import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { PerfectDayPage } from './pages/PerfectDayPage';
import './App.css';

// Sample data for development
const mockPerfectDay = {
  id: 'pd1',
  title: 'Amazing Tokyo Day',
  description: 'Perfect exploration of Tokyo',
  date: '2025-01-15',
  username: 'kouta',
  activities: [
    {
      id: 'act1',
      name: 'Morning Coffee',
      location: { type: 'custom_text', name: 'Local Cafe', area: 'Shibuya' },
      startTime: '09:00',
      duration: 60,
      description: 'Start the day right with excellent coffee',
      commentary: 'Best coffee in town!',
    },
    {
      id: 'act2',
      name: 'Tsukiji Fish Market',
      location: { type: 'google_places', name: 'Tsukiji Outer Market', area: 'Tsukiji' },
      startTime: '11:00',
      duration: 120,
      description: 'Fresh sushi and seafood exploration',
      commentary: 'The tuna was incredible!',
    },
    {
      id: 'act3',
      name: 'Senso-ji Temple',
      location: { type: 'google_places', name: 'Senso-ji Temple', area: 'Asakusa' },
      startTime: '14:30',
      duration: 90,
      description: 'Historic temple visit and traditional shopping',
      commentary: 'Beautiful architecture and peaceful atmosphere',
    },
  ],
};

function App() {
  const handleSave = (updatedPerfectDay: any) => {
    console.log('Saving perfect day:', updatedPerfectDay);
    // In real app, this would call the API
  };

  const handleShare = (perfectDayId: string) => {
    console.log('Sharing perfect day:', perfectDayId);
    // In real app, this would generate share links
  };

  return (
    <Router>
      <div className="bg-background text-foreground">
        <Routes>
          <Route
            path="/"
            element={
              <PerfectDayPage
                perfectDay={mockPerfectDay}
                isOwner={true}
                onSave={handleSave}
                onShare={handleShare}
              />
            }
          />
          <Route
            path="/day/:id"
            element={
              <PerfectDayPage
                perfectDay={mockPerfectDay}
                isOwner={false}
                onShare={handleShare}
              />
            }
          />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
