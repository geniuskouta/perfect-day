import React, { useState } from 'react';
import { cn } from '../../lib/utils';

interface LoginPageProps {
  onLogin: (username: string) => void;
  loading: boolean;
  error: string | null;
}

export const LoginPage: React.FC<LoginPageProps> = ({ onLogin, loading, error }) => {
  const [username, setUsername] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (username.trim()) {
      onLogin(username.trim());
    }
  };

  return (
    <div className="min-h-screen bg-background flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-foreground mb-2">Perfect Day</h1>
          <p className="text-muted-foreground">Share your perfect day timeline</p>
        </div>

        <div className="bg-card rounded-lg shadow-lg p-6 border border-border">
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label htmlFor="username" className="block text-sm font-medium text-foreground mb-2">
                Username
              </label>
              <input
                id="username"
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                className={cn(
                  'w-full p-3 rounded-md border border-input',
                  'bg-background text-foreground',
                  'focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2',
                  'placeholder:text-muted-foreground'
                )}
                placeholder="Enter your username"
                disabled={loading}
                required
              />
            </div>

            {error && (
              <div className="p-3 bg-destructive/10 border border-destructive/20 rounded-md">
                <p className="text-sm text-destructive">{error}</p>
              </div>
            )}

            <button
              type="submit"
              disabled={loading || !username.trim()}
              className={cn(
                'w-full p-3 rounded-md font-medium transition-colors',
                'bg-primary text-primary-foreground',
                'hover:bg-primary/90 disabled:opacity-50 disabled:cursor-not-allowed'
              )}
            >
              {loading ? 'Signing in...' : 'Sign In'}
            </button>
          </form>

          <div className="mt-6 text-center text-sm text-muted-foreground">
            <p>Simple username-based authentication</p>
            <p className="mt-1">Try: <code className="bg-muted px-1 rounded">demo</code> or create a new username</p>
          </div>
        </div>
      </div>
    </div>
  );
};