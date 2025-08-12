export type PermissionState = boolean | 'indeterminate';

export const hasPermission = (
  value: string | Array<string>,
  availablePermissions: Array<string>
): PermissionState => {
  if (!availablePermissions || availablePermissions.length === 0) return false;

  const requested = Array.isArray(value) ? value : [value];
  const perms = availablePermissions
    .filter((p): p is string => p != null)
    .map((p) => p.trim())
    .filter((p) => p.length > 0);

  for (const perm of requested) {
    // 1) Exact match (including exact wildcard pattern): granted
    if (perms.some((ap) => ap === perm)) return true;

    // 2) Higher permission via wildcard (e.g., "*" or "users.*"): indeterminate
    if (
      perms.some((ap) => {
        const cleaned = ap.replaceAll('.*', '').replaceAll('*', '');
        const matches = cleaned === perm || perm.startsWith(cleaned);
        return ap.includes('*') && matches;
      })
    ) {
      return 'indeterminate';
    }

    // 3) Non-wildcard prefix match (legacy behavior): granted
    if (
      perms.some((ap) => {
        const cleaned = ap.replaceAll('.*', '').replaceAll('*', '');
        return cleaned === perm || perm.startsWith(cleaned);
      })
    ) {
      return true;
    }
  }

  return false;
};
