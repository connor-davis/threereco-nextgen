import { createFileRoute } from '@tanstack/react-router';

import PermissionGuard from '@/components/guards/permission';

export const Route = createFileRoute('/_auth/roles/$id/')({
  component: () => (
    <PermissionGuard value="roles.view" isPage={true}>
      <RouteComponent />
    </PermissionGuard>
  ),
});

function RouteComponent() {
  return <div>Hello "/roles/$id/"!</div>;
}
