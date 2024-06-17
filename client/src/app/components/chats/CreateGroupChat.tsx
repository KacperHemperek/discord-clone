import React from "react";
import { PlusIcon } from "lucide-react";
import { CreateChatResponse, QueryKeys, useAllFriends, api } from "@app/api";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useToast } from "../../hooks/useToast";
import {
  UserSelectListContent,
  UserSelectListRoot,
  UserSelectListTrigger,
} from "@app/components/UserSelectList.tsx";

export default function CreateGroupChat() {
  const { data } = useAllFriends();
  const toast = useToast();
  const queryClient = useQueryClient();
  const [selectedIds, setSelectedIds] = React.useState<number[]>([]);

  const [open, setOpen] = React.useState(false);

  const { mutate, isPending } = useMutation({
    mutationFn: async (selectedIds: number[]) =>
      api.post<CreateChatResponse>("/chats/group", {
        body: JSON.stringify({
          userIds: selectedIds,
        }),
      }),
    onSuccess: async () => {
      setOpen(false);
      toast.success("Chat created successfully!");
      await queryClient.invalidateQueries({
        queryKey: QueryKeys.getAllChats(),
      });
    },
    onError: (error) => {
      toast.error(error.message);
    },
    onSettled: () => {
      setSelectedIds([]);
    },
  });

  function changeOpenState(val: boolean) {
    if (isPending) return;

    setOpen(val);
  }

  if (!data) return null;

  return (
    <UserSelectListRoot open={open} onOpenChange={changeOpenState}>
      <UserSelectListTrigger>
        <PlusIcon className="w-4 h-4 text-dc-neutral-300" />
      </UserSelectListTrigger>
      {open && (
        <UserSelectListContent
          selectedIds={selectedIds}
          setSelectedIds={setSelectedIds}
          users={data.friends}
          onSubmit={(sIds) => mutate(sIds)}
          submitLabel="Create Group Chat"
        />
      )}
    </UserSelectListRoot>
  );
}
