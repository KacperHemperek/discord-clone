import { HomeIcon, LogOut } from "lucide-react";
import { Link, Outlet, useMatch } from "react-router-dom";
import { useLogout } from "@app/api";
import { cn } from "../../utils/cn";
import ChatLinkList from "@app/components/chats/ChatLinkList.tsx";

function HomeLink() {
  const match = useMatch("/home/*");

  return (
    <Link
      to={"/home/friends"}
      className={cn(
        "p-3.5 rounded-full bg-dc-neutral-500 flex items-center justify-center transition",
        !!match && "bg-dc-neutral-600 rounded-2xl",
      )}
    >
      <HomeIcon className="w-5 h-5" />
    </Link>
  );
}

export default function HomeLayout() {
  const { mutate } = useLogout();

  function logout() {
    mutate();
  }

  return (
    <div className="flex h-screen max-h-screen bg-dc-neutral-900 overflow-hidden text-dc-neutral-50">
      <div className="flex flex-col bg-dc-neutral-1000 ">
        {/* Channels List */}
        <div className="flex-grow overflow-auto p-3 gap-2 flex flex-col">
          <HomeLink />
          <span className=" h-0.5 bg-dc-neutral-800 mx-2 rounded-full" />
          <button
            onClick={logout}
            className="p-3.5 rounded-full bg-dc-neutral-500 flex items-center justify-center"
          >
            <LogOut className="w-5 h-5" />
          </button>
        </div>
        {/* Settings */}
        <div className="flex flex-col gap-2 px-2 pb-2"></div>
      </div>
      <ChatLinkList />
      <Outlet />
    </div>
  );
}
