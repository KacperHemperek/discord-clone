import { createBrowserRouter } from 'react-router-dom';
import LoginPage from './pages/Login';
import RegisterPage from './pages/Register';
import BaseLayout from './layouts/BaseLayout';
import AuthGuard from './layouts/AuthGuard';
import FriendsLayout from './layouts/FriendsLayout';
import FriendRequestsPage from './pages/FriendRequests';
import AllFriendsPage from './pages/AllFriends';
import InviteUserPage from './pages/InviteUser';
import PrivateChat from './pages/PrivateChat';

export const router = createBrowserRouter([
  {
    path: '/',
    element: <AuthGuard />,
    children: [
      {
        path: 'home',
        element: <BaseLayout />,
        children: [
          {
            path: 'friends',
            element: <FriendsLayout />,
            children: [
              {
                path: '',
                element: <AllFriendsPage />,
              },
              {
                path: 'requests',
                element: <FriendRequestsPage />,
              },
              {
                path: 'invite',
                element: <InviteUserPage />,
              },
            ],
          },
          {
            path: 'chats/:chatId',
            element: <PrivateChat />,
          },
        ],
      },
      {
        path: 'login',
        element: <LoginPage />,
      },
      {
        path: 'register',
        element: <RegisterPage />,
      },
    ],
  },
]);
