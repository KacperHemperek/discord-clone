import * as P from "@radix-ui/react-popover";
import { cn } from "@app/utils/cn.ts";
import DCButton from "@app/components/Button.tsx";
import React from "react";
import { UserResponse } from "@app/api";
import { Checkbox, CheckboxIndicator } from "@radix-ui/react-checkbox";
import { CheckIcon } from "lucide-react";

type UserCheckboxItemProps = {
  id: number;
  email: string;
  username: string;
  selected: boolean;
  onCheck: (val: boolean) => void;
  avatar?: string;
};

function UserCheckboxItem({
  username,
  onCheck,
  selected,
}: UserCheckboxItemProps) {
  return (
    <label className="flex justify-between items-center cursor-pointer py-1 hover:bg-dc-neutral-850 px-2 -mx-2 rounded-sm">
      <div className="flex items-center gap-2">
        {/* TODO: replace with avatar of user in future */}
        <div className="w-8 h-8 rounded-full bg-dc-neutral-800" />
        <p className="truncate max-w-[200px]">{username}</p>
      </div>
      <Checkbox
        onCheckedChange={onCheck}
        checked={selected}
        className='p-1 border-[1.5px] border-dc-neutral-600 rounded-md data-[state="checked"]:border-dc-purple-500 data-[state="checked"]:text-dc-purple-500'
      >
        <div className="w-4 h-4 flex">
          <CheckboxIndicator asChild>
            <CheckIcon className="w-full h-full" strokeWidth={2.5} />
          </CheckboxIndicator>
        </div>
      </Checkbox>
    </label>
  );
}

export const UserSelectListTrigger = P.Trigger;

export const UserSelectListRoot = P.Root;

type UserSelectListProps = {
  users: UserResponse[];
  onSubmit: (selectedIds: number[]) => void;
  submitLabel: string;
  selectedIds: number[];
  setSelectedIds: React.Dispatch<React.SetStateAction<number[]>>;
};

export function UserSelectListContent({
  users,
  onSubmit,
  submitLabel,
  setSelectedIds,
  selectedIds,
}: UserSelectListProps) {
  const lastScrollPosition = React.useRef(0);
  const [showBottomBorder, setShowBottomBorder] = React.useState(false);

  React.useEffect(() => {
    return () => {
      setSelectedIds([]);
    };
  }, [setSelectedIds]);

  function handleScroll(e: React.UIEvent<HTMLDivElement, UIEvent>) {
    const target = e.currentTarget;

    if (target.scrollTop > lastScrollPosition.current) {
      setShowBottomBorder(true);
    }

    if (target.scrollTop === 0) {
      setShowBottomBorder(false);
    }

    lastScrollPosition.current = target.scrollTop;
  }

  function toggleUserSelection(val: boolean, userId: number) {
    if (val) {
      setSelectedIds((prev) => [...prev, userId]);
    } else {
      setSelectedIds((prev) => prev.filter((id) => id !== userId));
    }
  }

  return (
    <P.Portal>
      <P.Content
        className="bg-dc-neutral-900 rounded-md border border-dc-neutral-1000 w-[420px] flex flex-col text-dc-neutral-50 shadow-md"
        side="bottom"
        sideOffset={8}
        align="start"
      >
        <div
          className={cn(
            "flex flex-col px-4 pt-4 border-b transition-colors",
            showBottomBorder
              ? "border-b-dc-neutral-1000"
              : "border-b-transparent",
          )}
        >
          <h2 className="text-lg font-semibold text-dc-neutral-50">
            Choose friends
          </h2>
          <p className="text-xs font-medium text-dc-neutral-300 pb-4">
            You can add up to {users.length} friends
          </p>
        </div>
        <div
          className="flex flex-col px-4 overflow-y-auto max-h-36"
          onScroll={handleScroll}
        >
          <div className="grid pb-2 gap-1">
            {users.map((friend, i) => (
              <UserCheckboxItem
                {...friend}
                key={i}
                onCheck={(val) => toggleUserSelection(val, friend.id)}
                selected={selectedIds.includes(friend.id)}
              />
            ))}
          </div>
        </div>
        <div className="p-4 relative">
          <span className="absolute top-0 left-0 flex w-full px-2">
            <span className="h-[1px] flex-grow bg-dc-neutral-850"></span>
          </span>
          <DCButton onClick={() => onSubmit(selectedIds)} className="w-full">
            {submitLabel}
          </DCButton>
        </div>
      </P.Content>
    </P.Portal>
  );
}
