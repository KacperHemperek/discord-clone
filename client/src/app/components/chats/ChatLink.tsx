import { Link, useMatch } from "react-router-dom";
import { ChatType } from "@app/api";
import { cn } from "../../utils/cn";
import { useNotifications } from "@app/context/NotificationsProvider.tsx";
import { useMemo } from "react";
import { NotificationBadge } from "@app/components/NotificationBadge.tsx";

type ChatLinkProps = {
  name: string;
  id: number;
  type: ChatType;
  users: {
    id: number;
  }[];
};

export default function ChatLink({ id, name, type, users }: ChatLinkProps) {
  const link = `/home/chats/${id}`;

  const { newMessageNotifications: notifications } = useNotifications();

  const hasNotifications = useMemo(
    () => !!notifications && notifications.some((n) => n.data.chatId === id),
    [notifications, id],
  );

  const match = useMatch(link);

  if (type === ChatType.PRIVATE) {
    return (
      <Link
        to={link}
        className={cn(
          "w-52 p-2 rounded-sm font-medium text-dc-neutral-300 hover:text-dc-neutral-50 hover:bg-dc-neutral-900 flex gap-3 cursor-pointer transition-colors duration-100 items-center",
          !!match && "bg-dc-neutral-850 text-dc-neutral-50",
        )}
      >
        {/* TODO: replace div with actual user avatar */}
        <div className="w-8 h-8 min-w-[2rem] rounded-full bg-dc-neutral-800 relative">
          {hasNotifications && (
            <NotificationBadge className="-left-0.5 -top-0.5" />
          )}
        </div>
        <p className="truncate">{name}</p>
      </Link>
    );
  }

  return (
    <Link
      to={link}
      className={cn(
        "w-52 p-2 rounded-sm text-dc-neutral-300 hover:text-dc-neutral-50 hover:bg-dc-neutral-900 flex gap-3 cursor-pointer transition-colors duration-100 relative",
        !!match && "bg-dc-neutral-850 text-dc-neutral-50",
      )}
    >
      <div className="flex relative text-dc-neutral-300 items-center justify-center w-8 h-8 min-w-[2rem] rounded-full bg-dc-neutral-800">
        {hasNotifications && (
          <span className="p-1.5 rounded-full bg-dc-red-500 absolute -top-0.5 -left-0.5" />
        )}
        {name?.charAt(0).toUpperCase()}
      </div>
      <div className="flex-grow flex flex-col min-w-0 gap-1">
        <h4 className="font-medium leading-5 truncate">{name}</h4>
        <p className="text-sm text-dc-neutral-300 leading-3">
          {users.length} members
        </p>
      </div>
    </Link>
  );
}
