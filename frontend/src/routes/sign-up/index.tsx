import { Link, createFileRoute } from '@tanstack/react-router';

import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';

export const Route = createFileRoute('/sign-up/')({
  component: RouteComponent,
});

function RouteComponent() {
  return (
    <div className="flex flex-col items-center justify-center w-full h-full">
      <div className="flex flex-col w-full md:max-w-96 items-center justify-center gap-5 md:gap-10 p-5 md:p-10 m-5 md:m-10 border rounded-md bg-popover">
        <div className="flex flex-col w-full h-auto gap-5 items-center justify-center text-center">
          <img src="/logo.png" className="w-full h-20 object-contain" />

          <Label className="text-muted-foreground">
            Welcome to 3REco, please select your starting point below.
          </Label>
        </div>

        <div className="flex flex-col w-full h-auto gap-5 items-center justify-center">
          <Link to="/sign-up/individual" className="w-full">
            <Button className="w-full">I am an individual.</Button>
          </Link>
          <Link to="/sign-up/" className="w-full">
            <Button className="w-full">I am a business.</Button>
          </Link>
        </div>
      </div>
    </div>
  );
}
