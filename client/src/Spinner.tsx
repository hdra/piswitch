import * as React from 'react';
import './Spinner.css';

interface SpinnerProps {
  display: boolean;
}
const Spinner = function(props: SpinnerProps) {
  let activeClass = props.display ? 'active' : 'inactive';
  return (
    <div className={`Spinner ${activeClass}`}>
      <div className="sk-folding-cube">
        <div className="sk-cube1 sk-cube"/>
        <div className="sk-cube2 sk-cube"/>
        <div className="sk-cube4 sk-cube"/>
        <div className="sk-cube3 sk-cube"/>
      </div>
    </div>
  );
};

export default Spinner;
