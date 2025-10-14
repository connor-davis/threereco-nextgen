import { CalendarIcon } from 'lucide-react';
import { useState } from 'react';
import type { DateRange } from 'react-day-picker';

import { format } from 'date-fns';

import { Button } from '../ui/button';
import { Calendar } from '../ui/calendar';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '../ui/dialog';
import { Popover, PopoverContent, PopoverTrigger } from '../ui/popover';

export default function CollectionsExportDialog() {
  const [range, setRange] = useState<DateRange | undefined>({
    from: new Date(new Date().getFullYear(), new Date().getMonth(), 1),
    to: new Date(),
  });

  return (
    <Dialog>
      <DialogTrigger>
        <Button>Export</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Export Collections</DialogTitle>
          <DialogDescription>
            Please select the start and end dates for the export range.
          </DialogDescription>
        </DialogHeader>

        <Popover>
          <PopoverTrigger asChild>
            <Button variant="outline">
              <CalendarIcon />
              {range?.from && range?.to
                ? `${range.from.toLocaleDateString()} - ${range.to.toLocaleDateString()}`
                : 'June 2025'}
            </Button>
          </PopoverTrigger>
          <PopoverContent className="w-auto overflow-hidden p-0" align="end">
            <Calendar
              className="w-full"
              mode="range"
              defaultMonth={range?.from}
              selected={range}
              onSelect={setRange}
              startMonth={range?.from}
              fixedWeeks
              showOutsideDays
            />
          </PopoverContent>
        </Popover>

        <DialogFooter>
          <a
            href={`${import.meta.env.VITE_API_URL}/api/collections/export?startDate=${format(range?.from ?? new Date(new Date().getFullYear(), new Date().getMonth(), 1), 'yyyy-MM-dd')}&endDate=${format(range?.to ?? new Date(), 'yyyy-MM-dd')}`}
            target="_blank"
            rel="noreferrer"
            className="w-full"
          >
            <Button className="w-full">Export Collections</Button>
          </a>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
