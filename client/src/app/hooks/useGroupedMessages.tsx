import { ChatMessageGroupProps } from "../components/chats/ChatMessageGroup";
import { getDayAtMidnight, Time } from "../utils/dates";
import { Message } from "@app/api";

type ChatMessageWithoutSender = ChatMessageGroupProps["messages"][number];

function getMassageTimeDiffInMinutes(
  message: Message | ChatMessageWithoutSender,
  lastMessage: Message | ChatMessageWithoutSender,
) {
  const messageDate = new Date(message.createdAt);
  const lastMessageDate = new Date(lastMessage.createdAt);

  const diff = messageDate.getTime() - lastMessageDate.getTime();
  return diff / Time.minute;
}

function groupMessagesByAuthor(messages: Message[]) {
  return messages.reduce<ChatMessageGroupProps[]>((prev, message) => {
    const lastMessage = prev.at(-1);

    if (!lastMessage) {
      return [
        {
          sender: message.user,
          messages: [message],
        },
      ];
    }

    const lastMessageFromGroup = lastMessage.messages.at(-1);

    if (!lastMessageFromGroup) {
      // this should never happen but just in case we return the previous state
      return prev;
    }

    if (
      lastMessage.sender.id === message.user.id &&
      getMassageTimeDiffInMinutes(message, lastMessageFromGroup) < 30
    ) {
      lastMessage.messages.push(message);
      return prev;
    }

    return [
      ...prev,
      {
        sender: message.user,
        messages: [message],
      },
    ];
  }, []);
}

export function useGroupedMessages(messages: Message[]) {
  return messages
    .reduce<{ date: Date; messages: Message[] }[]>((prev, message) => {
      // group messages by day first
      const lastMessage = prev.at(-1);

      // if there is no last message, create a new group
      if (!lastMessage) {
        const day = getDayAtMidnight(message.createdAt);

        return [
          {
            date: day,
            messages: [message],
          },
        ];
      }

      const day = getDayAtMidnight(message.createdAt);
      const time = day.getTime();

      // if there is a last message, check if it is from the same day
      if (lastMessage.date.getTime() === time) {
        lastMessage.messages.unshift(message);
        return prev;
      } else {
        return [...prev, { date: day, messages: [message] }];
      }
    }, [])
    .map((messageGroupedByDay) => {
      const groupedMessages = groupMessagesByAuthor(
        messageGroupedByDay.messages,
      );

      return {
        date: messageGroupedByDay.date,
        messages: groupedMessages,
      };
    }, []);
}
