import { useMutation, useQueryClient } from "@tanstack/react-query";
import { MutationKeys } from "@app/api/mutationKeys.ts";
import { api, QueryKeys } from "@app/api";

export function useMarkFriendRequestNotificationsAsSeen() {
  const qc = useQueryClient();

  return useMutation({
    mutationKey: MutationKeys.markFriendRequestNotificationAsSeen(),
    mutationFn: async () =>
      await api.put("/notifications/friend-requests/mark-as-seen"),
    onSuccess: () => {
      void qc.invalidateQueries({
        queryKey: QueryKeys.getFriendRequestNotifications(),
      });
    },
  });
}
