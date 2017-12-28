import * as React from 'react';

interface ToggleIconProps {
  color: string;
  blur: boolean;
  height: number;
  width: number;
  onClick: (e: {}) => void;
}
const ToggleIcon = function(props: ToggleIconProps) {
  return (
    <svg
      className="toggle-button"
      viewBox="0 0 1024 1024"
      xmlns="http://www.w3.org/2000/svg"
      xmlnsXlink="http://www.w3.org/1999/xlink"
      version="1.1"
      overflow="visible"
      height={props.height}
      width={props.width}
      onClick={props.onClick}
    >
      <filter id="blur">
        <feGaussianBlur in="SourceGraphic" stdDeviation="90"/>
      </filter>
      <g fill={props.color}>
        <title id="fhsi-ant-poweroff-title">icon poweroff</title>
        <path
          // tslint:disable-next-line:max-line-length
          d="M512 384q-13 0-22.5-9.5T480 352V32q0-13 9.5-22.5T512 0t22.5 9.5T544 32v320q0 13-9.5 22.5T512 384zM638 55v6q0 20 20 27 14 5 28 11 81 34 143 96t96 143q35 83 35 174t-35 174q-34 81-96 143t-143 96q-83 35-174 35t-174-35q-81-34-143-96T99 686q-35-83-35-174t35-174q34-81 96-143t143-96q14-6 26-10 22-8 22-30 0-17-13.5-26.5T343 28Q191 81 95.5 213.5T0 512q0 139 68.5 257T255 955.5t257 68.5 257-68.5T955.5 769t68.5-257q0-167-97-300T677 27q-15-5-27 4t-12 24z"
        />
      </g>
      <g fill={props.color} filter={props.blur ? 'url(#blur)' : ''}>
        <title id="fhsi-ant-poweroff-title">icon poweroff</title>
        <path
          // tslint:disable-next-line:max-line-length
          d="M512 384q-13 0-22.5-9.5T480 352V32q0-13 9.5-22.5T512 0t22.5 9.5T544 32v320q0 13-9.5 22.5T512 384zM638 55v6q0 20 20 27 14 5 28 11 81 34 143 96t96 143q35 83 35 174t-35 174q-34 81-96 143t-143 96q-83 35-174 35t-174-35q-81-34-143-96T99 686q-35-83-35-174t35-174q34-81 96-143t143-96q14-6 26-10 22-8 22-30 0-17-13.5-26.5T343 28Q191 81 95.5 213.5T0 512q0 139 68.5 257T255 955.5t257 68.5 257-68.5T955.5 769t68.5-257q0-167-97-300T677 27q-15-5-27 4t-12 24z"
        />
      </g>
    </svg>
  );
};

interface ToggleButtonProps {
  active: boolean;
  onClick(e: {}): void;
}
const ToggleButton = function(props: ToggleButtonProps) {
  let color = props.active ? '#EEE' : '#6D6D6D';

  return (
    <div className="App-header">
      <ToggleIcon color={color} height={80} width={80} blur={props.active} onClick={props.onClick}/>
    </div>
  );
};

export default ToggleButton;