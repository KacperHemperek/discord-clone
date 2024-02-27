export enum FriendRequestStatus {
	pending = 'pending',
	accepted = 'accepted',
	declined = 'declined',
}

export enum FriendRequestType {
	newFriendInvite = 'NEW_FRIEND_INVITE',
	allFriendInvites = 'ALL_FRIEND_INVITES',
	friendInviteAccepted = 'FRIEND_INVITE_ACCEPTED',
	friendInviteDeclined = 'FRIEND_INVITE_DECLINED',
}

export const NEW_FRIEND_INVITE_TYPE = 'NEW_FRIEND_INVITE';

export const ALL_FRIEND_INVITES_TYPE = 'ALL_FRIEND_INVITES';
