import { QueryClientProvider } from '@tanstack/react-query';
import { Outlet, createRootRoute } from '@tanstack/react-router';

import { Toaster } from '@/components/ui/sonner';
import { queryClient } from '@/lib/utils';
import { AuthenticationProvider } from '@/providers/authentication';
import { ThemeProvider } from '@/providers/theme';

export const Route = createRootRoute({
  component: () => (
    <ThemeProvider defaultTheme="system" defaultAppearance="threereco">
      <QueryClientProvider client={queryClient}>
        <AuthenticationProvider>
          <div className="flex flex-col w-screen h-screen text-foreground bg-background">
            <Outlet />
          </div>
        </AuthenticationProvider>
      </QueryClientProvider>
      <Toaster />
    </ThemeProvider>
  ),
});
