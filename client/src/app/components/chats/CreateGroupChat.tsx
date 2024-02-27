import React from 'react';
import {
  Popover,
  PopoverTrigger,
  PopoverContent,
  PopoverPortal,
} from '@radix-ui/react-popover';
import { Checkbox, CheckboxIndicator } from '@radix-ui/react-checkbox';
import { ChatsTypes } from '@discord-clone-v2/types';
import { CheckIcon, PlusIcon } from 'lucide-react';
import { useAllFriends } from '../../hooks/reactQuery/useAllFriends';
import DCButton from '../Button';
import { cn } from '../../utils/cn';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '../../utils/api';
import { useToast } from '../../hooks/useToast';

function UserCheckboxItem({
  username,
  avatar,
  onCheck,
  selected,
}: {
  id: string;
  email: string;
  username: string;
  selected: boolean;
  onCheck: (val: boolean) => void;
  avatar?: string;
}) {
  return (
    <label className='flex justify-between items-center cursor-pointer py-1 hover:bg-dc-neutral-850 px-2 -mx-2 rounded-sm'>
      <div className='flex items-center gap-2'>
        {/* TODO: replace with avatar of user in future */}
        <div className='w-8 h-8 rounded-full bg-dc-neutral-800' />
        <p className='truncate max-w-[200px]'>{username}</p>
      </div>
      <Checkbox
        onCheckedChange={onCheck}
        checked={selected}
        className='p-1 border-[1.5px] border-dc-neutral-600 rounded-md data-[state="checked"]:border-dc-purple-500 data-[state="checked"]:text-dc-purple-500'
      >
        <div className='w-4 h-4 flex'>
          <CheckboxIndicator asChild>
            <CheckIcon className='w-full h-full' strokeWidth={2.5} />
          </CheckboxIndicator>
        </div>
      </Checkbox>
    </label>
  );
}

export default function CreateGroupChat() {
  const { data } = useAllFriends();
  const toast = useToast();
  const queryClient = useQueryClient();

  const lastScrollPosition = React.useRef(0);

  const [selectedIds, setSelectedIds] = React.useState<string[]>([]);
  const [open, setOpen] = React.useState(false);
  const [showBottomBorder, setShowBottomBorder] = React.useState(false);

  const { mutate, isPending } = useMutation({
    mutationFn: async () =>
      api.post<ChatsTypes.CreateChatWithUsersSuccessResponseType>('/chats', {
        body: JSON.stringify({
          userIds: selectedIds,
        } as ChatsTypes.CreateChatWithUsersBodyType),
      }),
    onSuccess: () => {
      setSelectedIds([]);
      setOpen(false);
      toast.success('Chat created successfully!');
      queryClient.invalidateQueries({ queryKey: ['chats'] });
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });

  function handleScroll(e: React.UIEvent<HTMLDivElement, UIEvent>) {
    const target = e.currentTarget;

    if (target.scrollTop > lastScrollPosition.current) {
      console.log('scrolling down');
      setShowBottomBorder(true);
    } else {
      console.log('scrolling up');
    }

    if (target.scrollTop === 0) {
      setShowBottomBorder(false);
    }

    lastScrollPosition.current = target.scrollTop;
  }

  function toggleUserSelection(val: boolean, userId: string) {
    if (val) {
      setSelectedIds((prev) => [...prev, userId]);
    } else {
      setSelectedIds((prev) => prev.filter((id) => id !== userId));
    }
  }

  function createGroupChat() {
    mutate();
  }

  function changeOpenState(val: boolean) {
    if (isPending) return;

    setOpen(val);
  }

  if (!data) return null;

  return (
    <Popover open={open} onOpenChange={changeOpenState}>
      <PopoverTrigger>
        <button>
          <PlusIcon className='w-4 h-4 text-dc-neutral-300' />
        </button>
      </PopoverTrigger>

      <PopoverPortal>
        <PopoverContent
          className='bg-dc-neutral-900 rounded-md border border-dc-neutral-1000 w-[420px] flex flex-col text-dc-neutral-50 shadow-md'
          side='bottom'
          sideOffset={8}
          align='start'
        >
          <div
            className={cn(
              'flex flex-col px-4 pt-4 border-b transition-colors',
              showBottomBorder
                ? 'border-b-dc-neutral-1000'
                : 'border-b-transparent',
            )}
          >
            <h2 className='text-lg font-semibold text-dc-neutral-50'>
              Choose friends
            </h2>
            <p className='text-xs font-medium text-dc-neutral-300 pb-4'>
              You can add up to {data?.friends.length} friends
            </p>
          </div>
          <div
            className='flex flex-col px-4 overflow-y-auto max-h-36'
            onScroll={handleScroll}
          >
            <div className='grid pb-2 gap-1'>
              {data.friends.map((friend, i) => (
                <UserCheckboxItem
                  {...friend}
                  key={i}
                  onCheck={(val) => toggleUserSelection(val, friend.id)}
                  selected={selectedIds.includes(friend.id)}
                />
              ))}
            </div>
          </div>
          <div className='p-4 relative'>
            <span className='absolute top-0 left-0 flex w-full px-2'>
              <span className='h-[1px] flex-grow bg-dc-neutral-850'></span>
            </span>
            <DCButton onClick={createGroupChat} className='w-full'>
              Create a group chat
            </DCButton>
          </div>
        </PopoverContent>
      </PopoverPortal>
    </Popover>
  );
}
