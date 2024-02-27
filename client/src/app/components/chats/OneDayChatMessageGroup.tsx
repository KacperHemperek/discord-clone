import React from 'react';
import ChatMessageGroup, { ChatMessageGroupProps } from './ChatMessageGroup';
import { useDateSeparatorFormatter } from '../../hooks/useDateSeparatorFormatter';

type OneDayChatMessageGroupProps = {
  date: Date;
  messages: ChatMessageGroupProps[];
};

export default function OneDayChatMessageGroup({
  date,
  messages,
}: OneDayChatMessageGroupProps) {
  const displayDate = useDateSeparatorFormatter(date);

  return (
    <section className='flex flex-col'>
      {/* Day separator */}
      <div className='flex gap-1 items-center py-4'>
        <span className='h-[1px] bg-dc-neutral-800 flex-grow' />
        <p className='text-dc-neutral-300 text-xs font-semibold'>
          {displayDate}
        </p>
        <span className='h-[1px] bg-dc-neutral-800 flex-grow' />
      </div>
      {/* Messages */}
      {messages.map((m) => (
        <ChatMessageGroup
          key={`chat__messages__group__${m.messages[0]?.id}`}
          {...m}
        />
      ))}
    </section>
  );
}
