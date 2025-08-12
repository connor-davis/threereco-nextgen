export const hasPermission = (
  value: string | Array<string>,
  availablePermissions: Array<string>
): boolean => {
  if (availablePermissions.length === 0) return false;

  for (const perm of Array.isArray(value) ? value : [value]) {
    return availablePermissions
      .filter((permission) => permission !== undefined)
      .map((permission) => permission.replaceAll('.*', '').replaceAll('*', ''))
      .some((permission) => permission === perm || perm.startsWith(permission));
  }

  return false;
};
