import { createFileRoute } from '@tanstack/react-router';

import PermissionGuard from '@/components/guards/permission';

export const Route = createFileRoute('/_auth/products/$id/')({
  component: () => (
    <PermissionGuard value="products.view" isPage={true}>
      <RouteComponent />
    </PermissionGuard>
  ),
});

function RouteComponent() {
  return <div>Hello "/products/$id/"!</div>;
}
