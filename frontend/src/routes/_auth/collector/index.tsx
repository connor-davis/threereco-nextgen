import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_auth/collector/')({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/_auth/collector/"!</div>;
}
