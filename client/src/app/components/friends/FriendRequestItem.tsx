import { Check, X } from "lucide-react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api, QueryKeys } from "@app/api";
import FriendListItemButton from "./FriendItemButton";
import { useToast } from "@app/hooks/useToast.tsx";
import { ClientError } from "@app/utils/clientError.ts";

export default function FriendRequestItem({
  id,
  username,
}: {
  id: number;
  username: string;
}) {
  const queryClient = useQueryClient();
  const toast = useToast();

  const { mutate: acceptMutation } = useMutation({
    mutationKey: ["accept-friend-request", id],
    mutationFn: async () => api.post(`/friends/requests/${id}/accept`),
    onSuccess: async () => {
      await queryClient.refetchQueries({
        queryKey: QueryKeys.getPendingFriendRequests(),
      });
    },
    onError: (err: ClientError) => {
      if (err.code !== 500) {
        toast.error(err.message);
        return;
      }
      toast.error("Failed to accept friend request");
    },
  });

  const { mutate: declineMutation } = useMutation({
    mutationKey: ["friend-request-decline", id],
    mutationFn: async () => api.put(`/friends/invites/${id}/decline`),
  });

  function acceptFriendRequest() {
    acceptMutation();
  }

  function declineFriendRequest() {
    declineMutation();
  }

  return (
    <div className="relative flex w-full group">
      {/* Top Border */}
      <div className="top-0 left-0 right-0 absolute h-[1px] bg-dc-neutral-850" />
      <div className="flex justify-between items-center flex-grow py-3 px-3 -mx-3 rounded-md group-hover:bg-dc-neutral-850 transition-colors duration-100">
        {/* User information */}
        <div className="flex flex-col">
          <h2 className="font-semibold">{username}</h2>
        </div>
        {/* Action Buttons */}
        <div className="flex gap-3">
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
