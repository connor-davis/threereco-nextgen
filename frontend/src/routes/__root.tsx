import { QueryClientProvider } from '@tanstack/react-query';
import { Outlet, createRootRouteWithContext } from '@tanstack/react-router';
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools';

import { Toaster } from '@/components/ui/sonner';
import { getUser, queryClient } from '@/lib/utils';
import { ThemeProvider } from '@/providers/theme';

export const Route = createRootRouteWithContext<{
  getUser: typeof getUser;
}>()({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <ThemeProvider defaultTheme="system" defaultAppearance="threereco">
      <QueryClientProvider client={queryClient}>
        <div className="flex flex-col w-screen h-screen text-foreground bg-background">
          <Outlet />
          {import.meta.env.MODE === 'development' && (
            <TanStackRouterDevtools position="bottom-right" />
          )}
        </div>
      </QueryClientProvider>
      <Toaster />
    </ThemeProvider>
  );
}
