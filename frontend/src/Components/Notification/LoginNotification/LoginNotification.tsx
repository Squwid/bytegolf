import React from 'react';
import { useHistory } from 'react-router';
import Notification from '../Notification';

export const LoginNotification: React.FC = () => {
  const history = useHistory();
  return (<Notification 
    type='warn'
    text='NOT LOGGED IN, CLICK HERE TO LOGIN AND SEE PAST SUBMISSIONS'
    style={{marginBottom: '1rem'}}
    onClick={() => history.push('/profile')}
  />);
}