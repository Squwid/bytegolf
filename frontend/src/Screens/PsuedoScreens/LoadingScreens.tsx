import { makeStyles } from '@material-ui/core';
import React from 'react';
import Nav from '../../Components/Nav/Nav';
import { PrimaryColor } from '../../Globals';
import LoadingIcon from '../../Logo/LoadingIcon/LoadingIcon';
import { NavType } from '../../Types';

type Props = {
  active?: NavType;
}

const LoadingScreen: React.FC<Props> = (props) => {
  const classes = useStyles();

  return (
    <div className={classes.loadingScreenPageWrapper}>
      <Nav active={props.active ? props.active : 'none'} />

      <div style={{margin: '0 auto', display: 'flex', width: '200px', flexDirection: 'column', alignItems: 'center'}}>
        <LoadingIcon style={{width: '200px', height: '200px'}}/>
        <p style={{fontSize: '1.5rem', fontFamily: 'FiraCode', color: PrimaryColor, padding: 0, marginTop: -20}}>LOADING...</p>
      </div>
    </div>
  )
}

const useStyles= makeStyles({
  loadingScreenPageWrapper: {
    height: '100vh',
    minHeight: '500px',
    backgroundColor: '#f5f5f5',
  },
});

export default LoadingScreen;