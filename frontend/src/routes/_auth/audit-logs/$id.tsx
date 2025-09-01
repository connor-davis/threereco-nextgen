import { Link, createFileRoute } from '@tanstack/react-router';
import { ArrowLeftIcon } from 'lucide-react';

import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';

export const Route = createFileRoute('/_auth/audit-logs/$id')({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="flex flex-col w-full h-full bg-popover border-t p-3 gap-3">
      <div className="flex items-center justify-between w-full h-auto">
        <div className="flex items-center gap-3">
          <Link to="/audit-logs">
            <Button variant="ghost" size="icon">
              <ArrowLeftIcon className="size-4" />
            </Button>
          </Link>

          <Label className="text-lg">View Audit Log</Label>
        </div>
        <div className="flex items-center gap-3"></div>
      </div>
    </div>
  );
}
