import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_auth/business/users/')({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="flex flex-col w-full h-full overflow-hidden px-3"></div>
  );
}
