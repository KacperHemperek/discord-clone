import { Link, useMatch } from "react-router-dom";
import { ChatType } from "@app/api";
import { cn } from "../../utils/cn";

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

  const match = useMatch(link);

  if (type === ChatType.PRIVATE) {
    return (
      <Link
        to={link}
        className={cn(
          "w-52 p-2 rounded-sm hover:bg-dc-neutral-900 flex gap-3 cursor-pointer transition-colors duration-100 items-center",
          !!match && "bg-dc-neutral-850",
        )}
      >
        {/* TODO: replace div with actual user avatar */}
        <div className="w-8 h-8 min-w-[2rem] rounded-full bg-dc-neutral-800" />
        <p className="truncate">{name}</p>
      </Link>
    );
  }

  return (
    <Link
      to={link}
      className={cn(
        "w-52 p-2 rounded-sm hover:bg-dc-neutral-900 flex gap-3 cursor-pointer transition-colors duration-100",
        !!match && "bg-dc-neutral-850",
      )}
    >
      <div className="flex items-center justify-center w-8 h-8 min-w-[2rem] rounded-full bg-dc-neutral-800">
        {name?.charAt(0).toUpperCase()}
      </div>
      <div className="flex-grow flex flex-col min-w-0 gap-1">
        <h4 className="font-medium leading-4 truncate">{name}</h4>
        <p className="text-sm text-dc-neutral-300 leading-3">
          {users.length} members
        </p>
      </div>
    </Link>
  );
}
