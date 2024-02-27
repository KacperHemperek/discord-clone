import React from 'react';
import { Navigate, Outlet, matchPath, useLocation } from 'react-router-dom';
import { useAuth } from '../../context/AuthProvider';

function UserLoading() {
  return (
    <div className='w-screen h-screen flex items-center justify-center bg-dc-neutral-800 text-center gap-2 flex-col'>
      <h1 className='text-xl font-semibold text-dc-neutral-50'>Loading...</h1>
      <p className='text-dc-neutral-300 max-w-[300px]'>
        We are checking if you are logged in give us a second
      </p>
    </div>
  );
}

export default function AuthGuard() {
  const { pathname } = useLocation();
  const { user, isLoadingUser } = useAuth();
  const filterMatch = (route: string) => {
    const match = matchPath(route, pathname);

    return !!match;
  };

  const protectedRoutes = ['/home/*'];

  const redirectUserIfLoggedIn = ['/login', '/register'];
  const matchesProtected = protectedRoutes.some(filterMatch);

  const matchesRedirectUserIfLoggedIn =
    redirectUserIfLoggedIn.some(filterMatch);

  const isRoot = pathname === '/';

  if (isLoadingUser) {
    return <UserLoading />;
  }

  if (isRoot && user) {
    return <Navigate to='/home/friends' />;
  }

  if (isRoot && !user) {
    return <Navigate to='/login' />;
  }

  if (user && matchesRedirectUserIfLoggedIn) {
    return <Navigate to='/home/friends' />;
  }

  if (!user && matchesProtected) {
    return <Navigate to='/login' />;
  }

  return <Outlet />;
}
