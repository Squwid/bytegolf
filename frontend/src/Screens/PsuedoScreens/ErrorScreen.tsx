import { makeStyles } from '@material-ui/core';
import React from 'react';
import Nav from '../../Components/Nav/Nav';
import { PrimaryColor } from '../../Globals';
import { NavType } from '../../Types';

type Props = {
  active?: NavType;
  text?: string;
}

const ErrorScreen: React.FC<Props> = (props) => {
  const classes = useStyles();
  const text = props.text ? (<p className={classes.notFoundText}>{props.text.toUpperCase()}</p>) :
    (<p className={classes.notFoundText}>AN UNEXPECTED ERROR OCCURRED</p>);
  const activeText = !!props.active ? (props.active) : 'none';

  return (
    <div className={classes.notFoundWrapper}>
      <Nav active={activeText} />
      <div className={classes.notFoundBodyWrapper}>
        <p style={{fontSize: '6rem', fontFamily: 'FiraCode', color: 'white', margin: 0, padding: 0}}>ERROR</p>
        {/* <img style={{margin:0, padding:0}} height='250px' width='500px' src={Logo} alt="Bytegolf not found logo" /> */}
        {text}
      </div>
    </div>
  )
}

const useStyles= makeStyles({
  notFoundWrapper: {
    height: '100vh',
    minHeight: '500px',
    backgroundColor: PrimaryColor,
  },
  notFoundBodyWrapper: {
    width: '800px',
    margin: '0 auto',
    marginTop: '100px',
    fontFamily: 'FiraCode',
    fontWeight: 'lighter',
    letterSpacing: '-.09rem',

    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center'
  },
  notFoundText: {
    color: 'white',
    fontSize: '1.6rem',
    margin: 0,
    padding: 0,
    textAlign: 'center'
  }
});

export default ErrorScreen;