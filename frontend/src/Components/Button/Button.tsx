import React from 'react';
import CSS from 'csstype';
import './Button.css';
import { CircularProgress, makeStyles } from '@material-ui/core';
import { PrimaryColor, ThirdColor } from '../../Globals';

type Props = {
  to?: string;
  fontSize?: CSS.Property.FontSize;
  color: CSS.Property.Color;
  activeColor?: CSS.Property.Color;

  disabled?: boolean;
  loading?: boolean;

  text: string;
  margins?: CSS.Property.Margin;
  active?: boolean;
  onPress?: () => void;
}

const Button: React.FC<Props> = (props) => {
  const margin = props.margins ? props.margins : 5;
  const activeText = props.activeColor ? props.activeColor : ThirdColor;
  const hoverColor = !!props.active ? props.color : 'lightgreen';

  const isLoading = !!props.loading;

  let button = makeStyles({
    btn: {
      fontFamily: 'FiraCode',
      fontSize: props.fontSize ? props.fontSize : '1rem',
      padding: '0.5rem',
      border: `3px ${props.color} solid`,
      backgroundColor: props.active ? PrimaryColor : 'white',
      cursor: 'pointer',
      color: ThirdColor,
      letterSpacing: '-.08rem',
      margin: margin,
      '&:hover': {
        backgroundColor: hoverColor
      },
      '&:active': {
        backgroundColor: props.color,
        color: activeText
      }
    },

    btnDisabled: {
      fontFamily: 'FiraCode',
      fontSize: props.fontSize ? props.fontSize : '1rem',
      padding: '0.5rem',
      border: `3px ${props.color} solid`,
      backgroundColor: 'lightgray',
      color: ThirdColor,
      cursor: 'default',
      letterSpacing: '-.08rem',
      margin: margin,
    },
  })

  const classes = button();

  return (
    <div className={props.disabled ? classes.btnDisabled : classes.btn} onClick={props?.onPress}>
      {props.text.toUpperCase()}
      {isLoading && <CircularProgress style={{marginLeft: '1rem'}} color='primary' size={props.fontSize?props.fontSize:'1rem'}/>}
    </div>
  );
}



export default Button;