import { createFileRoute } from '@tanstack/react-router';

import PermissionGuard from '@/components/guards/permission';

export const Route = createFileRoute('/users/create')({
  component: () => (
    <PermissionGuard value="users.create" isPage={true}>
      <RouteComponent />
    </PermissionGuard>
  ),
});

function RouteComponent() {
  return <div>Hello "/users/create"!</div>;
}
