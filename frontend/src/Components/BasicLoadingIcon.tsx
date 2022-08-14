import React from 'react';
import { PrimaryColor } from '../Globals';
import LoadingIcon from '../Logo/LoadingIcon/LoadingIcon';

export const BasicLoadingIcon = () => (
  <div style={{margin: '0 auto', display: 'flex', width: '200px', flexDirection: 'column', alignItems: 'center'}}>
    <LoadingIcon style={{width: '200px', height: '200px'}}/>
    <p style={{fontSize: '1.5rem', fontFamily: 'FiraCode', color: PrimaryColor, padding: 0, marginTop: -20}}>LOADING...</p>
  </div>
);