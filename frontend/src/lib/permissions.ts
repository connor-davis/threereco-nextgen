import type { User } from '@/api-client';

export const hasPermission = (
  value: string | string[],
  user: User
): boolean => {
  if (!user.roles) return false;

  for (const perm of Array.isArray(value) ? value : [value]) {
    return user.roles
      .flatMap((role) => role.permissions)
      .filter((permission) => permission !== undefined)
      .map((permission) => permission.replaceAll('.*', '').replaceAll('*', ''))
      .some((permission) => permission === perm || perm.startsWith(permission));
  }

  return false;
};
