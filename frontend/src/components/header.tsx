import { Link } from '@tanstack/react-router';

import { Label } from './ui/label';
import { SidebarTrigger } from './ui/sidebar';

export default function Header() {
  return (
    <div className="flex items-center justify-between gap-3">
      <div className="flex items-center gap-3 p-3">
        <SidebarTrigger />

        <Link to="/">
          <Label className="text-primary font-bold text-2xl">3rEco</Label>
        </Link>
      </div>
    </div>
  );
}
