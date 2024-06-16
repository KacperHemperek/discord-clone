import { VariantProps, cva } from "class-variance-authority";
import { Link, useMatch } from "react-router-dom";

type LinkVariants = VariantProps<typeof link>;

type LinkProps = {
  to: string;
  label: string;
  showBadge?: boolean;
} & LinkVariants;

const link = cva(
  "px-2 capitalize font-medium rounded-md transition-colors relative",
  {
    variants: {
      variant: {
        neutral: "hover:text-dc-neutral-100 hover:bg-dc-neutral-800",
        success: undefined,
      },
      state: {
        active: undefined,
        inactive: undefined,
      },
    },
    compoundVariants: [
      {
        variant: "neutral",
        state: "active",
        className: "bg-dc-neutral-700 text-dc-neutral-50",
      },
      {
        variant: "success",
        state: "active",
        className: "text-dc-green-400",
      },
      {
        variant: "success",
        state: "inactive",
        className: "bg-dc-green-600 text-dc-purple-50 hover:bg-dc-green-500",
      },
    ],
    defaultVariants: {
      variant: "neutral",
      state: "inactive",
    },
  },
);

export function FriendsNavLink({ to, label, showBadge, variant }: LinkProps) {
  const match = useMatch(to);

  return (
    <Link
      to={to}
      className={link({
        variant,
        state: match ? "active" : "inactive",
      })}
    >
      {showBadge && (
        <div className="p-1.5 rounded-full bg-dc-red-500 absolute -top-[4px] -right-[4px]" />
      )}
      {label}
    </Link>
  );
}
