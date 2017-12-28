import * as React from 'react';
import * as ReactModal from 'react-modal';
import Spinner from './Spinner';
import ToggleButton from './ToggleButton';
import ScheduleList from './ScheduleList';
import NewSchedule from './NewSchedule';
import { Schedule } from './types';
import './App.css';

interface AppState {
  showModal: boolean;
  loading: boolean;
  currentState: boolean;
  schedules: Array<Schedule>;
}

class App extends React.Component <{}, AppState> {
  constructor(props: {}) {
    super(props);
    this.state = {
      showModal: false,
      loading: true,
      currentState: false,
      schedules: []
    };
  }
  componentDidMount() {
    fetch('/api').then((response) => {
      response.json().then((data) => {
        this.setState({
          ...this.state,
          loading: false,
          currentState: data.currentState,
          schedules: data.schedules
        });
      });
    });
  }
  onScheduleDelete = (schedule: Schedule) => {
    this.setState({
      ...this.state,
      loading: true
    });

    fetch('/api/remove', {
      method: 'POST',
      body: JSON.stringify({id: schedule.id})
    }).then((response) => {
      if (response.ok) {
        this.setState({
          ...this.state,
          loading: false,
          schedules: this.state.schedules.filter((entry) => {
            return entry.id !== schedule.id;
          })
        });
      } else {
        response.text().then( (message) => {
          alert(message);
          this.setState({
            loading: false
          });
        });
      }
    }).catch((err) => {
      alert('An error has occurred: ' + err);
      this.setState({
        loading: false
      });
    });
  }
  onScheduleAdd = (newSchedule: Schedule) => {
    this.setState({
      ...this.state,
      loading: true
    });

    fetch('/api/add', {
      method: 'POST',
      body: JSON.stringify(newSchedule)
    }).then((response) => {
      if (response.ok) {
        this.setState({
          ...this.state,
          loading: false,
          schedules: [...this.state.schedules, newSchedule]
        });
        this.closeModal();
      } else {
        response.text().then( (message) => {
          alert(message);
          this.setState({
            loading: false
          });
        });
      }
    }).catch((err) => {
      alert('An error has occurred: ' + err);
      this.setState({
        loading: false
      });
    });
  }
  onTogglePress = () => {
    this.setState({
      ...this.state,
      loading: true
    });

    fetch('/api/toggle', {
      method: 'POST'
    }).then((response) => {
      if (response.ok) {
        this.setState({
          ...this.state,
          loading: false,
          currentState: !this.state.currentState
        });
      } else {
        response.text().then( (message) => {
          alert(message);
          this.setState({
            loading: false
          });
        });
      }
    }).catch((err) => {
      alert('An error has occurred: ' + err);
      this.setState({
        loading: false
      });
    });
  }
  openModal = () => {
    this.setState({
      ...this.state,
      showModal: true
    });
  }
  closeModal = () => {
    this.setState({
      ...this.state,
      showModal: false
    });
  }
  render() {
    return (
      <div className="App">
        <ToggleButton active={this.state.currentState} onClick={this.onTogglePress} />
        <ScheduleList schedules={this.state.schedules} onDelete={this.onScheduleDelete} />
        <button className="NewButton" onClick={this.openModal}>NEW</button>
        <ReactModal
          isOpen={this.state.showModal}
          onRequestClose={this.closeModal}
          className="Modal"
        >
          <NewSchedule onSave={this.onScheduleAdd}/>
        </ReactModal>
        <Spinner display={this.state.loading} />
      </div>
    );
  }
}

export default App;
