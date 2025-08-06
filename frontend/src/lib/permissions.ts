import type { User } from '@/api-client';

export const hasPermission = (
  value: string | string[],
  user: User
): boolean => {
  if (!user.roles) return false;

  const combinedPermissions = user.roles
    .flatMap((role) => role.permissions)
    .filter((permission) => permission !== undefined);

  let permissionsPassed = false;

  if (typeof value === 'string') {
    if (!permissionsPassed)
      permissionsPassed = combinedPermissions.includes(value);

    if (!permissionsPassed)
      permissionsPassed =
        combinedPermissions
          .filter((permission) => permission !== undefined)
          .filter((permission) => permission.endsWith('.*'))
          .map((permission) => permission.replaceAll('*', ''))
          .find((permission) => value.startsWith(permission)) !== undefined;
  } else {
    for (const listPermission of value) {
      if (!permissionsPassed)
        permissionsPassed = combinedPermissions.includes(listPermission);

      if (!permissionsPassed)
        permissionsPassed =
          combinedPermissions
            .filter((permission) => permission !== undefined)
            .filter((permission) => permission.endsWith('.*'))
            .map((permission) => permission.replaceAll('*', ''))
            .find((permission) => listPermission.startsWith(permission)) !==
          undefined;
    }
  }

  return permissionsPassed;
};
