import { Outlet, createFileRoute, redirect } from '@tanstack/react-router';

import type { User } from '@/api-client';
import Header from '@/components/header';
import BusinessSidebar from '@/components/sidebar/business';
import { SidebarProvider } from '@/components/ui/sidebar';

export const Route = createFileRoute('/_auth/business')({
  loader: async ({ context: { getUser } }) => {
    const { user } = await getUser();

    if (user && user.type !== 'business') {
      throw redirect({
        to: '/',
      });
    }

    return { user: user as User };
  },
  component: RouteComponent,
});

function RouteComponent() {
  const { user } = Route.useLoaderData();

  return (
    <SidebarProvider>
      <BusinessSidebar user={user} />

      <div className="flex flex-col w-full h-full overflow-hidden">
        <Header />
        <Outlet />
      </div>
    </SidebarProvider>
  );
}
