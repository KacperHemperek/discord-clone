import React from "react";

type RandomVariants<T extends object> = Record<number, T>;

export function useRandomVariant<T extends object>(
  variants: RandomVariants<T>,
) {
  const numberKeys = Object.keys(variants).map(Number);
  const min = Math.min(...numberKeys);
  const max = Math.max(...numberKeys);

  const variant = React.useRef(randomInt(min, max) as keyof typeof variants);

  function randomInt(min: number, max: number): number {
    return Math.floor(Math.random() * (max - min + 1)) + min;
  }

  return variants[variant.current] ?? variants[min];
}
