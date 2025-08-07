import { Link, createFileRoute } from '@tanstack/react-router';
import { ArrowLeftIcon } from 'lucide-react';

import PermissionGuard from '@/components/guards/permission';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';

export const Route = createFileRoute('/users/create')({
  component: () => (
    <PermissionGuard value="users.create" isPage={true}>
      <RouteComponent />
    </PermissionGuard>
  ),
});

function RouteComponent() {
  return (
    <div className="flex flex-col w-full h-full bg-popover border-t p-3 gap-3">
      <div className="flex items-center justify-between w-full h-auto">
        <div className="flex items-center gap-3">
          <Link to="/users">
            <Button variant="ghost" size="icon">
              <ArrowLeftIcon className="size-4" />
            </Button>
          </Link>

          <Label className="text-lg">Create User</Label>
        </div>
        <div className="flex items-center gap-3"></div>
      </div>
    </div>
  );
}
