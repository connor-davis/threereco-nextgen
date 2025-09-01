import { createFileRoute } from '@tanstack/react-router';

import PermissionGuard from '@/components/guards/permission';

export const Route = createFileRoute('/_auth/transactions/$id/')({
  component: () => (
    <PermissionGuard value="transactions.view" isPage={true}>
      <RouteComponent />
    </PermissionGuard>
  ),
});

function RouteComponent() {
  return <div>Hello "/transactions/$id/"!</div>;
}
