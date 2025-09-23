import { Outlet, createFileRoute, redirect } from '@tanstack/react-router';

import type { User } from '@/api-client';
import Header from '@/components/header';
import { SidebarProvider } from '@/components/ui/sidebar';

export const Route = createFileRoute('/_auth/collector')({
  loader: async ({ context: { getUser } }) => {
    const { user } = await getUser();

    if (user && user.type !== 'collector') {
      throw redirect({
        to: '/',
      });
    }

    return { user: user as User };
  },
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <SidebarProvider>
      <div className="flex flex-col w-full h-full overflow-hidden">
        <Header />
        <Outlet />
      </div>
    </SidebarProvider>
  );
}
