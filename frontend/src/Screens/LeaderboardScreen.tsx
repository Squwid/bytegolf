import React, { useEffect } from 'react';
import Button from '../Components/Button/Button';
import Nav from '../Components/Nav/Nav';
import { PrimaryColor } from '../Globals';
import Leaderboard from '../Components/Leaderboard/Leaderboard';
import './screen.css';
import { BasicHole } from '../Types';
import MySubmissions from '../Components/MySubmissions/MySubmissions';
import { ListHoles } from '../Store/Holes';
import { useQuery } from 'react-query';
import LoadingScreen from './PsuedoScreens/LoadingScreens';
import ErrorScreen from './PsuedoScreens/ErrorScreen';
import { useHistory, useLocation } from 'react-router-dom';
 

const LeaderboardScreen: React.FC = () => {
  // Use query params rather than 
  const loc = useLocation();
  const history = useHistory();

  const params = new URLSearchParams(loc.search);
  const paramHole = params.get('hole');
  const selectedHole: string = paramHole ? paramHole : '';

  useEffect(()=>{document.title = 'Bytegolf - Leaderboards'},[]);

  // TODO: Get holes & Leaders for each hole
  const holes = useQuery('Holes', () =>  ListHoles());
  if (holes.isLoading) return (<LoadingScreen active='leaderboards' />);
  if (holes.isError || !holes.data) return (<ErrorScreen active='leaderboards' text={`${holes.error}`}/>);

  const onHoleClick = (hole: BasicHole) => {
    if (hole.ID === selectedHole) return;
    
    const newParams = new URLSearchParams(loc.search);
    newParams.set('hole', hole.ID);
    history.push(`${loc.pathname}?${newParams.toString()}`)
  }

  const isHoleActive = (hole: BasicHole|null): boolean => {
    if (!hole) return false;
    if (hole.ID === selectedHole) return true;
    return false
  }

  const leaderboard = (): JSX.Element => {
    if (!selectedHole) return (<></>);

    return (
      <div className='screenContainer'>
        <p className='screenSubText'>LEADERS FOR {selectedHole.toUpperCase()}</p>
        <Leaderboard holeID={selectedHole} />
      </div>
    )
  }

  const submissions = (): JSX.Element => {
    if (!selectedHole) return (<></>);

    return (
      <div className='screenContainer'>
        <p className='screenSubText'>MY SUBMISSIONS</p>
        <MySubmissions hole={selectedHole} />
      </div>
    )
  }

  return (
    <div>
      <Nav active={'leaderboards'} />
      <p className='screenTitle'>LEADERBOARDS</p>
      <div className='holesBtnContainer'>
        {holes.data.map(hole => <Button onPress={() => onHoleClick(hole)} active={isHoleActive(hole)} color={PrimaryColor} activeColor='white' fontSize='1rem' text={hole.Name} />)}
      </div>
      
      {leaderboard()}
      {submissions()}
    </div>
  );
}

export default LeaderboardScreen;