import React from "react";

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogOverlay,
  DialogPortal,
  DialogTitle,
  DialogTrigger,
} from "@radix-ui/react-dialog";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { CommonResponsesTypes, FriendsTypes } from "@discord-clone-v2/types";

import DCButton from "../Button";
import { api } from "../../api";
import { ToastDuration, useToast } from "../../hooks/useToast";
import { ClientError } from "../../utils/clientError";

export default function RemoveFriendDialog({
  id,
  username,
  trigger,
  open,
  setOpen,
}: {
  id: string;
  username: string;
  trigger: React.ReactNode;
  open: boolean;
  setOpen: (open: boolean) => void;
}) {
  const queryClient = useQueryClient();
  const toast = useToast();

  const { mutate, isPending } = useMutation({
    mutationFn: async () =>
      await api.delete<CommonResponsesTypes.MessageSuccessResponseType>(
        `/friends/invites`,
        {
          body: JSON.stringify({
            friendId: id,
          } as FriendsTypes.RemoveFriendBodyType),
        },
      ),
    onSuccess: (data) => {
      toast.success(data.message, { duration: ToastDuration.short });
      queryClient.invalidateQueries({ queryKey: ["all-friends"] });
      queryClient.invalidateQueries({ queryKey: ["chats"] });
      setOpen(false);
    },
    onError: (error: ClientError) => {
      toast.error(error.message, { duration: ToastDuration.short });
    },
  });

  function removeFriend(e?: React.MouseEvent) {
    e?.stopPropagation();

    mutate();
  }

  function closeDialog(e?: React.MouseEvent) {
    e?.stopPropagation();
    if (isPending) return;
    setOpen(false);
  }

  function handleOpenChange(open: boolean) {
    if (open) {
      setOpen(true);
    } else {
      closeDialog();
    }
  }

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogTrigger asChild>{trigger}</DialogTrigger>
      <DialogPortal>
        <DialogOverlay className="fixed inset-0 bg-black/50" />
        <DialogContent className="flex flex-col max-w-md w-full fixed top-1/2 -translate-x-1/2 left-1/2 -translate-y-1/2 bg-dc-neutral-800 rounded-md overflow-hidden">
          <DialogTitle className="p-3 text-xl font-medium text-dc-neutral-50">
            Remove Friend '{username}'
          </DialogTitle>
          <DialogDescription className="px-3 pb-3 text-dc-neutral-200">
            Are you sure that you want to remove{" "}
            <span className="font-medium">{username}</span> from your friends
          </DialogDescription>

          <footer className="p-3 bg-dc-neutral-850 flex justify-end gap-2">
            <DCButton
              variant="link"
              size="lg"
              onClick={closeDialog}
              disabled={isPending}
            >
              Cancel
            </DCButton>
            <DCButton
              variant="danger"
              size="lg"
              onClick={removeFriend}
              disabled={isPending}
            >
              Remove Friend
            </DCButton>
          </footer>
        </DialogContent>
      </DialogPortal>
    </Dialog>
  );
}
