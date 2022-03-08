import React from 'react';
import { useParams } from 'react-router';
import { Difficulty } from '../Components/Difficulty';
import Nav from '../Components/Nav/Nav';
import { FormControl, MenuItem, Select } from '@material-ui/core';
import { AvailableLanguages, getLanguage, Language } from '../Languages';
import AceEditor from 'react-ace';
import { BackendURL, PrimaryColor, SecondaryColor } from '../Globals';
import Button from '../Components/Button/Button';
import MySubmissions from '../Components/MySubmissions/MySubmissions';
import Notification, { NotificationProps } from '../Components/Notification/Notification';
import { useHistory, useLocation } from 'react-router-dom';
import { GetHole } from '../Store/Holes';
import { useQuery, useQueryClient } from 'react-query';
import NotFoundScreen from './PsuedoScreens/NotFoundScreen';
import ErrorScreen from './PsuedoScreens/ErrorScreen';
import LoadingScreen from './PsuedoScreens/LoadingScreens';
import { GetClaims } from '../Store/Profile';
import './screen.css';
import 'ace-builds/src-noconflict/ace';
import 'ace-builds/src-noconflict/mode-python';
import 'ace-builds/src-noconflict/mode-java';
import 'ace-builds/src-noconflict/mode-javascript';
import 'ace-builds/src-noconflict/mode-c_cpp';
import 'ace-builds/src-noconflict/mode-php';
import 'ace-builds/src-noconflict/mode-rust';
import 'ace-builds/src-noconflict/mode-typescript';
import 'ace-builds/src-noconflict/mode-golang';
import 'ace-builds/src-noconflict/mode-powershell';
import 'ace-builds/src-noconflict/mode-batchfile';
import 'ace-builds/src-noconflict/theme-textmate';
import { SubmissionResponse } from '../Types';

