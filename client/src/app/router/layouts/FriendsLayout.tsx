import { Outlet } from "react-router-dom";
import { User } from "lucide-react";
import { useNotifications } from "../../context/NotificationsProvider.tsx";
import { FriendsNavLink } from "../../components/friends/FriendsNavLink";

export default function FriendsLayout() {
  const { hasUnseenFriendRequestNotifications: hasNewRequests } =
    useNotifications();

  return (
    <>
      <div className="flex-grow flex flex-col">
        <nav className="border-b flex border-dc-neutral-1000 w-full p-3 gap-4">
          <div className="flex gap-2 font-semibold items-center">
            <User />
            Friends
          </div>

          <div className="border-r border-dc-neutral-600" />

          <div className="flex gap-4">
            <FriendsNavLink to="/home/friends" label="All" />
            <FriendsNavLink
              to="/home/friends/requests"
              label="Requests"
              showBadge={hasNewRequests}
            />
            <FriendsNavLink
              to="/home/friends/invite"
              label="Invite Friend"
              variant="success"
            />
          </div>
        </nav>
        <Outlet />
      </div>
    </>
  );
}
