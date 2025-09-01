import { Outlet, createFileRoute } from '@tanstack/react-router';
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools';

import AuthenticationGuard from '@/components/guards/authentication';
import Header from '@/components/header';
import AppSidebar from '@/components/sidebar/app/sidebar';
import { SidebarProvider } from '@/components/ui/sidebar';

export const Route = createFileRoute('/_auth')({
  component: RouteComponent,
});

function RouteComponent() {
  return (
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
  );
}
