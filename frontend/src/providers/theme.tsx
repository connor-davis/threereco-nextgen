import { createContext, useContext, useEffect, useState } from 'react';

type Theme = 'light' | 'dark' | 'system';

type ThemeProviderProps = {
  children: React.ReactNode;
  defaultTheme?: Theme;
  defaultAppearance?: 'threereco';
  storageKey?: string;
};

type ThemeProviderState = {
  theme: Theme;
  appearance: 'threereco';
  setTheme: (theme: Theme) => void;
  setAppearance: (appearance: 'threereco') => void;
};

const initialState: ThemeProviderState = {
  theme: 'system',
  appearance: 'threereco',
  setTheme: () => null,
  setAppearance: () => null,
};

const ThemeProviderContext = createContext<ThemeProviderState>(initialState);

export function ThemeProvider({
  children,
  defaultTheme = 'system',
  defaultAppearance = 'threereco',
  storageKey = 'threereco-theme',
  ...props
}: ThemeProviderProps) {
  const [theme, setTheme] = useState<Theme>(
    () => (localStorage.getItem(storageKey) as Theme) || defaultTheme
  );
  const [appearance, setAppearance] = useState<'threereco'>(
    () =>
      (localStorage.getItem('threereco-appearance') as 'threereco') ||
      defaultAppearance
  );

  useEffect(() => {
    const root = window.document.documentElement;

    root.classList.remove('light', 'dark');

    if (theme === 'system') {
      const systemTheme = window.matchMedia('(prefers-color-scheme: dark)')
        .matches
        ? `dark`
        : `light`;

      root.classList.add(systemTheme);

      return;
    }

    root.classList.add(`${theme}`);
  }, [theme, appearance]);

  const value = {
    theme,
    setTheme: (theme: Theme) => {
      localStorage.setItem(storageKey, theme);
      setTheme(theme);
    },
    appearance,
    setAppearance: (appearance: 'threereco') => {
      localStorage.setItem('threereco-appearance', appearance);
      setAppearance(appearance);
    },
  };

  return (
    <ThemeProviderContext.Provider {...props} value={value}>
      {children}
    </ThemeProviderContext.Provider>
  );
}

export const useTheme = () => {
  const context = useContext(ThemeProviderContext);

  if (context === undefined)
    throw new Error('useTheme must be used within a ThemeProvider');

  return context;
};
