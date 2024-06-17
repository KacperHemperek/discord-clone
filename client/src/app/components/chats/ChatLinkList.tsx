import { Link, useMatch } from "react-router-dom";
import { AlertTriangleIcon, RefreshCwIcon, User } from "lucide-react";
import { useChats } from "@app/api";
import { useRandomVariant } from "@app/hooks/useRandomVariant.ts";

import { cn } from "../../utils/cn";
import CreateGroupChat from "./CreateGroupChat";
import { useAuth } from "../../context/AuthProvider";
import ChatLink from "./ChatLink";
import Button from "@app/components/Button.tsx";
import { useNotifications } from "@app/context/NotificationsProvider.tsx";
import { NotificationBadge } from "@app/components/NotificationBadge.tsx";

function FriendsLink() {
  const { friendRequestNotifications } = useNotifications();
  const link = "/home/friends";

  const linkToMatch = "/home/friends/*";

  const match = useMatch(linkToMatch);

  const hasNotifications =
    Boolean(friendRequestNotifications?.some((n) => !n.seen)) && !match;

  return (
    <Link
      to={link}
      className={cn(
        "w-52 p-2 rounded-sm relative text-dc-neutral-200 hover:text-dc-neutral-50 hover:bg-dc-neutral-900 flex flex-col gap-1 cursor-pointer transition-colors duration-100 mb-6",
        !!match && "bg-dc-neutral-850",
      )}
    >
      {hasNotifications && <NotificationBadge className="top-0.5 left-0.5" />}
      <h3 className="font-medium text-lg flex gap-2 items-center">
        <User size={20} />
        Friends
      </h3>
    </Link>
  );
}

export default function ChatLinkList() {
  const { data, isLoading, showLoading, error, refetch } = useChats();
  const { user } = useAuth();

  return (
    <div className=" bg-dc-neutral-950 max-h-screen flex flex-col">
      <div className="pt-2 px-2">
        <FriendsLink />
      </div>
      <div className="flex justify-between pl-4 pr-2 pb-2">
        <h4 className="uppercase text-xs font-semibold tracking-[0.02em] text-dc-neutral-300">
          Private messages
        </h4>
        <CreateGroupChat />
      </div>

      <div
        className={cn(
          "flex flex-col gap-1 flex-grow px-2 overflow-x-hidden overflow-y-auto",
          showLoading && isLoading && "overflow-y-hidden",
        )}
      >
        {error && (
          <div className="pt-12 flex flex-col items-center max-w-44 text-center mx-auto">
            <AlertTriangleIcon size={42} className="text-dc-red-500" />
            <h3 className="text-dc-red-500">Something went wrong!</h3>
            <p className="text-dc-neutral-300 text-sm pb-2">
              {error.code === 500
                ? "Could not retrieve your chats, you can try again"
                : error.message}
            </p>
            <Button
              size="sm"
              variant="primary"
              className="flex gap-1 items-center"
              onClick={() => refetch()}
            >
              <span>Try again</span> <RefreshCwIcon size={16} />
            </Button>
          </div>
        )}

        {isLoading &&
          showLoading &&
          [...Array(20).keys()].map((v) => (
            <ChatLinkSkeleton key={`chat__link_skeleton__${v}`} />
          ))}

        {data &&
          !isLoading &&
          !error &&
          data.chats.map((chat) => (
            <ChatLink
              type={chat.type}
              name={chat.name}
              id={chat.id}
              key={chat.id}
              users={chat.members}
            />
          ))}
      </div>
      <div className="bg-dc-neutral-1000/70 p-2 flex gap-3 w-full">
        <div className="w-10 h-10 min-w-[2.5rem] bg-dc-neutral-800 rounded-full" />
        <div className="min-w-0">
          <p className="truncate">{user?.username}</p>{" "}
          <p className="text-xs truncate text-dc-neutral-300">{user?.email}</p>{" "}
        </div>
      </div>
    </div>
  );
}

function ChatLinkSkeleton() {
  const variants = {
    0: {
      topText: "w-24",
      bottomText: "w-32",
    },
    1: {
      topText: "w-28",
      bottomText: "w-20",
    },
    2: {
      topText: "w-32",
      bottomText: "w-22",
    },
  };

  const variant = useRandomVariant(variants);

  return (
    <div className="w-52 p-2 rounded-sm flex gap-3 cursor-pointer transition-colors duration-100 animate-pulse">
      <div className="flex items-center justify-center w-8 h-8 min-w-[2rem] rounded-full bg-dc-neutral-800"></div>
      <div className="flex-grow flex flex-col min-w-0 gap-1">
        <div
          className={cn(
            "skeleton-text-base bg-dc-neutral-800 my-0.5",
            variant.topText,
          )}
        />
        <div
          className={cn(
            "skeleton-text-sm bg-dc-neutral-800 my-0.5",
            variant.bottomText,
          )}
        />
      </div>
    </div>
  );
}
