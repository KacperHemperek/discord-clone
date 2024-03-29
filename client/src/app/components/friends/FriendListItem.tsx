import React from "react";
import { MessageCircle, Trash } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import FriendListItemButton from "./FriendItemButton";
import RemoveFriendDialog from "./RemoveFriendDialog";
import { api, CreateChatResponse, QueryKeys } from "@app/api";
import { ClientError } from "../../utils/clientError";
import { useToast } from "../../hooks/useToast";

export default function FriendListItem({
  id,
  username,
}: {
  id: number;
  username: string;
}) {
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const toast = useToast();
  const [open, setOpen] = React.useState(false);

  const { mutate: createChatAndRedirect, isPending } = useMutation({
    mutationFn: async () =>
      api.post<CreateChatResponse>("/chats/private", {
        body: JSON.stringify({
          userId: id,
        }),
      }),

    onError: (error: ClientError) => {
      toast.error(error.message);
    },
    onSuccess: async (data) => {
      const { chatId } = data;
      await queryClient.invalidateQueries({
        queryKey: QueryKeys.getAllChats(),
      });
      navigate(`/home/chats/${chatId}`);
    },
  });

  return (
    <div
      onClick={() => !isPending && createChatAndRedirect()}
      className="relative flex w-full group cursor-pointer"
    >
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
            disabled={isPending}
            icon={<MessageCircle size={20} />}
            onClick={(e) => {
              e.stopPropagation();
              if (!isPending) {
                createChatAndRedirect();
              }
            }}
          />
          <RemoveFriendDialog
            open={open}
            setOpen={setOpen}
            userId={id}
            username={username}
            trigger={
              <FriendListItemButton
                disabled={isPending}
                icon={<Trash size={20} />}
                onClick={(e) => e.stopPropagation()}
              />
            }
          />
        </div>
      </div>
    </div>
  );
}
