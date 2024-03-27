export class Time {
  static second = 1000;

  static minute = 60 * Time.second;

  static hour = 60 * Time.minute;

  static day = Time.hour * 24;
}

export function getDayAtMidnight(inputDate: Date | string | number) {
  const date = new Date(inputDate);
  date.setHours(0, 0, 0, 0);

  return date;
}

export function formatShortDate(inputDate: Date | string | number) {
  const date = new Date(inputDate);

  return Intl.DateTimeFormat("en-UK", {
    day: "numeric",
    month: "numeric",
    year: "numeric",
  }).format(date);
}

export function formatMediumDate(inputDate: Date | string | number) {
  const date = new Date(inputDate);

  return Intl.DateTimeFormat("en-UK", {
    day: "numeric",
    month: "long",
    year: "numeric",
  }).format(date);
}

export function formatShortTime(inputDate: Date | string | number) {
  const date = new Date(inputDate);

  return Intl.DateTimeFormat("en-UK", {
    hour: "numeric",
    minute: "numeric",
  }).format(date);
}
