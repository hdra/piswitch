export const DAYS: Array<Day> = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];

export type Day = 'Sunday' | 'Monday' | 'Tuesday' | 'Wednesday' | 'Thursday' | 'Friday' | 'Saturday';

export const COMMANDS: Array<ScheduleCommand> = ['On', 'Off', 'Toggle'];

export type ScheduleCommand = 'On' | 'Off' | 'Toggle';

export interface Schedule {
  id: string;
  time: number; // Integer between 0 and 86400
  days: Array<Day>;
  command: ScheduleCommand;
}
