import * as React from 'react';
import { DAYS, COMMANDS, Day, ScheduleCommand, Schedule } from './types';
import './NewSchedule.css';

interface NewScheduleState {
  time: string;
  days: Array<Day>;
  command: ScheduleCommand;
}

interface NewScheduleProps {
  onSave: (schedule: Schedule) => void;
}

export default class NewSchedule extends React.Component<NewScheduleProps, NewScheduleState> {
  state = {
    time: '',
    days: [] as Array<Day>,
    command: 'On' as ScheduleCommand
  };

  onTimeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    let time = e.target.value;
    this.setState({
      ...this.state,
      time: time
    });
  }

  onDaysChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    let day = e.target.value as Day;
    if (this.state.days.indexOf(day) >= 0) {
      this.setState({
        ...this.state,
        days: this.state.days.filter( (d) => d !== day )
      });
    } else {
      this.setState({
        ...this.state,
        days: [...this.state.days, day]
      });
    }
  }

  onCommandClick = (command: ScheduleCommand) => {
    this.setState({
      ...this.state,
      command: command
    });
  }

  isValid = () => {
    if (this.state.time !== '' && this.state.days.length > 0) {
      return true;
    }
    return false;
  }

  save = () => {
    if (!this.isValid()) {
      return;
    }
    let parts = this.state.time.split(':');
    let hour = parseInt(parts[0], 10);
    let minutes = parseInt(parts[1], 10);

    let hash = this.state.time + this.state.command + this.state.days.join('');

    let schedule: Schedule = {
      ...this.state,
      time: hour * 3600 + minutes * 60,
      id: hash
    };
    this.props.onSave(schedule);
  }

  render() {
    return (
      <div className="NewSchedule">
        <input type="time" className="TimeInput" onChange={this.onTimeChange} />
        <div className="DaysInput">
          {
            DAYS.map( (day) => {
              return (
                <label key={day}>
                  <input
                    type="checkbox"
                    value={day}
                    checked={this.state.days.indexOf(day) >= 0}
                    onChange={this.onDaysChange}
                  />
                  {day.substr(0, 3)}
                </label>
              );
            })
          }
        </div>
        <div className="CommandInput">
          {
            COMMANDS.map( (command) => {
              return (
                <button
                  key={command}
                  className={`Button ${command} ${this.state.command === command ? 'active' : ''}`}
                  onClick={() => this.onCommandClick(command)}
                >
                  {command}
                </button>
              );
            })
          }
        </div>
        <div className="FormButton">
          <button
            onClick={() => this.save()}
            className={`Button Expanded ${this.isValid() ? 'active' : ''}`}
          >
            SAVE
          </button>
        </div>
      </div>
    );
  }
}
