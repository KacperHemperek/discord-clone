export const DAY_IN_MS = 1000 * 60 * 60 * 24;

export const MINUTE_IN_MS = 1000 * 60;

export function getDayAtMidnight(inputDate: Date | string | number) {
  const date = new Date(inputDate);
  date.setHours(0, 0, 0, 0);

  return date;
}

export function formatShortDate(inputDate: Date | string | number) {
  const date = new Date(inputDate);

  return Intl.DateTimeFormat('en-UK', {
    day: 'numeric',
    month: 'numeric',
    year: 'numeric',
  }).format(date);
}

export function formatMediumDate(inputDate: Date | string | number) {
  const date = new Date(inputDate);

  return Intl.DateTimeFormat('en-UK', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  }).format(date);
}

export function formatShortTime(inputDate: Date | string | number) {
  const date = new Date(inputDate);

  return Intl.DateTimeFormat('en-UK', {
    hour: 'numeric',
    minute: 'numeric',
  }).format(date);
}
