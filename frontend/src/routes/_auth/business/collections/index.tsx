import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_auth/business/collections/')({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/_auth/business/collections/"!</div>;
}
