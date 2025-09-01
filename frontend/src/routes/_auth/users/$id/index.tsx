import { createFileRoute } from '@tanstack/react-router';

import PermissionGuard from '@/components/guards/permission';

export const Route = createFileRoute('/_auth/users/$id/')({
  component: () => (
    <PermissionGuard value="users.view" isPage={true}>
      <RouteComponent />
    </PermissionGuard>
  ),
});

function RouteComponent() {
  return <div>Hello "/users/$id/"!</div>;
}
