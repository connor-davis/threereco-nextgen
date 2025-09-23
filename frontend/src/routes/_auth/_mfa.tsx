import { Outlet, createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_auth/_mfa')({
  component: RouteComponent,
});

function RouteComponent() {
  return <Outlet />;
}
