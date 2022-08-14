import React, { useEffect } from 'react';
import MySubmissions from '../Components/MySubmissions/MySubmissions';
import Nav from '../Components/Nav/Nav';
import './screen.css';

const RecentScreen: React.FC = () => {
  useEffect(() => {document.title=`Bytegolf - Recent Submissions`}, []);

  return (
    <div>
      <Nav active='recent'/>
      <p className='screenTitle'>RECENT SUBMISSIONS</p>
      <div className='screenContainer'>
        <MySubmissions />
      </div>
    </div>
  )
}

export default RecentScreen;