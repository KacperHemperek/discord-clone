import {
  api,
  ChatType,
  GetAllChats,
  GetChat,
  QueryKeys,
  SuccessMessageResponse,
} from "@app/api";
import React from "react";
import { useChatId } from "@app/hooks/useChatId.ts";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useToast } from "@app/hooks/useToast.tsx";
import { ClientError } from "@app/utils/clientError.ts";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { UserPlus } from "lucide-react";

const nameChangeSchema = z.object({
  newName: z.string(),
});

type NameChangeFormData = z.infer<typeof nameChangeSchema>;

function NameChangeElement({
  name,
  chatId,
  disabled,
}: {
  name: string;
  chatId: number;
  disabled?: boolean;
}) {
  const [oldChatName, setOldChatName] = React.useState<string | null>(null);

  const queryClient = useQueryClient();
  const toast = useToast();
  const { mutate, isPending } = useMutation<
    SuccessMessageResponse,
    ClientError,
    NameChangeFormData
  >({
    mutationFn: async (data: NameChangeFormData) =>
      api.put<SuccessMessageResponse>(`/chats/${chatId}/update-name`, {
        body: JSON.stringify(data),
      }),
    onMutate: async (inputData) => {
      updateChatName(inputData.newName);
      const oldData = queryClient.getQueryData<GetAllChats>(
        QueryKeys.getAllChats(),
      );

      if (!oldData) return;

      const newChatList = oldData.chats.map((c) => {
        if (c.id === chatId) {
          return {
            ...c,
            name: inputData.newName,
          };
        }

        return c;
      });

      const newData: GetAllChats = {
        chats: newChatList,
      };

      queryClient.setQueryData(QueryKeys.getAllChats(), newData);
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: QueryKeys.getAllChats(),
      });
    },
    onError: async (err) => {
      if (err.code === 400) {
        toast.error("New chat name must be between 6 and 32 characters");
      } else {
        toast.error("Unknown error when changing chat name");
      }
      await queryClient.invalidateQueries({
        queryKey: QueryKeys.getAllChats(),
      });

      if (!oldChatName) return;

      updateChatName(oldChatName);
    },
  });

  const form = useForm<NameChangeFormData>({
    resolver: zodResolver(nameChangeSchema),
    defaultValues: {
      newName: name,
    },
  });

  React.useEffect(() => {
    form.setValue("newName", name);
  }, [chatId, name, form]);

  function updateChatName(name: string) {
    queryClient.setQueryData(
      QueryKeys.getChat(chatId),
      (oldData: GetChat): GetChat => ({
        ...oldData,
        name,
      }),
    );
  }

  return (
    <form
      onSubmit={form.handleSubmit((data) => {
        setOldChatName(name);
        mutate(data);
      })}
    >
      <input
        type="text"
        {...form.register("newName")}
        className="px-2 rounded-md bg-transparent border ring-0 border-transparent font-semibold hover:border-dc-neutral-950 focus:bg-dc-neutral-950 focus:border-dc-neutral-1000 disabled:hover:bg-transparent disabled:hover:border-transparent disabled:cursor-text"
        size={form.watch("newName").length}
        disabled={disabled || isPending}
      />
      <button className="hidden" type="submit" />
    </form>
  );
}

type PrivateChatHeader = {
  name: string;
  type: ChatType;
};

export function PrivateChatHeader({ name, type }: PrivateChatHeader) {
  const chatId = useChatId();

  return (
    <nav className="border-b flex justify-between border-dc-neutral-1000 w-full p-3 gap-4">
      {!!name && !!chatId && (
        <NameChangeElement
          name={name}
          disabled={type === ChatType.PRIVATE}
          chatId={chatId}
        />
      )}

      <div className="flex gap-2 pr-16">
        <button className="text-dc-neutral-300 hover:text-dc-neutral-50 transition-colors duration-100">
          <UserPlus />
        </button>
      </div>
    </nav>
  );
}
