import { ChatsTypes } from '@discord-clone-v2/types';
import { useChatMessageDateFormatter } from '../../hooks/useChatDateFormatter';
import { formatShortTime } from '../../utils/dates';

function ChatMessage({
  message,
}: {
  message: ChatMessageGroupProps['messages'][number];
}) {
  return (
    <div className='relative group'>
      <p>{message.text}</p>
      <p className='absolute text-dc-neutral-300 font-semibold text-xs top-1 -left-11 hidden group-hover:block group-first:hidden group-first:group-hover:hidden'>
        {formatShortTime(message.createdAt)}
      </p>
    </div>
  );
}

export type ChatMessageGroupProps = {
  sender: ChatsTypes.ChatMessage['sender'];
  messages: Omit<ChatsTypes.ChatMessage, 'sender'>[];
};

export default function ChatMessageGroup({
  sender,
  messages,
}: ChatMessageGroupProps) {
  const displayDate = useChatMessageDateFormatter(messages[0]!.createdAt);

  return (
    <article className='flex gap-4 py-2 first-of-type:pt-0 last-of-type:pb-0'>
      {/* Avatar in future */}
      <div className='w-10 h-10 min-w-[2.5rem] min-h-[2.5rem] bg-dc-neutral-700 rounded-full ' />
      <div className='flex flex-col w-full'>
        <div className='flex items-end gap-2'>
          <p className='font-semibold'>{sender.username}</p>
          <p className='text-dc-neutral-300 text-xs font-medium pb-0.5'>
            {displayDate}
          </p>
        </div>
        <div className='flex flex-col'>
          {messages.map((m) => (
            <ChatMessage key={`single__chat__message__${m.id}`} message={m} />
          ))}
        </div>
      </div>
    </article>
  );
}
