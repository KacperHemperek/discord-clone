import {
  DAY_IN_MS,
  formatShortDate,
  formatShortTime,
  getDayAtMidnight,
} from '../utils/dates';

export function useChatMessageDateFormatter(inputDate: Date | string | number) {
  const date = new Date(inputDate);

  const dateToFormat = new Date(inputDate);
  const time = formatShortTime(dateToFormat);

  const today = getDayAtMidnight(new Date());
  const day = getDayAtMidnight(date);

  if (day.getTime() === today.getTime()) {
    return `Today at ${time}`;
  }

  if (day.getTime() === today.getTime() - DAY_IN_MS) {
    return `Yesterday at ${time}`;
  }

  const formattedDate = formatShortDate(dateToFormat);

  return `${formattedDate} ${time}`;
}
