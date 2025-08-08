import { Link } from '@tanstack/react-router';

import { SidebarTrigger } from './ui/sidebar';

export default function Header() {
  return (
    <div className="flex items-center justify-between gap-3">
      <div className="flex items-center gap-3 p-3">
        <SidebarTrigger />

        <Link to="/">
          <img src="/logo.png" className="w-full h-8 object-contain" />
        </Link>
      </div>
    </div>
  );
}
