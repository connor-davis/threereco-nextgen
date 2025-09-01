import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/sign-up/business')({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/sign-up/business"!</div>;
}
