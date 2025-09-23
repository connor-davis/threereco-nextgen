import { createFileRoute } from '@tanstack/react-router';

import { Label } from '@/components/ui/label';

export const Route = createFileRoute('/_auth/admin/')({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="flex flex-col w-full h-full items-center justify-center">
      <Label className="text-2xl text-muted-foreground">Welcome to 3REco</Label>
    </div>
  );
}
