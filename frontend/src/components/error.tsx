import { AlertCircleIcon } from 'lucide-react';

import { Alert, AlertDescription, AlertTitle } from './ui/alert';

export default function ErrorMessage({ error }: { error: Error }) {
  return (
    <div className="flex flex-col w-full h-full items-center justify-center">
      <Alert variant="destructive" className="max-w-lg">
        <AlertCircleIcon />
        <AlertTitle>{error.name}</AlertTitle>
        <AlertDescription>{error.message}</AlertDescription>
      </Alert>
    </div>
  );
}
