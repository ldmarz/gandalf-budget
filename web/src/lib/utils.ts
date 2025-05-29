import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export const formatCurrency = (amount: number | null | undefined) => {
  if (amount === null || amount === undefined) {
    return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(0);
  }
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(amount);
};
