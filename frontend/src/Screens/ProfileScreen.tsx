import { makeStyles } from '@material-ui/core';
import React from 'react';
import { useQuery } from 'react-query';
import { useParams } from 'react-router';
import { useLocation } from 'react-router-dom';
import Nav from '../Components/Nav/Nav';
import { PrimaryColor, ThirdColor } from '../Globals';
import { GetProfile } from '../Store/Profile';
import ErrorScreen from './PsuedoScreens/ErrorScreen';
import LoadingScreen from './PsuedoScreens/LoadingScreens';
import NotFoundScreen from './PsuedoScreens/NotFoundScreen';

const ProfileScreen: React.FC = () => {
  const classes = useStyles();
  const { bgid } = useParams<{bgid: string}>();
  const { pathname } = useLocation();

  React.useEffect(() => {
    window.scrollTo(0, 0);
    document.title = `Bytegolf - Profile`;
  }, [pathname]);

  const profile = useQuery(['Profile', bgid], () => GetProfile(bgid));
  if (profile.isLoading) return (<LoadingScreen active='profile' />);
  if (profile.isError) return (<ErrorScreen active='profile' text={`${profile.error}`} />);
  if (!profile.data) return (<NotFoundScreen active='profile' text={`Profile ${bgid} was not found!`} />);
  
  return (
    <div>
      <Nav active={'profile'}/>
      <div className={classes.profileTopHeader}>
        <div className={classes.profileBody}>
          <div className={classes.profileDetailsContainer}>
            <img className={classes.profileImage} height="200px" width="200px" src={`https://avatars.githubusercontent.com/u/${profile.data.BGID}?v=4`} alt={`${profile.data.GithubUser.ID} github`} />

            <div className={classes.profileTextWrapper}>
              <p style={{fontWeight: 'bolder', fontSize: '2rem', padding: '0', margin: '0'}}>{profile.data.GithubUser.Login.toUpperCase()}</p>
              <p className={classes.profileGithub} onClick={() => window.open(profile.data?.GithubUser.URL, "_blank")}>{profile.data?.GithubUser.URL}</p>
            </div>
          </div>

          <div>
            <p style={{color: ThirdColor, fontSize: '2rem', marginTop: '100px'}}>MORE PROFILES COMING SOON!</p>
          </div>
        </div>
      </div>
    </div>
  );
}

const useStyles = makeStyles({
  profileBody: {
    textAlign: 'center',
    width: '980px',
    margin: '0 auto',
    height: 'auto',
    fontFamily: 'FiraCode',
    color: 'white'
    // backgroundColor: 'darkgray'
  },
  profileDetailsContainer: {
    height: '200px',
    // backgroundColor: 'lightblue',
    width: '100%',

    display: 'flex',
    flexDirection: 'row',

  },
  profileTopHeader: {
    backgroundColor: PrimaryColor,
    height: '160px'
  },
  profileImage: {
    borderRadius: '50%',
    marginTop: '50px',
  },
  profileTextWrapper: {
    height: '100%',
    margin: '0',
    padding: '0',
    paddingLeft: '10px',
    
    display: 'flex',
    justifyContent: 'center',
    // backgroundColor: 'orange',
    alignItems: 'flex-start',
    flexDirection: 'column',
  },
  profileGithub: {
    fontWeight: 'lighter',
    fontSize: '1rem',
    margin: '0',
    padding: '0',
    letterSpacing: '-.06rem',

    '&:hover': {
      cursor: 'pointer',
      textDecoration: 'underline',
      color: 'lightgreen'
    },
    '&:active': {
      cursor: 'pointer',
      textDecoration: 'underline',
      color: ThirdColor
    }
  }

});

export default ProfileScreen;