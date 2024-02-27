import { formatMediumDate } from '../utils/dates';

export function useDateSeparatorFormatter(date: Date | string | number) {
  return formatMediumDate(new Date(date));
}
