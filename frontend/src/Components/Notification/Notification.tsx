import React from 'react';
import CSS from 'csstype';
import { makeStyles } from '@material-ui/core';
import { PrimaryColor, SecondaryColor } from '../../Globals';

export type NotificationProps = {
  type: 'info'|'error'|'warn';
  text: string;
  // timeout?: number;
  width?: CSS.Property.Width;
  height?: CSS.Property.Height;
  style?: CSS.Properties;
  onClick?: ()=>void;
}



const Notification: React.FC<NotificationProps> = (props) => {
  const classes = makeStyles({
    info: {
      color: PrimaryColor,
      backgroundColor: '#BCFFC3',
      border: `1px ${PrimaryColor} solid`
    },
    error: {
      color: SecondaryColor,
      backgroundColor: '#FFCFC4',
      border: `1px ${SecondaryColor} solid`
    },
    warn: {
      color: '#ab900c',
      backgroundColor: '#f0e999',
      border: '1px #ab900c solid'
    }
  })();

  return(
    <div className={classes[props.type]} style={{
      height: !!props.height ? props.height : 'auto',
      width: !!props.width ? props.width : 'auto',
      fontFamily: 'FiraCode',
      paddingLeft: '10px',
      paddingRight: '10px',
      borderRadius: '10px',
      cursor: 'pointer',
      ...props.style
    }} onClick={props?.onClick}>
      <p>{props.text.toUpperCase()}</p>
    </div>

  )
}


export default Notification;