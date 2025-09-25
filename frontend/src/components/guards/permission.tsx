import type { User } from '@/api-client';
import { Label } from '@/components/ui/label';
import { hasPermission } from '@/lib/permissions';

export default function PermissionGuard({
  children,
  value,
  isPage = false,
  user,
}: {
  children: React.ReactNode;
  value: string | string[];
  isPage?: boolean;
  user?: User;
}) {
  if (!user && !isPage) return null;
  if (!user && isPage)
    return (
      <div className="flex flex-col w-full h-full items-center justify-center">
        <div className="flex flex-col w-auto h-auto gap-3">
          <Label>Access Denied</Label>
          <Label className="text-muted-foreground">
            You do not have permission to view this page.
          </Label>
        </div>
      </div>
    );

  if (
    !hasPermission(value, [
      ...(user?.permissions ?? []),
      ...(user?.roles ?? []).flatMap((role) => role.permissions),
    ]) &&
    !isPage
  )
    return null;
  if (
    !hasPermission(value, [
      ...(user?.permissions ?? []),
      ...(user?.roles ?? []).flatMap((role) => role.permissions),
    ]) &&
    isPage
  )
    return (
      <div className="flex flex-col w-full h-full items-center justify-center">
        <div className="flex flex-col w-auto h-auto gap-3">
          <Label>Access Denied</Label>
          <Label className="text-muted-foreground">
            You do not have permission to view this page.
          </Label>
        </div>
      </div>
    );

  return children;
}
