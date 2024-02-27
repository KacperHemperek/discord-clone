import { Search } from 'lucide-react';
import React from 'react';

export default function DCSearchBar({
  setValue,
  value,
}: {
  value: string;
  setValue: React.Dispatch<React.SetStateAction<string>>;
}) {
  return (
    <div className='bg-dc-neutral-1000 px-3 py-1.5 rounded-md flex items-center gap-3 placeholder:text-dc-neutral-300'>
      <input
        type='text'
        className='bg-transparent outline-none flex-grow ring-0'
        placeholder='Search'
        value={value}
        onChange={(e) => setValue(e.target.value)}
      />
      <Search size={20} className='text-dc-neutral-300' />
    </div>
  );
}
