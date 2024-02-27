import React from 'react';
import { Check, X } from 'lucide-react';
import { useMutation } from '@tanstack/react-query';
import { api } from '../../utils/api';
import { useFriendRequests } from '../../context/FriendRequestsProvider';
import FriendListItemButton from './FriendItemButton';
import { CommonResponsesTypes } from '@discord-clone-v2/types';

export default function FriendRequestItem({
  id,
  userId,
  username,
  avatar,
}: {
  id: string;
  username: string;
  userId: string;
  avatar?: string;
}) {
  const { removeRequest } = useFriendRequests();

  const { mutate: acceptMutation } = useMutation({
    mutationKey: ['friend-request-accept', id],
    mutationFn: async () =>
      api.put<CommonResponsesTypes.MessageSuccessResponseType>(
        `/friends/invites/${id}/accept`,
      ),
    onSuccess: () => {
      removeRequest(id);
    },
  });

  const { mutate: declineMutation } = useMutation({
    mutationKey: ['friend-request-decline', id],
    mutationFn: async () =>
      api.put<CommonResponsesTypes.MessageSuccessResponseType>(
        `/friends/invites/${id}/decline`,
      ),

    onSuccess: () => {
      removeRequest(id);
    },
  });

  function acceptFriendRequest() {
    acceptMutation();
  }

  function declineFriendRequest() {
    declineMutation();
  }

  return (
    <div className='relative flex w-full group'>
      {/* Top Border */}
      <div className='top-0 left-0 right-0 absolute h-[1px] bg-dc-neutral-850' />
      <div className='flex justify-between items-center flex-grow py-3 px-3 -mx-3 rounded-md group-hover:bg-dc-neutral-850 transition-colors duration-100'>
        {/* User information */}
        <div className='flex flex-col'>
          <h2 className='font-semibold'>{username}</h2>
        </div>
        {/* Action Buttons */}
        <div className='flex gap-3'>
          <FriendListItemButton
            onClick={acceptFriendRequest}
            icon={<Check size={20} />}
          />
          <FriendListItemButton
            onClick={declineFriendRequest}
            icon={<X size={20} />}
          />
        </div>
      </div>
    </div>
  );
}
