import React from 'react';

type FriendListItemButtonProps = {
  icon: React.ReactNode;
  onClick: (e: React.MouseEvent<HTMLButtonElement>) => void;
  disabled?: boolean;
};

const FriendListItemButton = React.forwardRef<
  HTMLButtonElement,
  FriendListItemButtonProps
>(({ icon, onClick, disabled }, ref) => {
  return (
    <button
      ref={ref}
      onClick={onClick}
      className='p-2 rounded-full bg-dc-neutral-950 group disabled:opacity-50 disabled:cursor-not-allowed transition-opacity duration-100'
    >
      {icon}
    </button>
  );
});

export default FriendListItemButton;