const PlayScreen: React.FC = () => {
  const { pathname } = useLocation();
  const { holeID } = useParams<{holeID: string}>();
  const history = useHistory();

  const [ submitting, setSubmitting ] = React.useState(false);
  const [ activeLanguage, setActiveLanguage ] = React.useState<Language>(getLanguage('python3'));
  const [ len, setLen ] = React.useState(0);
  const [ script, setScript ] = React.useState('');
  const [ alerts, setAlerts ] = React.useState<NotificationProps[]>([]);

  const queryClient = useQueryClient();

  const onSubmit = () => {
    if (len === 0) return setAlerts([{type: 'info', text: 'LENGTH CANNOT BE 0'}]);

    if (submitting) return;
    setSubmitting(true);

    const url = `${BackendURL()}/submit/${holeID}`;
    const body = {
      script: script,
      language: activeLanguage.language,
      version: activeLanguage.version,
    };
    fetch(url, {credentials: 'include', body: JSON.stringify(body), method: 'POST'})
      .then(async resp => {
        setSubmitting(false);
        setAlerts([]);

        if (resp.status === 401) return setAlerts([...alerts, {type: 'warn', text: 'NOT LOGGED IN, CLICK HERE TO LOGIN TO PLAY', onClick: ()=>history.push('/profile')}]);
        if (resp.status === 500 || resp.status !== 200) return setAlerts([...alerts, {type: 'error', text: `INTERNAL SERVER ERROR, TRY AGAIN LATER`}]);
        const sr: SubmissionResponse = await resp.json();

        queryClient.invalidateQueries(['Submissions', holeID]);

        if (sr.BestScore) {
          queryClient.invalidateQueries(['BestSubmission', holeID]);
          setAlerts([{type: 'info', text: `NEW HIGH SCORE OF ${sr.Length} (${sr.ID.substr(0, 8)})`}]);
        } else {
          if (sr.Correct) {
            setAlerts([{type: 'info', text: `CORRECT SUBMISSION OF LENGTH ${sr.Length} (${sr.ID.substr(0, 8)})`}]);
          } else {
            setAlerts([{type: 'error', text: `INCORRECT SUBMISSION OF LENGTH ${sr.Length} (${sr.ID.substr(0, 8)})`}]);
          }
        }

      })
      .catch(err => {
        setSubmitting(false);
        setAlerts([{type: 'error', text: `INTERNAL SERVER ERROR, TRY AGAIN LATER`}]);
      });
    
    


    //   setSubmitting(false);
    //   setAlerts([...alerts, {type: 'info', text: 'NOT LOGGED IN, CLICK HERE TO LOGIN TO PLAY'}]);
    // }, 3000);
  }

  React.useEffect(() => {
    window.scrollTo(0, 0);
    document.title = `Bytegolf - Play`;    
  }, [pathname]);

  React.useEffect(() => {
    
  }, [submitting, alerts]);
  
  const claims = useQuery('Claims', () => GetClaims());
  const hole = useQuery(['Hole', holeID], () => GetHole(holeID));

  if (hole.isLoading || claims.isLoading) return (<LoadingScreen active='play'/>);
  if (hole.isError || claims.isError) return (<ErrorScreen active='play' text={`${hole.error}`} />);
  if (!hole.data) return (<NotFoundScreen active='play' text={`Hole ${holeID} was not found`} />);
  // if (!claims.data) setAlerts([...alerts, {type: 'warn', text: 'NOT LOGGED IN, CLICK HERE TO LOGIN TO PLAY', onClick: () => history.push('/profile')}]);

  const selectButton = (): JSX.Element => {
    return (
      <FormControl style={{width: '25%'}} variant='outlined'>
        <Select
          labelId='language-select'
          id='language-selection'
          value={activeLanguage.language}
          onChange={(e: React.ChangeEvent<{value: unknown}>) => {setActiveLanguage(getLanguage(e.target.value as string))}}
        >
          {AvailableLanguages.map((lang, i) => <MenuItem key={i} value={lang.language}>{lang.name}</MenuItem>)}
        </Select>
      </FormControl>
    );
  }

  return (
    <div>
      <Nav active='play'/>
      <div className='screenContainer'>
        <p className='screenTitle' style={{fontSize: '2.3rem', letterSpacing: '-0.09rem'}}>{hole.data.Name.toUpperCase()}</p>
        <p className='screenSubText' style={{fontSize: '1.5rem'}}>{hole.data?.Question}</p>
        <Difficulty style={{fontSize:'1.3rem', textAlign: 'center'}} difficulty={hole.data.Difficulty} />
        <div style={{display: 'flex', justifyContent: 'center', flexWrap: 'nowrap', flexDirection: 'row'}}>
          <p style={{marginRight: '1rem'}}>LANGUAGE:</p>
          {selectButton()}
        </div>

        {/* Notifications, but multiple just in case */}
        <div style={{margin: '0 auto', marginTop: '1rem', display: 'flex', flexDirection: 'column', width: '80%'}}>
          {alerts.map((a: NotificationProps, i: number) => (<Notification key={i} {...a} style={{marginBottom: '1rem'}} />))}
        </div>

        <div style={{margin: '0 auto', display: 'flex', flexDirection: 'row-reverse', flexWrap: 'nowrap', width: '80%'}}>
          <p style={{margin: 0, padding: 0, paddingRight: '10px'}}>BYTES: {len}</p>
        </div>
        <AceEditor
          mode={activeLanguage.editorValue}
          theme='textmate'
          style={{width: '80%', height: '600px', margin: '0 auto'}}
          readOnly={submitting}
          wrapEnabled={true}
          onChange={(val: string) => { setScript(val); setLen(val.length)}}
          value={script}
        />

        <div style={{display: 'flex', flexDirection: 'row', justifyContent: 'center', height: 'auto'}}>
          <Button disabled={submitting} loading={submitting} color={PrimaryColor} activeColor='white' fontSize='1.1rem' onPress={onSubmit} text='SUBMIT' />
          <Button disabled={submitting} color={SecondaryColor} activeColor='white' fontSize='1.1rem' text='CLEAR' onPress={() => {setScript(''); setLen(0)}} />
        </div>

        <div style={{marginTop: '2rem'}}>
          <p className='screenSubText' style={{marginBottom: '20px'}}>PREVIOUS SUBMISSIONS</p>
          <MySubmissions hole={holeID} />
        </div>


      </div>
    </div>
  )
}

export default PlayScreen;