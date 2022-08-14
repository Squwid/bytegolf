import React from 'react';
import { makeStyles, Theme, withStyles } from '@material-ui/core/styles';
import { Modal, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from '@material-ui/core';
import { createStyles } from '@material-ui/styles';
import { PrimaryColor, SecondaryColor } from '../../Globals';
import AceEditor from 'react-ace';
import Chip, { ChipProps } from '../Chip/Chip';
import { LoginNotification } from '../Notification/LoginNotification/LoginNotification';
import { useQuery } from 'react-query';
import { GetFullSubmission, GetMyBestSubmission, GetMySubmissions } from '../../Store/Subs';
import { BasicShortSubmission } from '../../Types';
import Notification from '../Notification/Notification';
import { BasicLoadingIcon } from '../BasicLoadingIcon';
import TimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en'
import TestCases from '../TestCases/TestCases';

import 'ace-builds/src-noconflict/ace';
import 'ace-builds/src-noconflict/mode-golang';
import 'ace-builds/src-noconflict/mode-javascript';
import 'ace-builds/src-noconflict/mode-typescript';
import 'ace-builds/src-noconflict/mode-python';
import 'ace-builds/src-noconflict/mode-c_cpp';
import 'ace-builds/src-noconflict/theme-crimson_editor';


const modalStyles = makeStyles({
  modal: {
    minHeight: '300px',
    maxHeight: '75%',
    height: 'auto',
    width: '850px',
    backgroundColor: 'white',
    padding: '15px',
    margin: '0 auto',
    marginTop: '100px',
    overflowY: 'scroll',
    // borderRadius: '5px',

    fontFamily: 'FiraCode'
  },
  centerText: {
    textAlign: 'center',
    margin: '0'
  },
  chipHolder: {
    height: 'auto',
    width: '80%',
    margin: '0 auto',
    marginTop: '15px',
    marginBottom: '10px',

    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-around',
    flexWrap: 'wrap'
  }
});

const SubmissionModal: React.FC<{id?: string, open: boolean, onClose: () => void}> = (props) => {
  TimeAgo.addLocale(en)
  const ta = new TimeAgo('en-US');

  const id = props.id ? props.id : '';
  const classes = modalStyles();
  
  const sub = useQuery(['Submissions', id], () => GetFullSubmission(id));
  if (sub.isLoading) return (
    <Modal
      open={props.open && !!props.id}
      aria-labelledby="Title"
      aria-describedby="Description"
      onClose={props.onClose}
    >
      <div className={classes.modal}>
        <BasicLoadingIcon />
      </div>
    </Modal>
  );
  if (sub.isError) return (
    <Modal
      open={props.open && !!props.id}
      aria-labelledby="Title"
      aria-describedby="Description"
      onClose={props.onClose}
    >
      <div className={classes.modal}>
      <Notification type='error' text={`${sub.error}`} style={{marginBottom: '1rem'}}/>
      </div>
    </Modal>
  );
  if (sub.data === 'not found' || sub.data === 'not logged in' || !sub.data) return (
    <Modal
      open={props.open && !!props.id}
      aria-labelledby="Title"
      aria-describedby="Description"
      onClose={props.onClose}
    >
      <div className={classes.modal}>
      <Notification type='warn' text={"SUBMISSION WAS NOT FOUND"} style={{marginBottom: '1rem'}}/>
      </div>
    </Modal>
  );

  let chips: ChipProps[] = [];
  if (sub.data) {
    chips = [
      {ckey: 'HOLE NAME', value: sub.data.HoleName.toUpperCase(), bgColor: PrimaryColor, secondaryTextColor: 'white'},
      {ckey: 'SCORE', value: `${sub.data.Length}`, bgColor: PrimaryColor, secondaryTextColor: 'white'},
      {ckey: 'LANGUAGE', value: sub.data.Language, bgColor: PrimaryColor, secondaryTextColor: 'white'},
      {ckey: 'CORRECT', value: `${sub.data.Correct}`.toUpperCase(), bgColor: sub.data.Correct ? PrimaryColor : SecondaryColor, secondaryTextColor: 'white'},
      {ckey: 'SUBMITTED', value: ta.format(new Date(sub.data.SubmittedTime), 'round-minute').toUpperCase(), bgColor: PrimaryColor, secondaryTextColor: 'white'}
      
    ]
  }

  return (
    <Modal
      open={props.open && !!props.id}
      aria-labelledby="Title"
      aria-describedby="Description"
      onClose={props.onClose}
    >
      <div className={classes.modal}>
        <p className={classes.centerText} style={{fontWeight: 'lighter', fontSize: '1.5rem', marginBottom: '1rem'}}>{sub.data.ID}</p>
        <div>
          <div className={classes.chipHolder}>
            {chips.map(chip => <Chip key={chip.ckey} {...chip} style={{marginRight: '5px', marginBottom: '10px'}}/>)}
          </div>
          <AceEditor 
            mode={sub.data.Language}
            theme="crimson_editor"
            style={{width: '80%', height: 'auto', minHeight: '500px', margin: '0 auto'}}
            readOnly={true}
            defaultValue={sub.data.Script}
            wrapEnabled={true}
          />

          <TestCases
            holeID={sub.data.HoleID}
            tests={sub.data.Tests}
            sizeStyle={{width: '80%', height: 'auto', margin: '0 auto'}}
          />
        </div>
      </div>

    </Modal>
  );
}

type Props = {
  hole?: string;
}
const MySubmissions: React.FC<Props> = (props) => {
  TimeAgo.addLocale(en)
  const ta = new TimeAgo('en-US');

  const classes = useStyles();
  const [subModal, setSubModal] = React.useState<string|undefined>(undefined);
  
  const bestSubmission = useQuery(['BestSubmission', props?.hole], ()=>GetMyBestSubmission(props?.hole));

  let submissions = useQuery(['Submissions', props?.hole], () => GetMySubmissions(props?.hole));
  if (submissions.isLoading || bestSubmission.isLoading) return (<BasicLoadingIcon />);
  if (submissions.isError) return (<Notification type='error' text={`${submissions.error}`} />);
  if (!submissions.data) return (<LoginNotification />)
  if (submissions.data.length === 0) return (<p style={{fontFamily:'FiraCode', fontWeight: 'lighter', textAlign: 'center'}}>NO SUBMISSIONS YET</p>);

  const singleHole = !!props.hole;

  // Remove the best submission from the array of submissions and place at the top only if a single hole data is shown
  if (bestSubmission.data && singleHole) submissions.data = submissions.data.filter(sub => sub.ID !== bestSubmission.data?.ID);
  
  const onClick = (id: string) => {
    if (subModal) return;
    setSubModal(id);
  }
  
  const onClose = () => {
    setSubModal(undefined);
  }
  
  if (singleHole) {  
    const row = (sub: BasicShortSubmission, best?: boolean) => (
      <TRow className={best ? classes.best : sub.Correct ? classes.correct : classes.incorrect} key={sub.ID} onClick={() => onClick(sub.ID)}>
        <TCell padding={'none'} style={{padding:'5px'}} component="th" scope="row">{sub.ID.substr(0, 8)}</TCell>
        <TCell padding={'none'} style={{padding:'5px'}} align='right' scope="row">{ta.format(new Date(sub.SubmittedTime), 'round-minute').toUpperCase()}</TCell>
        <TCell padding={'none'} style={{padding:'5px'}} align='right' scope="row">{sub.Length}</TCell>
        <TCell padding={'none'} style={{padding:'5px'}} align='right' scope="row">{sub.Language}</TCell>
      </TRow>
    );
    
    return (
      <>
      <TableContainer>
        <Table className={classes.table}>
          <TableHead>
            <TableRow>
              <TCell>ID</TCell>
              <TCell align='right'>SUBMITTED</TCell>
              <TCell align='right'>SCORE</TCell>
              <TCell align='right'>LANGUAGE</TCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {bestSubmission.data && (row(bestSubmission.data, true))}
            {submissions.data.map(sub => row(sub))}
          </TableBody>
        </Table>
      </TableContainer>
      <SubmissionModal id={subModal} onClose={onClose} open={subModal!==undefined}/>
      </>
    )
  } else {
    // TODO: Add pagination to the table at some point

    const row = (sub: BasicShortSubmission) => (
      <TRow className={sub.Correct ? classes.correct : classes.incorrect} key={sub.ID} onClick={() => onClick(sub.ID)}>
        <TCell padding={'none'} style={{padding:'5px'}} component="th" scope="row">{sub.ID.substr(0, 8)}</TCell>
        <TCell padding={'none'} style={{padding:'5px'}} align='left' scope="row">{sub.HoleName.substr(0,35)}</TCell>
        <TCell padding={'none'} style={{padding:'5px'}} align='right' scope="row">{ta.format(new Date(sub.SubmittedTime), 'round-minute').toUpperCase()}</TCell>
        <TCell padding={'none'} style={{padding:'5px'}} align='right' scope="row">{sub.Length}</TCell>
        <TCell padding={'none'} style={{padding:'5px'}} align='right' scope="row">{sub.Language}</TCell>
      </TRow>
    );
    
    return (
      <>
      <TableContainer>
        <Table className={classes.table}>
          <TableHead>
            <TableRow>
              <TCell>ID</TCell>
              <TCell align='left'>HOLE NAME</TCell>
              <TCell align='right'>SUBMITTED</TCell>
              <TCell align='right'>SCORE</TCell>
              <TCell align='right'>LANGUAGE</TCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {submissions.data.map(sub => row(sub))}
          </TableBody>
        </Table>
      </TableContainer>
      <SubmissionModal id={subModal} onClose={onClose} open={subModal!==undefined}/>
      </>
    )
  }

}

const TCell = withStyles((theme: Theme) => createStyles({
  head: {
    backgroundColor: SecondaryColor,
    color: theme.palette.common.white,
    fontFamily: 'FiraCode',
    fontWeight: 'bold',
  },
  body: {
    fontSize: '1rem',
    fontFamily: 'FiraCode',
    fontWeight: 'lighter',
    letterSpacing: '-.09rem'
  }
})) (TableCell);

const TRow = withStyles((theme: Theme) => createStyles({
  root: {
    cursor: 'pointer',
  }
})) (TableRow);

const useStyles = makeStyles({
  table: {
    width: '980px',
    border: '1px solid #CCC',
    marginBottom: '5rem',
    marginTop: '1rem'
  },
  correct: {
    '&:hover': {
      backgroundColor: 'lightgreen'
    },
    '&:active': {
      backgroundColor: '#2BDB7C'
    },
    backgroundColor: '#BCFFC3'
  },
  incorrect: {
    backgroundColor: '#FFCFC4',
    '&:hover': {
      backgroundColor: '#FFB6A5'
    },
    '&:active': {
      backgroundColor: '#FF9D87'
    }
  },
  best: {
    backgroundColor: '#FFF4BB',
    '&:hover': {
      backgroundColor: '#FFED8D'
    },
    '&:active': {
      backgroundColor: '#FFE34F'
    }
  }
});

export default MySubmissions;