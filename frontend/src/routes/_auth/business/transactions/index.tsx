import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_auth/business/transactions/')({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/_auth/business/transactions/"!</div>;
}
