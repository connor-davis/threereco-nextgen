import { createFileRoute } from '@tanstack/react-router';

import PermissionGuard from '@/components/guards/permission';

export const Route = createFileRoute('/_auth/materials/$id/')({
  component: () => (
    <PermissionGuard value="materials.view" isPage={true}>
      <RouteComponent />
    </PermissionGuard>
  ),
});

function RouteComponent() {
  return <div>Hello "/materials/$id/"!</div>;
}
