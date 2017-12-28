import * as React from 'react';
import { Day, DAYS, Schedule } from './types';
import './ScheduleList.css';

interface ScheduleEntryProps {
  schedule: Schedule;
  onDelete(entry: Schedule): void;
}

function padTime(time: number): string {
  let str = time.toString();
  if (str.length < 2) {
    return '0' + str;
  }
  return str;
}

const renderTime = function(time: number) {
  let hour = Math.floor(time / 3600);
  let minutes = Math.floor((time - (hour * 3600)) / 60);
  let seconds = time - (hour * 3600) - (minutes * 60);
  return (
    <span className="TimeDisplay">{`${padTime(hour)}:${padTime(minutes)}:${padTime(seconds)}`}</span>
  );
};

const ScheduleEntry = function(props: ScheduleEntryProps) {
  return (
    <div className="ScheduleEntry">
      <div className="LeftContainer">
        <div className="ScheduleTime">
          {renderTime(props.schedule.time)}
        </div>
        <div className="ScheduleDays">
          {
            DAYS.map( (day) => {
              let className = props.schedule.days.indexOf(day as Day) === -1 ?
                'InactiveDay' : 'ActiveDay';

              return <span key={day} className={className}>{day.substring(0, 3)}</span>;
            })
          }
          <span className={`ScheduleCommand ${props.schedule.command}`}>
            | {props.schedule.command}
          </span>
        </div>
      </div>
      <div className="RightContainer">
        <div className="ScheduleDeleteButton" onClick={() => props.onDelete(props.schedule)}>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            xmlnsXlink="http://www.w3.org/1999/xlink"
            viewBox="0 0 742 742"
            version="1.1"
            height="40"
            width="30"
          >
            <path
              // tslint:disable-next-line:max-line-length
              d="M733 9q-9-9-22.5-9T688 9L371 326 54 9q-9-9-22.5-9T9 9 0 31.5 9 54l317 317L9 688q-9 9-9 22.5T9 733t22.5 9 22.5-9l317-317 317 317q9 9 22.5 9t22.5-9 9-22.5-9-22.5L416 371 733 54q9-9 9-22.5T733 9z"
            />
          </svg>
        </div>
      </div>
    </div>
  );
};

interface ScheduleListProps {
  schedules: Array<Schedule>;
  onDelete: (entry: Schedule) => void;
}

const ScheduleList = function(props: ScheduleListProps) {
  return (
    <section className="ScheduleList">
      {
        props.schedules.map((schedule) => {
          return (
            <ScheduleEntry key={schedule.id} schedule={schedule} onDelete={props.onDelete}/>
          );
        })
      }
    </section>
  );
};

export default ScheduleList;
