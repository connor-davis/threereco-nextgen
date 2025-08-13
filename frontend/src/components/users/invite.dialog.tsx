import { postApiOrganizationsInvitesSendByEmailMutation } from '@/api-client/@tanstack/react-query.gen';
import { useMutation } from '@tanstack/react-query';
import { useRouter } from '@tanstack/react-router';
import { type ReactNode, useState } from 'react';

import { toast } from 'sonner';

import type { ErrorResponse } from '@/api-client';
import { apiClient } from '@/lib/utils';

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '../ui/alert-dialog';
import { Button } from '../ui/button';
import { Input } from '../ui/input';

export default function InviteUserByEmailDialog({
  children,
}: {
  children?: ReactNode;
}) {
  const router = useRouter();

  const [email, setEmail] = useState<string>('');

  const inviteUser = useMutation({
    ...postApiOrganizationsInvitesSendByEmailMutation({
      client: apiClient,
    }),
    onError: (error: ErrorResponse) =>
      toast.error(error.error, {
        description: error.message,
        duration: 2000,
      }),
    onSuccess: () => {
      toast.success('Success', {
        description: 'The invite has been sent successfully.',
        duration: 2000,
      });

      return router.invalidate();
    },
  });

  return (
    <AlertDialog>
      <AlertDialogTrigger>
        {children ?? (
          <Button variant="outline" size="icon">
            Invite User
          </Button>
        )}
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Invite User</AlertDialogTitle>
          <AlertDialogDescription>
            Please enter the email address of the user you want to invite.
          </AlertDialogDescription>
        </AlertDialogHeader>

        <Input
          type="text"
          autoComplete="off"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder="Enter email address"
        />

        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <AlertDialogAction
            disabled={!email || !email.includes('@') || email.length < 5}
            onClick={() =>
              inviteUser.mutate({
                path: {
                  email,
                },
              })
            }
          >
            Invite User
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
