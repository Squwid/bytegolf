import React from 'react';
import CSS from 'csstype';
import { makeStyles } from '@material-ui/core';

export type ChipProps = {
  ckey: string;
  value: string;
  secondaryTextColor: CSS.Property.Color;
  bgColor: CSS.Property.BackgroundColor;
  style?: CSS.Properties;
}

const Chip: React.FC<ChipProps> = (props) => {
  let chipStyles = makeStyles({
    chip: {
      fontFamily: 'FiraCode',
      border: `3px ${props.bgColor} solid`,
      borderRight: 'none',
      display: 'flex',
      flexDirection: 'row',
      justifyContent: 'flex-start',
      alignItems: 'center',

      height: '20px',
      padding: 0,
      margin: 0
    },
    chipText: {
      fontSize: '.9rem',
      padding: 0,
      margin: 0
    },
    chipTextWrapper: {
      paddingRight: '5px',
      paddingLeft: '5px',
      height: '100%',
      
      display: 'flex',
      alignItems: 'center'
    }
  })();

  
  return (
    <div className={chipStyles.chip} style={props.style}>
      <div className={chipStyles.chipTextWrapper} style={{backgroundColor: 'white'}}>
        <p className={chipStyles.chipText} >{props.ckey.toUpperCase()}</p>
      </div>

      <div className={chipStyles.chipTextWrapper} style={{color: props.secondaryTextColor,backgroundColor: props.bgColor}}>
        <p className={chipStyles.chipText} style={{fontWeight: 'bolder'}}>{props.value.toUpperCase()}</p>
      </div>
    </div>
  );
}

export default Chip;