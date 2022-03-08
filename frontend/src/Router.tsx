import { BrowserRouter, Route, Switch } from 'react-router-dom';
import HolesScreen from './Screens/HolesScreen';
import HomeScreen from './Screens/HomeScreen';
import LeaderboardScreen from './Screens/LeaderboardScreen';
import NotFoundScreen from './Screens/PsuedoScreens/NotFoundScreen';
import PlayScreen from './Screens/PlayScreen';
import ProfileRedirectScreen from './Screens/PsuedoScreens/ProfileRedirectScreen';
import ProfileScreen from './Screens/ProfileScreen';
import RecentScreen from './Screens/RecentScreen';

const Router = () => {
  return (
    <BrowserRouter>
      <Switch>
        <Route exact path={'/'} component={HomeScreen} />
        <Route exact path={'/play'} component={HolesScreen} />
        <Route exact path={'/play/:holeID'} component={PlayScreen} />
        <Route exact path={'/leaderboard'} component={LeaderboardScreen} />
        <Route exact path={'/recent'} component={RecentScreen} />
        <Route exact path={'/profile'} component={ProfileRedirectScreen} />
        <Route exact path={'/profile/:bgid'} component={ProfileScreen} />

        <Route component={NotFoundScreen} />
      </Switch>
    </BrowserRouter>
  )
}

export default Router;