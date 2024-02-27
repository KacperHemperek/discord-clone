import { Link, useMatch } from 'react-router-dom';
import { User } from 'lucide-react';
import { useChats } from '../../hooks/reactQuery/useChats';
import { cn } from '../../utils/cn';
import CreateGroupChat from './CreateGroupChat';
import { useAuth } from '../../context/AuthProvider';
import ChatLink from './ChatLink';

function FriendsLink() {
  const link = '/home/friends';

  const linkToMatch = '/home/friends/*';

  const match = useMatch(linkToMatch);

  return (
    <Link
      to={link}
      className={cn(
        'w-52 p-2 rounded-sm hover:bg-dc-neutral-900 flex flex-col gap-1 cursor-pointer transition-colors duration-100 mb-6',
        !!match && 'bg-dc-neutral-850',
      )}
    >
      <h3 className='font-medium text-lg flex gap-2 items-center'>
        <User size={20} />
        Friends
      </h3>
    </Link>
  );
}

export default function ChatLinkList() {
  const { data } = useChats();
  const { user } = useAuth();

  return (
    <div className=' bg-dc-neutral-950 max-h-screen flex flex-col'>
      <div className='pt-2 px-2'>
        <FriendsLink />
      </div>
      <div className='flex justify-between pl-4 pr-2 pb-2'>
        <h4 className='uppercase text-xs font-semibold tracking-[0.02em] text-dc-neutral-300'>
          Private messages
        </h4>
        <CreateGroupChat />
      </div>
      <div className='flex flex-col gap-1 flex-grow px-2 overflow-y-auto'>
        {data?.chats.map((chat) => <ChatLink {...chat} key={chat.id} />)}
      </div>
      <div className='bg-dc-neutral-1000/70 p-2 flex gap-3 w-full'>
        <div className='w-10 h-10 min-w-[2.5rem] bg-dc-neutral-800 rounded-full' />
        <div className='min-w-0'>
          <p className='truncate'>{user?.username}</p>{' '}
          <p className='text-xs truncate text-dc-neutral-300'>{user?.email}</p>{' '}
        </div>
      </div>
    </div>
  );
}
