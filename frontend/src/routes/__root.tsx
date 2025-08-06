import { QueryClientProvider } from '@tanstack/react-query';
import { Outlet, createRootRoute } from '@tanstack/react-router';
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools';

import AuthenticationGuard from '@/components/guards/authentication';
import Header from '@/components/header';
import AppSidebar from '@/components/sidebar/app/sidebar';
import { SidebarProvider } from '@/components/ui/sidebar';
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
            <SidebarProvider>
              <AuthenticationGuard>
                <div className="flex w-screen h-screen overflow-hidden">
                  <AppSidebar />

                  <div className="flex flex-col w-full h-full overflow-hidden">
                    <Header />
                    <Outlet />
                    {import.meta.env.MODE === 'development' && (
                      <TanStackRouterDevtools position="bottom-right" />
                    )}
                  </div>
                </div>
              </AuthenticationGuard>
            </SidebarProvider>
          </div>
        </AuthenticationProvider>
      </QueryClientProvider>
      <Toaster />
    </ThemeProvider>
  ),
});
