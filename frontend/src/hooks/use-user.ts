import { useRouter } from '@tanstack/react-router';
import { useEffect, useState } from 'react';

import { type User, getApiAuthenticationCheck } from '@/api-client';
import { apiClient } from '@/lib/utils';

export default function useUser() {
  const router = useRouter();
  const [user, setUser] = useState<User | false>(false);

  useEffect(() => {
    const disposeable = setTimeout(async () => {
      const response = await getApiAuthenticationCheck({
        client: apiClient,
      });

      const { data, error } = response;

      if (error) {
        setUser(false);

        return router.invalidate();
      }

      setUser(data!.item as User);

      return router.invalidate();
    }, 0);

    return () => clearTimeout(disposeable);
  }, []);

  return user;
}
